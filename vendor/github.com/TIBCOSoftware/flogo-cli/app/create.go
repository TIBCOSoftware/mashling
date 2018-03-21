package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/TIBCOSoftware/flogo-cli/util"
	"path/filepath"
)

var optCreate = &cli.OptionInfo{
	Name:      "create",
	UsageLine: "create AppName",
	Short:     "create a flogo project",
	Long: `Creates a flogo project.

Options:
    -flv     specify the flogo dependency constraints as comma separated value (for example github.com/TIBCOSoftware/flogo-lib@0.0.0,github.com/TIBCOSoftware/flogo-contrib@0.0.0)
    -f       specify the flogo.json to create project from
    -vendor  specify existing vendor directory to copy

 `,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdCreate{option: optCreate, currentDir: getwd})
}

type cmdCreate struct {
	option      *cli.OptionInfo
	constraints string
	fileName    string
	vendorDir   string
	currentDir  func() (dir string, err error)
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdCreate) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdCreate) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&(c.constraints), "flv", "", "flogo library constraints")
	fs.StringVar(&(c.fileName), "f", "", "flogo app file")
	fs.StringVar(&(c.vendorDir), "vendor", "", "vendor dir")
}

// Exec implementation of cli.Command.Exec
func (c *cmdCreate) Exec(args []string) error {

	var appJson string
	var appName string
	var err error

	if c.fileName != "" {

		if fgutil.IsRemote(c.fileName) {

			appJson, err = fgutil.LoadRemoteFile(c.fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Error loading app file '%s' - %s\n\n", c.fileName, err.Error())
				os.Exit(2)
			}
		} else {
			appJson, err = fgutil.LoadLocalFile(c.fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Error loading app file '%s' - %s\n\n", c.fileName, err.Error())
				os.Exit(2)
			}

			if len(args) != 0 {
				appName = args[0]
			}
		}
	} else {
		if len(args) == 0 {
			fmt.Fprint(os.Stderr, "Error: Application name not specified\n\n")
			cmdUsage(c)
		}

		if len(args) != 1 {
			fmt.Fprint(os.Stderr, "Error: Too many arguments given\n\n")
			cmdUsage(c)
		}

		appName = args[0]
		appJson = tplSimpleApp
	}

	currentDir, err := c.currentDir()

	if err != nil {
		return err
	}

	appDir := filepath.Join(currentDir, appName)

	return CreateApp(SetupNewProjectEnv(), appJson, appDir, appName, c.vendorDir, c.constraints)
}

func getwd() (dir string, err error) {
	return os.Getwd()
}

var tplSimpleApp = `{
  "name": "AppName",
  "type": "flogo:app",
  "version": "0.0.1",
  "description": "My flogo application description",
  "triggers": [
    {
      "id": "my_rest_trigger",
      "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/rest",
      "settings": {
        "port": "9233"
      },
      "handlers": [
        {
          "actionId": "my_simple_flow",
          "settings": {
            "method": "GET",
            "path": "/test"
          }
        }
      ]
    }
  ],
  "actions": [
    {
      "id": "my_simple_flow",
      "name": "my simple flow",
      "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
      "data": {
        "flow": {
          "name": "my simple flow",
          "attributes": [],
          "rootTask": {
            "id": 1,
            "type": 1,
            "tasks": [
              {
                "id": 2,
                "type": 1,
                "activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
                "name": "log",
                "attributes": [
                  {
                    "name": "message",
                    "value": "Simple Log",
                    "type": "string"
                  }
                ]
              }
            ],
            "links": [
            ]
          }
        }
      }
    }
  ]
}`
