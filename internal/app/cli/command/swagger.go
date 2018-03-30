package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	host    string
	trigger string
	output  string
)

var swaggerCommand = &cobra.Command{
	Use:   "swagger",
	Short: "Creates a swagger 2.0 doc",
	Long:  `Creates a swagger 2.0 doc based off of the HTTP triggers in the mashling.json configuration file`,
	Run:   swagger,
}

func init() {
	swaggerCommand.Flags().StringVarP(&host, "host", "H", "localhost", "the hostname where this mashling will be deployed")
	swaggerCommand.Flags().StringVarP(&trigger, "trigger", "t", "", "the trigger name to target (default is all))")
	swaggerCommand.Flags().StringVarP(&output, "output", "o", "", "the output file to write the swagger.json to (default is stdout)")
	rootCommand.AddCommand(swaggerCommand)
}

func swagger(command *cobra.Command, args []string) {
	if host == "" {
		log.Fatal("host is required")
	}

	err := loadGateway()
	if err != nil {
		log.Fatal(err)
	}

	docs, err := gateway.Swagger(host, trigger)
	if err != nil {
		log.Fatal("unable to generate Swagger representation: ", err)
	}
	if output == "" {
		fmt.Fprintf(os.Stdout, "%s\n", string(docs))
	} else {
		err := ioutil.WriteFile(output, docs, 0644)
		if err != nil {
			log.Fatal("not able write Swagger to output file: ", err)
		}
	}
}
