package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/TIBCOSoftware/mashling/internal/app/cli/command"
	model "github.com/TIBCOSoftware/mashling/internal/pkg/model"
	gwerrors "github.com/TIBCOSoftware/mashling/internal/pkg/model/errors"
)

var (
	flogoLibRevision = "52e50da7cdbe38eada4b1449992a296c48ae2349"
	mashlingRevision = "8cc2d8fe27b7252afb5744156096d2d77f792f5c"
	// Version is grabbed from the git repo at build time.
	Version = "dev"
	// BuildDate is set with the current timestamp at build time.
	BuildDate   = "unknown"
	gateway     model.Gateway
	config      string
	envVarName  string
	loadFromEnv bool
	name        string
	native      bool
	targetOS    string
)

func init() {
	flag.StringVar(&config, "config", "mashling.json", "mashling gateway configuration")
	flag.StringVar(&envVarName, "env-var-name", "MASHLING_CONFIG", "name of the environment variable that contain sthe base64 encoded mashling gateway configuration")
	flag.BoolVar(&loadFromEnv, "load-from-env", false, "load the mashling gateway configuration from an environment variable")
	flag.StringVar(&name, "name", "mashling-custom", "customized mashling-gateway name")
	flag.BoolVar(&native, "native", false, "build the customized binary natively instead of using Docker")
	flag.StringVar(&targetOS, "os", "", "target OS to build for (default is the host OS, valid values are windows, darwin, and linux)")
}

func main() {
	var err error
	var deps []string
	log.Println("[mashling] CLI Version: ", Version)
	log.Println("[mashling] Build Date: ", BuildDate)
	flag.Parse()

	// Load the appropriate Gateway instance from the configuration.
	if loadFromEnv {
		// Loading base64 encoded configuration from the env.
		gateway, err = model.LoadFromEnv(envVarName)
	} else {
		// Loading the configuration JSON from the specified file.
		gateway, err = model.LoadFromFile(config)
	}
	if err != nil {
		// Attempt to remedy any errors found, particularly missing dependencies.
		if gateway == nil {
			os.Exit(1)
		}
		for _, errd := range gateway.Errors() {
			switch e := errd.(type) {
			case *gwerrors.MissingDependency:
				log.Println("Missing dependencies found: ", strings.Join(e.MissingDependencies, " "))
				deps = append(deps, e.MissingDependencies...)
			default:
				fmt.Printf("Do not know how to handle error type %T!\n", e)
				os.Exit(1)
			}
		}
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	err = command.Create(filepath.Join(pwd, name), deps, native, targetOS)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
