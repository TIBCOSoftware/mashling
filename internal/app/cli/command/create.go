//go:generate go run ../generate/stub_generator.go

package command

import (
	"bytes"
	"fmt"
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
	var dockerCmd, dockerContainerID string
	if dockerCmd, err = exec.LookPath("docker"); native || err != nil {
		// Docker does not exist, try native toolchain.
		log.Println("Docker not found or native option specified, using make natively...")
		dockerCmd = ""
	} else {
		log.Println("Docker found, using it to build...")
		cmd = exec.Command(dockerCmd, "run", "--rm", "-d", "-t", "mashling/mashling-compile")
		output, cErr := cmd.Output()
		if cErr != nil {
			log.Println(string(output))
			log.Fatal(cErr)
		}
		dockerContainerID = strings.TrimSpace(string(output))
		defer func() {
			log.Println("Stopping container: ", dockerContainerID)
			// Stop running container.
			cmd = exec.Command(dockerCmd, "stop", dockerContainerID)
			cmd.Dir = name
			output, cErr = cmd.CombinedOutput()
			if cErr != nil {
				log.Println(string(output))
				log.Fatal(cErr)
			}
		}()
		log.Println("Copying default source code into container:", dockerContainerID)
		// Copy default source into container.
		cmd = exec.Command(dockerCmd, "cp", name+"/.", dockerContainerID+":/mashling/")
		cmd.Dir = name
		output, cErr = cmd.CombinedOutput()
		if cErr != nil {
			log.Println(string(output))
			log.Fatal(cErr)
		}
	}
	// Setup environment
	log.Println("Setting up project...")
	if dockerCmd != "" {
		cmd = exec.Command(dockerCmd, "exec", dockerContainerID, "/bin/bash", "-c", "go run build.go setup")
	} else {
		cmd = exec.Command("go", "run", "build.go", "setup")
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
		buffer.WriteString("-newdeps=\"")
		buffer.WriteString(strings.Join(util.UniqueStrings(deps), " "))
		buffer.WriteString("\"")
		if dockerCmd != "" {
			cmd = exec.Command(dockerCmd, "exec", dockerContainerID, "/bin/bash", "-c", "go run build.go depadd "+buffer.String())
		} else {
			cmd = exec.Command("go", "run", "build.go", "depadd", buffer.String())
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
		cmd = exec.Command(dockerCmd, "exec", dockerContainerID, "/bin/bash", "-c", "go run build.go allgatewayprep")
	} else {
		cmd = exec.Command("go", "run", "build.go", "allgatewayprep")
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
		cmd = exec.Command(dockerCmd, "exec", dockerContainerID, "/bin/bash", "-c", fmt.Sprintf("go run build.go releasegateway -os=%s -arch=%s", targetOS, targetArch))
	} else {
		cmd = exec.Command("go", "run", "build.go", "releasegateway", "-os="+targetOS, "-arch="+targetArch)
	}
	cmd.Dir = name
	output, cErr = cmd.CombinedOutput()
	if cErr != nil {
		log.Println(string(output))
		log.Fatal(cErr)
	}
	if dockerCmd != "" {
		log.Println("Copying out created source code and binary from container...")
		// Copy out created source directory from running container.
		cmd = exec.Command(dockerCmd, "cp", dockerContainerID+":/mashling/.", name)
		cmd.Dir = name
		output, cErr = cmd.CombinedOutput()
		if cErr != nil {
			log.Println(string(output))
			log.Fatal(cErr)
		}
	}
}
