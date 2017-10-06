/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/mashling/cli/cli"
)

var optVersion = &cli.OptionInfo{
	Name:      "version",
	UsageLine: "version",
	Short:     "version of mashling",
	Long:      "Displays the version of mashling",
}

//Mashling version
const version = "0.2.0"

func init() {
	CommandRegistry.RegisterCommand(&cmdVersion{option: optVersion})
}

type cmdVersion struct {
	option        *cli.OptionInfo
	versionNumber string
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdVersion) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdVersion) AddFlags(fs *flag.FlagSet) {
	// no flags
}

// Exec implementation of cli.Command.Exec
func (c *cmdVersion) Exec(args []string) error {

	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "usage: mashling version \n\nToo many arguments given.\n")
		os.Exit(2)
	} else {
		c.versionNumber = version
		fmt.Printf("mashling version %s\n", c.versionNumber)
	}

	return nil
}
