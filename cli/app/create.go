/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/TIBCOSoftware/mashling/cli/cli"
	"github.com/TIBCOSoftware/mashling/lib/model"
)

var optCreate = &cli.OptionInfo{
	Name:      "create",
	UsageLine: "create AppName",
	Short:     "Create a mashling gateway",
	Long: `Creates a mashling gateway.
	
	Options:
		-f       	specify the mashling.json to create gateway project from
		-pingport	specify the port for ping functionality
	 `,
}

type GbManifest struct {
	Version      int          `json:"version"`
	Dependencies []Dependency `json:"dependencies"`
}

type Dependency struct {
	Importpath string `json:"importpath"`
	Repository string `json:"repository"`
	Revision   string `json:"revision"`
	Branch     string `json:"branch"`
}

func init() {
	CommandRegistry.RegisterCommand(&cmdCreate{option: optCreate})
}

type cmdCreate struct {
	option   *cli.OptionInfo
	fileName string
	pingport string
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdCreate) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdCreate) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&(c.fileName), "f", "", "gateway app file")
	fs.StringVar(&(c.pingport), "pingport", "", "ping port")
}

// Exec implementation of cli.Command.Exec
func (c *cmdCreate) Exec(args []string) error {

	var (
		gatewayJSON    string
		gatewayName    string
		defaultAppFlag bool
		err            error
	)

	if c.fileName != "" {
		defaultAppFlag = false
		if fgutil.IsRemote(c.fileName) {

			gatewayJSON, err = fgutil.LoadRemoteFile(c.fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Error loading app file '%s' - %s\n\n", c.fileName, err.Error())
				os.Exit(2)
			}
		} else {
			gatewayJSON, err = fgutil.LoadLocalFile(c.fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Error loading app file '%s' - %s\n\n", c.fileName, err.Error())
				os.Exit(2)
			}

			if len(args) != 0 {
				gatewayName = args[0]
			}
		}
	} else {
		if len(args) == 0 {
			fmt.Fprint(os.Stderr, "Error: Gateway name not specified\n\n")
			cmdUsage(c)
		}

		if len(args) != 1 {
			fmt.Fprint(os.Stderr, "Error: Too many arguments given\n\n")
			cmdUsage(c)
		}

		gatewayName = args[0]
		mashling, err := model.CreateMashlingSampleModel()
		if err != nil {
			return err
		}
		data, err := json.MarshalIndent(mashling, "", "\t")
		if err != nil {
			return err
		}
		gatewayJSON = string(data)
		defaultAppFlag = true
	}

	currentDir, err := os.Getwd()

	if err != nil {
		return err
	}

	appDir := filepath.Join(currentDir, gatewayName)

	isValidJSON, err := IsValidGateway(gatewayJSON)

	if !isValidJSON {
		fmt.Print("Mashling creation aborted \n")
		return err
	}

	return CreateMashling(SetupNewProjectEnv(), gatewayJSON, defaultAppFlag, appDir, gatewayName, c.pingport)
}
