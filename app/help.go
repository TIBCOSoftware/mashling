// Command, OptionInfo and command execution pattern derived from
// github.com/constabulary/gb, released under MIT license
// https://github.com/constabulary/gb/blob/master/LICENSE
package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/TIBCOSoftware/mashling-cli/cli"
)

var optHelp = &cli.OptionInfo{
	Name:      "help",
	UsageLine: "help [command]",
	Short:     "Get help for a command",
	Long: `Get help for a mashling command.

`,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdHelp{option: optHelp})
}

type cmdHelp struct {
	option *cli.OptionInfo
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdHelp) OptionInfo() *cli.OptionInfo {
	return c.option
}

// Exec implementation of cli.Command.Exec
func (c *cmdHelp) AddFlags(fs *flag.FlagSet) {
	//op op
}

// Exec implementation of cli.Command.Exec
func (c *cmdHelp) Exec(args []string) error {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return nil
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: mashling help command\n\nToo many arguments given.\n")
		os.Exit(2)
	}

	arg := args[0]

	cmd, exists := CommandRegistry.Command(arg)

	if exists {
		cli.PrintCmdHelp("", cmd)
		return nil
	}

	tool, exists := cli.GetTool(arg)

	if exists {
		fgutil.RenderTemplate(os.Stdout, "{{.Long}}\n\n", tool.OptionInfo())
		tool.PrintUsage(os.Stdout)
		return nil
	}

	fmt.Fprintf(os.Stderr, "Unknown help command %#q. Run 'gateway help'.\n", arg)
	os.Exit(2)

	return nil
}
