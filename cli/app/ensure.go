package app

import (
	"flag"
	"fmt"
	"os"

	api "github.com/TIBCOSoftware/flogo-cli/app"
	"github.com/TIBCOSoftware/flogo-cli/dep"
	"github.com/TIBCOSoftware/mashling/cli/cli"
)

var optEnsure = &cli.OptionInfo{
	Name:      "ensure",
	UsageLine: "ensure [-add][-update][-no-vendor | -vendor-only][-v]",
	Short:     "Ensure gets a project into a complete, reproducible, and likely compilable state",
	Long: `Ensure gets a project into a complete, reproducible, and likely compilable state:

  Options:
    -add              add new dependencies, or populate Gopkg.toml with constraints for existing dependencies (default: false)
    -no-vendor        update Gopkg.lock (if needed), but do not update vendor/ (default: false)
    -update           update the named dependencies (or all, if none are named) in Gopkg.lock to the latest allowed by Gopkg.toml (default: false)
    -v                enable verbose logging (default: false)
    -vendor-only      populate vendor/ from Gopkg.lock without updating it first (default: false)

`,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdEnsure{option: optEnsure})
}

type cmdEnsure struct {
	option     *cli.OptionInfo
	add        string
	update     bool
	noVendor   bool
	verbose    bool
	vendorOnly bool
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdEnsure) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdEnsure) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&(c.add), "add", "", "add")
	fs.BoolVar(&(c.update), "update", false, "update")
	fs.BoolVar(&(c.noVendor), "no-vendor", false, "no-vendor")
	fs.BoolVar(&(c.verbose), "verbose", false, "verbose")
	fs.BoolVar(&(c.vendorOnly), "vendor-only", false, "vendor-only")
}

// Exec implementation of cli.Command.Exec
func (c *cmdEnsure) Exec(args []string) error {

	if len(args) != 0 {
		fmt.Fprint(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	rootDir, err := os.Getwd()

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to determine working directory\n\n")
		os.Exit(2)
	}

	// Create args
	ensureArgs := []string{}
	if len(c.add) > 0 {
		ensureArgs = append(ensureArgs, "-add", c.add)
	}
	if c.update {
		ensureArgs = append(ensureArgs, "-update")
	}
	if c.verbose {
		ensureArgs = append(ensureArgs, "-v")
	}
	if c.noVendor {
		ensureArgs = append(ensureArgs, "-no-vendor")
	} else if c.vendorOnly {
		ensureArgs = append(ensureArgs, "vendor-only")
	}

	depManager := dep.New(api.SetupExistingProjectEnv(rootDir))

	return Ensure(depManager, ensureArgs...)
}
