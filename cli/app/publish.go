/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"flag"
	"fmt"
	"github.com/TIBCOSoftware/mashling/cli/cli"
	"github.com/go-ini/ini"
	"os"
	"os/user"
	"strconv"
)

type masheryCredFileStruct struct {
	ApiKey     string
	ApiSecret  string
	Username   string
	Password   string
	AreaDomain string
	AreaId     string
	PublicHost string
	IoDocs     bool
	TestPlan   bool
}

var masheryCredFile = ".mashery.conf"

var optPublish = &cli.OptionInfo{
	Name:      "publish",
	UsageLine: "publish",
	Short:     "Publish to mashery",
	Long: `Publish http triggers to mashery. The mashery creds file can be
provided with a -creds switch; if it's not provided, mashling will look for 
<HOME>/.mashery.conf. 

The file should contain:
	ApiKey=xxxyyyzzz
	ApiSecret=aaabbbccc
	Username=someuser
	Password=somepassword
	AreaDomain=somedomain.example.com
	AreaId=xxxyyyzzz
	PublicHost=somewhere.example.com
	IoDocs=false
	TestPlan=false

AreaDomain: the public domain of the Mashery gateway
AreaId:     the Mashery area id
PublicHost: the publicly available hostname where this mashling will be deployed

Options:
    -configFile  specify the mashling json
    -creds       path to Mashery configuration file if the dot file is not used
    -mock        true to mock, where it will simply display the transformed swagger doc; false to actually publish to Mashery (default is false)
    -apitemplate json file that contains defaults for api/endpoint settings in mashery
    -skipVerify  true to skip SSL verification; default is false
 `,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdPublish{option: optPublish})
}

type cmdPublish struct {
	option           *cli.OptionInfo
	fileName         string
	masheryCredsFile string
	mock             string
	apiTemplate      string
	skipVerify       string
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdPublish) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdPublish) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&(c.fileName), "configFile", "mashling.json", "gateway app file")
	fs.StringVar(&(c.masheryCredsFile), "creds", "", "mashery creds file")
	fs.StringVar(&(c.mock), "mock", "false", "mock")
	fs.StringVar(&(c.apiTemplate), "apitemplate", "", "api template file")
	fs.StringVar(&(c.skipVerify), "skipVerify", "false", "skip SSL verification")
}

// parseConfigFile parse file with mashery configuration
func parseConfigFile(file string) (*ini.File, error) {
	cfg, err := ini.InsensitiveLoad(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		panic("Not able load mashery configuration file " + file)
	}

	return cfg, err
}

// Exec implementation of cli.Command.Exec
func (c *cmdPublish) Exec(args []string) error {
	currentUser, _ := user.Current()

	var cfg *ini.File
	var err error
	if c.masheryCredsFile == "" {
		cfg, err = parseConfigFile(currentUser.HomeDir + "/" + masheryCredFile)
	} else {
		cfg, err = parseConfigFile(c.masheryCredsFile)
	}

	mashery := new(masheryCredFileStruct)
	err = cfg.MapTo(mashery)

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Not able read current directory. \n\n")
		return err
	}

	gatewayJSON, _, err := GetGatewayJSON(c.fileName)

	skipVerify, err := strconv.ParseBool(c.skipVerify)
	if err != nil {
		panic("Invalid option for -skipVerify")
	}

	user := ApiUser{mashery.Username, mashery.Password, mashery.ApiKey, mashery.ApiSecret, mashery.AreaId, mashery.AreaDomain, false, skipVerify}

	var apiTemplateJSON string
	if c.apiTemplate != "" {
		apiTemplateJSON, _, err = GetGatewayJSON(c.apiTemplate)
	}

	b, err := strconv.ParseBool(c.mock)
	if err != nil {
		panic("Invalid option for -mock")
	}

	return PublishToMashery(&user, currentDir, gatewayJSON, mashery.PublicHost, b, mashery.IoDocs, mashery.TestPlan, apiTemplateJSON)
}
