package command

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	cliCommand.AddCommand(validateCommand)
}

var validateCommand = &cobra.Command{
	Use:   "validate",
	Short: "Validates a mashling.json configuration file",
	Long:  `Validates a provided mashling.json configuration file based off of the supported Mashling schema versions`,
	Run:   validate,
}

func validate(command *cobra.Command, args []string) {
	err := loadGateway()
	if err != nil {
		log.Fatal("Invalid configuration file: ", err)
	} else {
		log.Println("Valid.")
	}
}
