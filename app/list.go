package app

import (
	"flag"
	"fmt"
	"os"

	"encoding/json"
	"github.com/TIBCOSoftware/mashling-cli/cli"
)

var optList = &cli.OptionInfo{
	Name:      "list",
	UsageLine: "list [links|triggers|handlers]",
	Short:     "list installed components in the mashling gateway recipe",
	Long: `List installed components in the mashling gateway recipe.

Options:
    -json generate output as json`,
}

const (
	ctLinks    = "links"
	ctTriggers = "triggers"
	ctHandlers = "handlers"
)

func init() {
	CommandRegistry.RegisterCommand(&cmdList{option: optList})
}

type cmdList struct {
	option     *cli.OptionInfo
	outputJson *bool
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdList) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdList) AddFlags(fs *flag.FlagSet) {
	c.outputJson = fs.Bool("json", true, "generate output as json")
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

	var components interface{}
	if cType == LINK {
		components, err = ListLinks(SetupExistingProjectEnv(appDir), cType)
	} else {
		components, err = ListComponents(SetupExistingProjectEnv(appDir), cType)
	}

	if err != nil {
		return err
	}

	json, err := json.MarshalIndent(components, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, string(json))

	return nil
}
