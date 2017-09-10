package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/mashling-cli/cli"
)

var optList = &cli.OptionInfo{
	Name:      "list",
	UsageLine: "list [triggers|handlers|links|all]",
	Short:     "list installed components in the mashling gateway recipe",
	Long:      `List installed components in the mashling gateway recipe.`,
}

const (
	ctTriggers = "triggers"
	ctHandlers = "handlers"
	ctLinks    = "links"
	ctAll      = "all"
)

func init() {
	CommandRegistry.RegisterCommand(&cmdList{option: optList})
}

type cmdList struct {
	option *cli.OptionInfo
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdList) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdList) AddFlags(fs *flag.FlagSet) {
}

// Exec implementation of cli.Command.Exec
func (c *cmdList) Exec(args []string) error {

	if len(args) > 1 {
		fmt.Fprint(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	var cType ComponentType

	if len(args) == 1 {
		listCT := args[0]

		switch listCT {
		case ctLinks:
			cType = LINK
		case ctTriggers:
			cType = TRIGGER
		case ctHandlers:
			cType = HANDLER
		case ctAll:
			cType = ALL
		default:
			fmt.Fprintf(os.Stderr, "Error: Unknown component type - %s\n\n", listCT)
			cmdUsage(c)
		}
	} else {
		cmdUsage(c)
	}

	appDir, err := os.Getwd()

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to determine working directory\n\n")
		os.Exit(2)
	}

	gwDetails, err := GetGatewayDetails(SetupExistingProjectEnv(appDir), cType)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, gwDetails)
	return nil
}
