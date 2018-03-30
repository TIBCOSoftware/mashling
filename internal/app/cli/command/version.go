package command

import (
	"fmt"

	"github.com/TIBCOSoftware/mashling/internal/app/version"
	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Prints the mashling-cli version",
	Long:  `Prints the mashling-cli version and build details`,
	Run:   ver,
}

func init() {
	rootCommand.AddCommand(versionCommand)
}

func ver(command *cobra.Command, args []string) {
	fmt.Println("Version: ", version.Version)
	fmt.Println("Build Date: ", version.BuildDate)
}
