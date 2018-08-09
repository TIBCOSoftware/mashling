//go:generate go run ../generate/stub_generator.go

package command

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/TIBCOSoftware/mashling/internal/app/cli/assets"
	"github.com/TIBCOSoftware/mashling/internal/pkg/grpcsupport"
	gwerrors "github.com/TIBCOSoftware/mashling/internal/pkg/model/errors"
	"github.com/TIBCOSoftware/mashling/pkg/files"
	"github.com/TIBCOSoftware/mashling/pkg/strings"
	"github.com/spf13/cobra"
)

const (
	// ImportPath is the root import path regardless of location.
	ImportPath = "github.com/TIBCOSoftware/mashling"
	// DockerImage is the Docker image used to run the creation process.
	DockerImage = "mashling/mashling-compile:0.4.0"
)

func init() {
	createCommand.Flags().StringVarP(&name, "name", "n", "mashling-custom", "customized mashling-gateway name")
	createCommand.Flags().StringVarP(&protoPath, "protoPath", "p", "", "path to proto file for grpc service")
	createCommand.Flags().BoolVarP(&native, "native", "N", false, "build the customized binary natively instead of using Docker")
	createCommand.Flags().StringVarP(&targetOS, "os", "O", "", "target OS to build for (default is the host OS, valid values are windows, darwin, and linux)")
	createCommand.Flags().StringVarP(&targetArch, "arch", "A", "", "target architecture to build for (default is amd64, arm64 is only compatible with Linux)")
	cliCommand.AddCommand(createCommand)
}

var (
	protoPath           string
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
			case *gwerrors.UndefinedReference:
				log.Fatalf("%s: %s", e.Type(), e.Details())
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
	fullPathName := filepath.Join(name, "src", ImportPath)

	Env := os.Environ()
	Env = append(Env, "GOPATH="+name)
	Env = append(Env, "PATH="+os.Getenv("PATH")+":"+filepath.Join(name, "bin"))

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
	if _, err = os.Stat(fullPathName); os.IsNotExist(err) {
		err = os.MkdirAll(fullPathName, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	if _, err = os.Stat(filepath.Join(name, "bin")); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join(name, "bin"), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	stub, err := assets.Asset("stub.zip")
	if err != nil {
		log.Fatal(err)
	}
	err = files.UnpackBytes(stub, fullPathName)
	if err != nil {
		log.Fatal(err)
	}

	//grpc Support code
	gRPCFlag := false
	if len(protoPath) != 0 {
		gRPCFlag = true
	}
	if gRPCFlag {
		log.Println("Generating grpc support files using proto file: ", protoPath)
		grpcsupport.AssignValues(name)
		err := grpcsupport.GenerateSupportFiles(protoPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	//grpc Support code end

	var cmd *exec.Cmd
	var dockerCmd, dockerContainerID string
	if dockerCmd, err = exec.LookPath("docker"); native || err != nil {
		// Docker does not exist, try native toolchain.
		log.Println("Docker not found or native option specified, using go natively...")
		dockerCmd = ""
	} else {
		log.Println("Docker found, using it to build...")
		cmd = exec.Command(dockerCmd, "run", "--rm", "-d", "-t", DockerImage)
		cmd.Dir = name
		cmd.Env = Env
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
			cmd.Env = Env
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
		cmd.Env = Env
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
	cmd.Dir = fullPathName
	cmd.Env = Env
	output, cErr := cmd.CombinedOutput()
	if cErr != nil {
		log.Println(string(output))
		log.Fatal(cErr)
	}
	// Run dep add for all identified new dependencies
	if len(deps) > 0 {
		// Turn deps into a string
		log.Println("Installing missing dependencies...")
		depString := strings.Join(util.UniqueStrings(deps), " ")
		if dockerCmd != "" {
			cmd = exec.Command(dockerCmd, "exec", dockerContainerID, "/bin/bash", "-c", "dep ensure -add "+depString)
		} else {
			cmd = exec.Command("dep", "ensure", "-add", depString)
		}
		cmd.Dir = fullPathName
		cmd.Env = Env
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
	cmd.Dir = fullPathName
	cmd.Env = Env
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
	cmd.Dir = fullPathName
	cmd.Env = Env
	output, cErr = cmd.CombinedOutput()
	if cErr != nil {
		log.Println(string(output))
		log.Fatal(cErr)
	}
	if dockerCmd != "" {
		log.Println("Copying out created source code and binary from container...")
		// Copy out created source directory from running container.
		cmd = exec.Command(dockerCmd, "cp", dockerContainerID+":/mashling/src/"+ImportPath+"/.", filepath.Join(name, "src", ImportPath))
		cmd.Dir = name
		cmd.Env = Env
		output, cErr = cmd.CombinedOutput()
		if cErr != nil {
			log.Println(string(output))
			log.Fatal(cErr)
		}
	}
	// Copy release folder contents to top level
	err = filepath.Walk(filepath.Join(name, "src", ImportPath, "release"), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			err = files.CopyFile(path, filepath.Join(name, info.Name()))
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
