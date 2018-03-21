package app

import (
	"flag"

	"fmt"
	"github.com/TIBCOSoftware/flogo-cli/cli"
	"os"
)

var optPrepare = &cli.OptionInfo{
	Name:      "prepare",
	UsageLine: "prepare [-o][-e]",
	Short:     "prepare the flogo application",
	Long: `[Deprecated, use 'build [-gen]' instead]Prepare the flogo application.

Options:
    -o   optimize for directly referenced contributions
    -e   embed application configuration into executable
`,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdPrepare{option: optPrepare})
}

type cmdPrepare struct {
	option      *cli.OptionInfo
	optimize    bool
	embedConfig bool
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdPrepare) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdPrepare) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.optimize), "o", false, "optimize prepare")
	fs.BoolVar(&(c.embedConfig), "e", false, "embed config")
}

// Exec implementation of cli.Command.Exec
func (c *cmdPrepare) Exec(args []string) error {

	appDir, err := os.Getwd()

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to determine working directory\n\n")
		os.Exit(2)
	}

	options := &PrepareOptions{OptimizeImports: c.optimize, EmbedConfig: c.embedConfig}
	return PrepareApp(SetupExistingProjectEnv(appDir), options)
}
