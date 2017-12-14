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
	Short:     "Version of mashling",
	Long:      "Displays the version of mashling",
}

const notset = "not set"
const refs = "ref:refs/heads/"

//Version is Mashling Version
var Version = "0.3.0"

//MashlingMasterGitRev is mashling git tag
var MashlingMasterGitRev = notset

//MashlingLocalGitRev is mashling git tag
var MashlingLocalGitRev = notset

//FlogoGitRev is flogo-lib git tag
var FlogoGitRev = notset

//SchemaVersion is mashling schema version
var SchemaVersion = GetAllSupportedSchemas()

//GitBranch is git repository checked in
var GitBranch = notset

//GitBranch is git repository checked in
var GITInfo = notset

//DisplayLocalChanges is to check local changes exist flag
var DisplayLocalChanges = false

//DetachedMode is to check git repo in detached or not
var DetachedMode = false

//BranchName to store branch name
var BranchName = ""

func init() {
	CommandRegistry.RegisterCommand(&cmdVersion{option: optVersion})
	if strings.Compare(MashlingMasterGitRev, MashlingLocalGitRev) != 0 {
		DisplayLocalChanges = true
	}
	if strings.Compare(GITInfo, notset) != 0 {
		if strings.Contains(GITInfo, refs) {
			BranchName = GITInfo[len(refs):len(GITInfo)]
		} else if len(GITInfo) == 40 {
			DetachedMode = true
		}
	}
}

type cmdVersion struct {
	option               *cli.OptionInfo
	versionNumber        string
	mashlingMasterGitRev string
	schemaVersion        string
	flogoGitRev          string
	mashlingLocalGitRev  string
	gitBranch            string
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
		c.mashlingMasterGitRev = MashlingMasterGitRev
		c.mashlingLocalGitRev = MashlingLocalGitRev
		c.schemaVersion = SchemaVersion
		c.flogoGitRev = FlogoGitRev
		c.gitBranch = GitBranch

		fmt.Printf(" mashling CLI version %s\n", c.versionNumber)
		fmt.Printf(" supported schema version %s\n", c.schemaVersion)

		if len(c.gitBranch) != 0 {
			fmt.Printf(" git branch %s \n", c.gitBranch)
		} else if DetachedMode {
			fmt.Println(" git is in detached state")
		} else {
			fmt.Printf(" git local branch %s \n", BranchName)
		}

		if len(c.mashlingMasterGitRev) != 0 {
			fmt.Printf(" mashling CLI revision %s\n", c.mashlingMasterGitRev)
		}

		if DisplayLocalChanges {
			fmt.Printf(" mashling local revision %s\n", c.mashlingLocalGitRev)
		}
		fmt.Printf(" flogo-lib revision %s\n", c.flogoGitRev)
	}

	return nil
}
