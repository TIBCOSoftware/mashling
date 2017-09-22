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

	incompleteArgs := []string{"-u", "username", "-p", "password", "-uuid", "xxxyyy"}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	if err := cli.ExecCommand(fs, cmd, incompleteArgs); err == nil {
		t.Error("All the required switches are not in place.")
	}
}
