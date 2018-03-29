//go:generate go run ../generate/stub_generator.go

package command

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/TIBCOSoftware/mashling/internal/app/cli/assets"
	"github.com/TIBCOSoftware/mashling/pkg/files"
	"github.com/TIBCOSoftware/mashling/pkg/strings"
)

var supportedTargetOS = map[string]bool{"windows": true, "darwin": true, "linux": true}

// Create builds a custom mashling-gateway project directory populated with
// dependencies listed in the provided Mashling config file.
func Create(name string, deps []string, native bool, targetOS string) (err error) {
	if targetOS == "" {
		targetOS = runtime.GOOS
	}
	if _, ok := supportedTargetOS[targetOS]; !ok {
		return errors.New("invalid target OS type specified")
	}
	if _, err = os.Stat(name); os.IsNotExist(err) {
		err = os.MkdirAll(name, 0755)
		if err != nil {
			return err
		}
	}
	stub, err := assets.Asset("stub.zip")
	if err != nil {
		return err
	}
	err = files.UnpackBytes(stub, name)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	var dockerCmd string
	if dockerCmd, err = exec.LookPath("docker"); native || err != nil {
		// Docker does not exist, try native toolchain.
		log.Println("Docker not found or native option specified, using make natively...")
		dockerCmd = ""
	} else {
		log.Println("Docker found, using it to build...")
	}
	// Setup environment
	log.Println("Setting up project...")
	if dockerCmd != "" {
		cmd = exec.Command(dockerCmd, "run", "-v", name+":/mashling", "--rm", "-t", "jeffreybozek/mashling:compile", "/bin/bash", "-c", "make setup")
	} else {
		cmd = exec.Command("make", "setup")
	}
	cmd.Dir = name
	output, cErr := cmd.CombinedOutput()
	if cErr != nil {
		log.Println(string(output))
		return cErr
	}
	// Run dep add for all identified new dependencies
	if len(deps) > 0 {
		// Turn deps into a string
		log.Println("Installing missing dependencies...")
		var buffer bytes.Buffer
		buffer.WriteString("NEWDEPS=\"")
		buffer.WriteString(strings.Join(util.UniqueStrings(deps), " "))
		buffer.WriteString("\"")
		if dockerCmd != "" {
			cmd = exec.Command(dockerCmd, "run", "-v", name+":/mashling", "--rm", "-t", "jeffreybozek/mashling:compile", "/bin/bash", "-c", "make depadd "+buffer.String())
		} else {
			cmd = exec.Command("make", "depadd", buffer.String())
		}
		cmd.Dir = name
		output, cErr = cmd.CombinedOutput()
		if cErr != nil {
			log.Println(string(output))
			return cErr
		}
	}
	// Run make targets to generate appropriate code
	log.Println("Generating assets for customized Mashling...")
	if dockerCmd != "" {
		cmd = exec.Command(dockerCmd, "run", "-v", name+":/mashling", "--rm", "-t", "jeffreybozek/mashling:compile", "/bin/bash", "-c", "make assets generate fmt")
	} else {
		cmd = exec.Command("make", "assets", "generate", "fmt")
	}
	cmd.Dir = name
	output, cErr = cmd.CombinedOutput()
	if cErr != nil {
		log.Println(string(output))
		return cErr
	}
	// Run make build target to build for appropriate OS
	log.Println("Building customized Mashling binary...")
	if dockerCmd != "" {
		cmd = exec.Command(dockerCmd, "run", "-e", "GOOS="+targetOS, "-v", name+":/mashling", "--rm", "-t", "jeffreybozek/mashling:compile", "/bin/bash", "-c", "make buildgateway")
	} else {
		cmd = exec.Command("make", "buildgateway")
		env := os.Environ()
		env = append(env, "GOOS="+targetOS)
		cmd.Env = env
	}
	cmd.Dir = name
	output, cErr = cmd.CombinedOutput()
	if cErr != nil {
		log.Println(string(output))
		return cErr
	}
	return nil
}
