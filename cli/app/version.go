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
	"strings"

	"github.com/TIBCOSoftware/mashling/cli/cli"
)

var optVersion = &cli.OptionInfo{
	Name:      "version",
	UsageLine: "version",
	Short:     "version of mashling",
	Long:      "Displays the version of mashling",
}

//Version is Mashling Version
var Version = "not set"

//MashlingGitTag is mashling git tag
var MashlingGitTag = "not set"

//FlogoGitTag is flogo-lib git tag
var FlogoGitTag = "not set"

//MashlingLocalGitTag is mashling local git tag
var MashlingLocalGitTag = "not set"

//SchemaVersion is mashling schema version
var SchemaVersion = GetAllSupportedSchemas()

//GitDiffCheck is to check any local changes made to build mashling cli
var GitDiffCheck = ""

func init() {
	CommandRegistry.RegisterCommand(&cmdVersion{option: optVersion})
}

type cmdVersion struct {
	option              *cli.OptionInfo
	versionNumber       string
	mashlingGitTag      string
	schemaVersion       string
	flogoGitTag         string
	mashlingLocalGitTag string
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
		c.versionNumber = Version
		c.mashlingGitTag = MashlingGitTag
		c.schemaVersion = SchemaVersion
		c.flogoGitTag = FlogoGitTag
		c.mashlingLocalGitTag = MashlingLocalGitTag

		fmt.Printf(" mashling version %s\n", c.versionNumber)
		fmt.Printf(" supported schema version %s\n", c.schemaVersion)
		fmt.Printf(" mashling revision %s\n", c.mashlingGitTag)
		fmt.Printf(" flogo-lib revision %s\n", c.flogoGitTag)

		if strings.Compare(MashlingLocalGitTag, MashlingGitTag) != 0 {
			fmt.Printf(" mashling local revision %s", c.mashlingLocalGitTag)
			if len(GitDiffCheck) != 0 {
				fmt.Print("++")
			}
			fmt.Println("")
		} else {
			if len(GitDiffCheck) != 0 {
				fmt.Print(" mashling cli files changed\n")
			}
		}

	}

	return nil
}
