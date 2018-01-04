/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"errors"
	"flag"
	"fmt"
	"github.com/TIBCOSoftware/mashling/cli/cli"
	"os"
	"strconv"
)

var optPublish = &cli.OptionInfo{
	Name:      "publish",
	UsageLine: "publish",
	Short:     "Publish to mashery",
	Long: `Publish http triggers to mashery.

Options:
    -configFile      specify the mashling json
    -apiKey          the api key (required)
    -apiSecret       the api secret key (required)
    -username        username (required)
    -password        password (required)
    -areaDomain      the public domain of the Mashery gateway (required)
    -areaId          the Mashery area id  (required)
    -publicHost      the publicly available hostname where this mashling will be deployed (required)
    -iodocs          true to create iodocs,  (default is false)
    -testplan        true to create package, plan and test app/key,  (default is false)
    -mock            true to mock, where it will simply display the transformed swagger doc; false to actually publish to Mashery (default is false)
    -apitemplate     json file that contains defaults for api/endpoint settings in mashery
 `,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdPublish{option: optPublish})
}

type cmdPublish struct {
	option      *cli.OptionInfo
	fileName    string
	apiKey      string
	apiSecret   string
	username    string
	password    string
	areaId      string
	areaDomain  string
	mock        string
	host        string
	iodocs      string
	testplan    string
	apiTemplate string
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdPublish) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdPublish) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&(c.apiKey), "apiKey", "", "api key")
	fs.StringVar(&(c.apiSecret), "apiSecret", "", "api secret")
	fs.StringVar(&(c.username), "username", "", "username")
	fs.StringVar(&(c.password), "password", "", "password")
	fs.StringVar(&(c.areaId), "areaId", "", "areaId")
	fs.StringVar(&(c.areaDomain), "areaDomain", "", "areaDomain")
	fs.StringVar(&(c.fileName), "configFile", "mashling.json", "gateway app file")
	fs.StringVar(&(c.mock), "mock", "false", "mock")
	fs.StringVar(&(c.iodocs), "iodocs", "false", "iodocs")
	fs.StringVar(&(c.testplan), "testplan", "false", "testplan")
	fs.StringVar(&(c.apiTemplate), "apitemplate", "", "api template file")
	fs.StringVar(&(c.host), "publicHost", "", "the publicly available hostname where this mashling will be deployed")
}

// Exec implementation of cli.Command.Exec
func (c *cmdPublish) Exec(args []string) error {
	if c.apiKey == "" || c.apiSecret == "" || c.username == "" || c.password == "" ||
		c.areaId == "" || c.areaDomain == "" {
		return errors.New("Error: api key, api secret, username, password, areaId and areaDomain are required")
	}

	if c.host == "" {
		return errors.New("Error: host is required")
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Not able read current directory. \n\n")
		return err
	}

	gatewayJSON, _, err := GetGatewayJSON(c.fileName)

	user := ApiUser{c.username, c.password, c.apiKey, c.apiSecret, c.areaId, c.areaDomain, false}

	var apiTemplateJSON string
	if c.apiTemplate != "" {
		apiTemplateJSON, _, err = GetGatewayJSON(c.apiTemplate)
	}

	b, err := strconv.ParseBool(c.mock)
	if err != nil {
		panic("Invalid option for -mock")
	}
	d, err := strconv.ParseBool(c.iodocs)
	if err != nil {
		panic("Invalid option for -iodocs")
	}
	e, err := strconv.ParseBool(c.testplan)
	if err != nil {
		panic("Invalid option for -testplan")
	}
	return PublishToMashery(&user, currentDir, gatewayJSON, c.host, b, d, e, apiTemplateJSON)
}
