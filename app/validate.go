package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/TIBCOSoftware/mashling-cli/cli"
)

var optValidate = &cli.OptionInfo{
	Name:      "validate",
	UsageLine: "validate gatewayJson",
	Short:     "validate gateway JSON",
	Long:      "validate gateway JSON",
}

type cmdValidate struct {
	option *cli.OptionInfo
}

func init() {
	CommandRegistry.RegisterCommand(&cmdValidate{option: optValidate})
}

func (c *cmdValidate) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdValidate) AddFlags(fs *flag.FlagSet) {
}

func (c *cmdValidate) Exec(args []string) error {
	var gatewayJson string
	var err error
	var fileName string

	if len(args) < 1 {
		fileName = ""
	} else {
		fileName = args[0]
	}

	if fileName != "" {
		if fgutil.IsRemote(fileName) {

			gatewayJson, err = fgutil.LoadRemoteFile(fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Error loading app file '%s' - %s\n\n", fileName, err.Error())
				os.Exit(2)
			}
		} else {
			gatewayJson, err = fgutil.LoadLocalFile(fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Error loading app file '%s' - %s\n\n", fileName, err.Error())
				os.Exit(2)
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error: Gateway JSON file not specified\n\n")
		cmdUsage(c)
		os.Exit(2)
	}

	//currentDir, err := os.Getwd()

	if err != nil {
		return err
	}

	return ValidateGateway(gatewayJson)
}
