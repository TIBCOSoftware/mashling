package command

import (
	"log"

	"github.com/TIBCOSoftware/mashling/internal/pkg/model"
	"github.com/spf13/cobra"
)

func init() {
	cliCommand.PersistentFlags().StringVarP(&config, "config", "c", "mashling.json", "mashling gateway configuration")
	cliCommand.PersistentFlags().StringVarP(&envVarName, "env-var-name", "e", "MASHLING_CONFIG", "name of the environment variable that contains the base64 encoded mashling gateway configuration")
	cliCommand.PersistentFlags().BoolVarP(&loadFromEnv, "load-from-env", "l", false, "load the mashling gateway configuration from an environment variable")
}

var (
	config      string
	envVarName  string
	loadFromEnv bool
	gateway     model.Gateway
)

var cliCommand = &cobra.Command{
	Use:   "mashling-cli",
	Short: "mashling-cli is a CLI to help build mashling-gateway instances",
	Long:  "A CLI to build custom mashling-gateway instances, publish configurations to Mashery, and more. Complete documentation is available at https://github.com/TIBCOSoftware/mashling",
}

// Execute executes registered commands.
func Execute() {
	if err := cliCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}

func loadGateway() (err error) {
	// Load the appropriate Gateway instance from the configuration.
	if loadFromEnv {
		// Loading base64 encoded configuration from the env.
		gateway, err = model.LoadFromEnv(envVarName)
	} else {
		// Loading the configuration JSON from the specified file.
		gateway, err = model.LoadFromFile(config)
	}
	return err
}
