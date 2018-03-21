package app

import (
	"flag"
	"fmt"
	"os"

	"encoding/json"
	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/TIBCOSoftware/flogo-cli/config"
)

var optList = &cli.OptionInfo{
	Name:      "list",
	UsageLine: "list [-json] [actions|triggers|activities]",
	Short:     "list installed contributions",
	Long: `Lists installed contributions.

Options:
    -json generate output as json
`,
}

const (
	ctActions   = "actions"
	ctTriggers  = "triggers"
	ctActivites = "activities"
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
	c.outputJson = fs.Bool("json", false, "generate output as json")
}

// Exec implementation of cli.Command.Exec
func (c *cmdList) Exec(args []string) error {

	if len(args) > 1 {
		fmt.Fprint(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	var cType config.ContribType

	if len(args) == 1 {
		listCT := args[0]

		switch listCT {
		case ctActions:
			cType = config.ACTION
		case ctTriggers:
			cType = config.TRIGGER
		case ctActivites:
			cType = config.ACTIVITY
		case "flow-models":
			cType = config.FLOW_MODEL
		default:
			fmt.Fprintf(os.Stderr, "Error: Unknown contribution type - %s\n\n", listCT)
			cmdUsage(c)
		}
	}

	appDir, err := os.Getwd()

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to determine working directory\n\n")
		os.Exit(2)
	}

	dependencies, err := ListDependencies(SetupExistingProjectEnv(appDir), cType)

	if err != nil {
		return err
	}

	if *c.outputJson {
		depJson, err := json.MarshalIndent(dependencies, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, string(depJson))
	} else {
		byType := make(map[string][]string)

		//aggregate by ContribType
		for _, dependency := range dependencies {

			switch dependency.ContribType {
			case config.ACTION:
				byType[ctActions] = append(byType[ctActions], dependency.Ref)
			case config.TRIGGER:
				byType[ctTriggers] = append(byType[ctTriggers], dependency.Ref)
			case config.ACTIVITY:
				byType[ctActivites] = append(byType[ctActivites], dependency.Ref)
			default:
				byType[dependency.ContribType.String()] = append(byType[dependency.ContribType.String()], dependency.Ref)
			}
		}

		for ct, refs := range byType {

			fmt.Fprintf(os.Stdout, "%s:\n", ct)

			for _, ref := range refs {

				fmt.Fprintf(os.Stdout, "  %s\n", ref)
			}
		}
	}

	return nil
}
