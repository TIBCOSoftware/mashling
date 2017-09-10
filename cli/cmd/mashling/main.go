package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/mashling-cli/app"
	"github.com/TIBCOSoftware/mashling-cli/cli"
)

var (
	fs = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
)

func init() {
	fs.Usage = app.Usage
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "FATAL: "+format+"\n", args...)
	os.Exit(1)
}

func main() {
	args := os.Args
	if len(args) < 2 || args[1] == "-h" {
		app.Usage()
	}
	name := args[1]

	var remainingArgs []string

	cmd, exists := app.CommandRegistry.Command(name)

	if !exists {
		fmt.Fprintf(os.Stderr, "FATAL: unknown command %q\n\n", name)
		app.Usage()
	}
	remainingArgs = args[2:]

	if err := cli.ExecCommand(fs, cmd, remainingArgs); err != nil {
		fatalf("command %q failed: %v", name, err)
	}

	os.Exit(0)
}
