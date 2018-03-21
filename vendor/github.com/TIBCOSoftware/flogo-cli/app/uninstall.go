package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-cli/cli"
)

var optUninstall = &cli.OptionInfo{
	Name:      "uninstall",
	UsageLine: "uninstall contribution",
	Short:     "uninstall a flogo contribution",
	Long: `Uninstalls a flogo contribution.
`,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdUninstall{option: optUninstall})
}

type cmdUninstall struct {
	option *cli.OptionInfo
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdUninstall) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdUninstall) AddFlags(fs *flag.FlagSet) {
}

// Exec implementation of cli.Command.Exec
func (c *cmdUninstall) Exec(args []string) error {

	if len(args) == 0 {
		fmt.Fprint(os.Stderr, "Error: contribution not specified\n\n")
		cmdUsage(c)
	}

	if len(args) != 1 {
		fmt.Fprint(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	contribPath := args[0]

	appDir, err := os.Getwd()

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to determine working directory\n\n")
		os.Exit(2)
	}

	return UninstallDependency(SetupExistingProjectEnv(appDir), contribPath)
}
