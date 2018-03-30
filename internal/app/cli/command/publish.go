package command

import (
	"github.com/spf13/cobra"
)

var publishCommand = &cobra.Command{
	Use:   "publish",
	Short: "Publishes to supported platforms",
	Long:  `Publishes details of the mashling.json configuration file to various support platforms (currently Mashery and Consul)`,
}

func init() {
	rootCommand.AddCommand(publishCommand)
}
