package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/mashling-cli/cli"
)

var optBuild = &cli.OptionInfo{
	Name:      "build",
	UsageLine: "build",
	Short:     "build mashling gateway from mashling.json",
	Long:      "build mashling gateway from gateway description file - mashling.json",
}

func init() {
	CommandRegistry.RegisterCommand(&cmdBuild{option: optBuild})
}

type cmdBuild struct {
	option *cli.OptionInfo
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdBuild) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdBuild) AddFlags(fs *flag.FlagSet) {
	// c.outputJson = fs.Bool("json", true, "generate output as json")
}

// Exec implementation of cli.Command.Exec
func (c *cmdBuild) Exec(args []string) error {

	if len(args) > 1 {
		fmt.Fprint(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	fmt.Fprint(os.Stdout, "Building mashling gateway from mashling.json ...\n\n")

	return nil
}
