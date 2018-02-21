/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/TIBCOSoftware/mashling/cli/assets"
	"github.com/TIBCOSoftware/mashling/cli/cli"
	"github.com/TIBCOSoftware/mashling/lib/model"
	mutil "github.com/TIBCOSoftware/mashling/lib/util"
)

var optCreate = &cli.OptionInfo{
	Name:      "create",
	UsageLine: "create AppName",
	Short:     "Create a mashling gateway",
	Long: `Creates a mashling gateway.
   
   Options:
	   -f       	specify the mashling.json to create gateway project from
	   -mlc     	specify the mashling dependency constraints as comma separated value (for example github.com/TIBCOSoftware/flogo-lib@0.0.0,github.com/TIBCOSoftware/flogo-contrib@0.0.0)
	   -vendor  	specify existing vendor directory to copy
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
	CommandRegistry.RegisterCommand(&cmdCreate{option: optCreate, currentDir: getwd})
}

type cmdCreate struct {
	option      *cli.OptionInfo
	constraints string
	fileName    string
	vendorDir   string
	pingport    string
	currentDir  func() (dir string, err error)
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdCreate) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdCreate) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&(c.constraints), "mlc", "", "mashling library constraints")
	fs.StringVar(&(c.fileName), "f", "", "gateway app file")
	fs.StringVar(&(c.vendorDir), "vendor", "", "vendor dir")
	fs.StringVar(&(c.pingport), "pingport", "", "ping port")
}

// Exec implementation of cli.Command.Exec
func (c *cmdCreate) Exec(args []string) error {

	var (
		gatewayJSON string
		gatewayName string
		manifest    io.Reader
	)

	_, err := os.Stat("manifest")
	if err == nil {
		var file *os.File
		file, err = os.Open("manifest")
		if err != nil {
			return err
		}
		defer file.Close()
		manifest = file
	}

	if c.fileName != "" {

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

		if manifest == nil {
			manifest = bytes.NewBuffer(assets.MustAsset("assets/default_manifest"))
		}
	}

	currentDir, err := c.currentDir()

	if err != nil {
		return err
	}

	appDir := filepath.Join(currentDir, gatewayName)

	isValidJSON, err := IsValidGateway(gatewayJSON)

	if !isValidJSON {
		fmt.Print("Mashling creation aborted \n")
		return err
	}

	return CreateMashling(SetupNewProjectEnv(), gatewayJSON, manifest, appDir, gatewayName, c.vendorDir, c.pingport, c.constraints, func() error {
		// Load GB manifest file to extract flogo-lib and mashling repository revisions.
		manifestFile, err := ioutil.ReadFile(filepath.Join(appDir, "vendor", "manifest"))
		if err != nil {
			return err
		}
		var manifestContents GbManifest
		json.Unmarshal(manifestFile, &manifestContents)
		// Extract dependency revisions.
		var flogoLibRev, mashlingRev string
		for _, dep := range manifestContents.Dependencies {
			if flogoLibRev != "" && mashlingRev != "" {
				break
			} else if dep.Repository == "https://github.com/TIBCOSoftware/flogo-lib" && flogoLibRev == "" {
				flogoLibRev = dep.Revision
			} else if dep.Repository == "https://github.com/TIBCOSoftware/mashling" && mashlingRev == "" {
				mashlingRev = dep.Revision
			}
		}
		// Load the main.go file so we can inject extract meta data output.
		gatewayMain, err := ioutil.ReadFile(filepath.Join(appDir, "src", strings.ToLower(gatewayName), "main.go"))
		if err != nil {
			return err
		}
		lines := strings.Split(string(gatewayMain), "\n")
		fileContent := ""
		// Create src payload.
		var extraSrc bytes.Buffer
		// Add the ASCII banner.
		banner, err := assets.Asset("assets/banner.txt")
		if err != nil {
			// Asset was not found.
			return err
		}

		schemaVersion, err := getSchemaVersion(gatewayJSON)
		if err != nil {
			return err
		}

		if strings.Compare(os.Getenv(mutil.Mashling_Ping_Embed_Config_Property), "TRUE") == 0 {
			mashCliOutput := fmt.Sprintf("\n\tmashlingCliRev :=  \"%s\"", MashlingMasterGitRev)
			extraSrc.WriteString(string(mashCliOutput))

			mashCliOutput = fmt.Sprint("\n\tutil.PingDataPntr.MashlingCliRev = mashlingCliRev \n")
			extraSrc.WriteString(string(mashCliOutput))

			mashCliOutput = fmt.Sprintf("\n\tmashlingCliVersion :=  \"%s\"", Version)
			extraSrc.WriteString(string(mashCliOutput))

			mashCliOutput = fmt.Sprint("\n\tutil.PingDataPntr.MashlingCliVersion = mashlingCliVersion \n")
			extraSrc.WriteString(string(mashCliOutput))

			if DisplayLocalChanges {
				mashCliOutput = fmt.Sprintf("\n\tmashlingLocRev :=  \"%s\"", MashlingLocalGitRev)
				extraSrc.WriteString(string(mashCliOutput))

				mashCliOutput = fmt.Sprint("\n\tutil.PingDataPntr.MashlingCliLocalRev = mashlingLocRev \n")
				extraSrc.WriteString(string(mashCliOutput))
			}

			mashCliOutput = fmt.Sprint("\n\tappVersion := app.Version")
			extraSrc.WriteString(string(mashCliOutput))

			mashCliOutput = fmt.Sprint("\n\tutil.PingDataPntr.AppVersion = appVersion \n")
			extraSrc.WriteString(string(mashCliOutput))

			mashCliOutput = fmt.Sprintf("\n\tschemaVersion :=  \"%s\"", schemaVersion)
			extraSrc.WriteString(string(mashCliOutput))

			mashCliOutput = fmt.Sprint("\n\tutil.PingDataPntr.SchemaVersion = schemaVersion \n")
			extraSrc.WriteString(string(mashCliOutput))

			mashCliOutput = fmt.Sprintf("\n\tflogolibRev :=  \"%s\"", flogoLibRev)
			extraSrc.WriteString(string(mashCliOutput))

			mashCliOutput = fmt.Sprint("\n\tutil.PingDataPntr.FlogolibRev = flogolibRev \n")
			extraSrc.WriteString(string(mashCliOutput))

			mashCliOutput = fmt.Sprintf("\n\tmashlingRev :=  \"%s\"", mashlingRev)
			extraSrc.WriteString(string(mashCliOutput))

			mashCliOutput = fmt.Sprint("\n\tutil.PingDataPntr.MashlingRev = mashlingRev \n")
			extraSrc.WriteString(string(mashCliOutput))

			mashCliOutput = fmt.Sprintf("\n\tappDesc := app.Description")
			extraSrc.WriteString(string(mashCliOutput))

			mashCliOutput = fmt.Sprint("\n\tutil.PingDataPntr.AppDescrption = appDesc \n\n")
			extraSrc.WriteString(string(mashCliOutput))

		}

		mashlingCliOutput := fmt.Sprintf("\n\tmashlingTxt :=  \"\\n[mashling] mashling CLI version %s\"", Version)
		extraSrc.WriteString(string(mashlingCliOutput))

		mashlingCliOutput = fmt.Sprintf("\n\tmashlingTxt = mashlingTxt + \"\\n[mashling] mashling CLI revision %s\"", MashlingMasterGitRev)
		extraSrc.WriteString(string(mashlingCliOutput))

		if DisplayLocalChanges {
			mashlingCliOutput = fmt.Sprintf("\n\tmashlingTxt = mashlingTxt + \"\\n[mashling] mashling local revision %s\"", MashlingLocalGitRev)
			extraSrc.WriteString(string(mashlingCliOutput))
		}

		mashlingCliOutput = fmt.Sprintf("\n\tmashlingTxt = mashlingTxt + \"\\n\\n\"")
		extraSrc.WriteString(string(mashlingCliOutput))

		mashlingCliOutput = fmt.Sprintf("\n\tfmt.Printf(\"%%s\\n\", mashlingTxt)\n")
		extraSrc.WriteString(string(mashlingCliOutput))

		bannerOutput := fmt.Sprintf("\tbannerTxt := `%s`\n\tfmt.Printf(\"%%s\\n\", bannerTxt)\n", banner)
		extraSrc.WriteString(string(bannerOutput))

		// Append file version output.
		versionOutput := fmt.Sprintf("\tfmt.Printf(\"[mashling] App Version: %%s\\n\", app.Version)\n")
		extraSrc.WriteString(versionOutput)
		// Append schema version output.
		schemaString := fmt.Sprintf("\tfmt.Printf(\"[mashling] Schema Version: %s\\n\")\n", schemaVersion)
		extraSrc.WriteString(schemaString)
		// Append flogo-lib and mashling revisions
		if flogoLibRev != "" {
			flogoLibString := fmt.Sprintf("\tfmt.Printf(\"[mashling] flogo-lib revision: %s\\n\")\n", flogoLibRev)
			extraSrc.WriteString(flogoLibString)
		}
		if mashlingRev != "" {
			mashlingString := fmt.Sprintf("\tfmt.Printf(\"[mashling] mashling revision: %s\\n\")\n", mashlingRev)
			extraSrc.WriteString(mashlingString)
		}
		// Append app description.
		descriptionOutput := fmt.Sprintf("\tfmt.Printf(\"[mashling] App Description: %%s\\n\", app.Description)\n")
		extraSrc.WriteString(descriptionOutput)
		// Cycle through the file contents, inject source, then rewrite the file.
		for _, line := range lines {
			if strings.Contains(line, "e.Start()") {
				fileContent += extraSrc.String()
			}
			fileContent += line
			fileContent += "\n"
		}
		return ioutil.WriteFile(filepath.Join(appDir, "src", strings.ToLower(gatewayName), "main.go"), []byte(fileContent), 0644)
	})
}

func getwd() (dir string, err error) {
	return os.Getwd()
}
