//go:generate go run ../generate/stub_generator.go

package command

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/TIBCOSoftware/mashling/internal/app/cli/assets"
	gwerrors "github.com/TIBCOSoftware/mashling/internal/pkg/model/errors"
	"github.com/TIBCOSoftware/mashling/pkg/files"
	"github.com/TIBCOSoftware/mashling/pkg/strings"
	"github.com/spf13/cobra"
)

func init() {
	createCommand.Flags().StringVarP(&name, "name", "n", "mashling-custom", "customized mashling-gateway name")
	createCommand.Flags().BoolVarP(&native, "native", "N", false, "build the customized binary natively instead of using Docker")
	createCommand.Flags().StringVarP(&targetOS, "os", "O", "", "target OS to build for (default is the host OS, valid values are windows, darwin, and linux)")
	createCommand.Flags().StringVarP(&targetArch, "arch", "A", "", "target architecture to build for (default is amd64, arm64 is only compatible with Linux)")
	cliCommand.AddCommand(createCommand)
}

var (
	name                string
	native              bool
	targetOS            string
	targetArch          string
	supportedTargetOS   = map[string]bool{"windows": true, "darwin": true, "linux": true}
	supportedTargetArch = map[string]bool{"amd64": true, "arm64": true}
)

var createCommand = &cobra.Command{
	Use:   "create",
	Short: "Creates a customized mashling-gateway",
	Long:  `Create a reusable customized mashling-gateway binary based off of the dependencies listed in your mashling.json configuration file`,
	Run:   create,
}

// Create builds a custom mashling-gateway project directory populated with
// dependencies listed in the provided Mashling config file.
func create(command *cobra.Command, args []string) {
	var deps []string

	err := loadGateway()

	if err != nil {
		// Attempt to remedy any errors found, particularly missing dependencies.
		if gateway == nil {
			log.Fatal(err)
		}
		for _, errd := range gateway.Errors() {
			switch e := errd.(type) {
			case *gwerrors.MissingDependency:
				log.Println("Missing dependencies found: ", strings.Join(e.MissingDependencies, " "))
				deps = append(deps, e.MissingDependencies...)
			default:
				log.Fatalf("Do not know how to handle error type %T!\n", e)
			}
		}
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	name = filepath.Join(pwd, name)

	if targetOS == "" {
		targetOS = runtime.GOOS
	}
	if targetArch == "" {
		targetArch = "amd64"
	}
	if _, ok := supportedTargetOS[targetOS]; !ok {
		log.Fatal("invalid target OS type specified")
	}
	if _, ok := supportedTargetArch[targetArch]; !ok {
		log.Fatal("invalid target architecture type specified")
	}
	if targetArch == "arm64" && targetOS != "linux" {
		log.Fatal("arm64 architecture is only valid with linux")
	}
	if _, err = os.Stat(name); os.IsNotExist(err) {
		err = os.MkdirAll(name, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	stub, err := assets.Asset("stub.zip")
	if err != nil {
		log.Fatal(err)
	}
	err = files.UnpackBytes(stub, name)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(cErr)
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
			log.Fatal(cErr)
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
		log.Fatal(cErr)
	}
	// Run make build target to build for appropriate OS
	log.Println("Building customized Mashling binary...")
	if dockerCmd != "" {
		cmd = exec.Command(dockerCmd, "run", "-e", "GOOS="+targetOS, "-e", "GOARCH="+targetArch, "-v", name+":/mashling", "--rm", "-t", "jeffreybozek/mashling:compile", "/bin/bash", "-c", "make buildgateway")
	} else {
		cmd = exec.Command("make", "buildgateway")
		env := os.Environ()
		env = append(env, "GOOS="+targetOS)
		env = append(env, "GOARCH="+targetArch)
		cmd.Env = env
	}
	cmd.Dir = name
	output, cErr = cmd.CombinedOutput()
	if cErr != nil {
		log.Println(string(output))
		log.Fatal(cErr)
	}
}
