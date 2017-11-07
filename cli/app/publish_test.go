/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"flag"
	"github.com/TIBCOSoftware/mashling/cli/cli"
	"os"
	"testing"
)

func TestPublishCommand(t *testing.T) {
	cmd, exists := CommandRegistry.Command("publish")

	if !exists {
		t.Error("Publish command should be registered.")
	}

	incompleteArgs := []string{"-u", "username", "-p", "password", "-areaId", "xxxyyy"}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	if err := cli.ExecCommand(fs, cmd, incompleteArgs); err == nil {
		t.Error("All the required switches are not in place.")
	}
}
