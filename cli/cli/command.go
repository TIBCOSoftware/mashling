// Command, OptionInfo, help and command execution pattern derived from
// github.com/constabulary/gb, released under MIT license
// https://github.com/constabulary/gb/blob/master/LICENSE
package cli

import (
	"bufio"
	"flag"
	"io"
	"os"
	"sort"
	"sync"

	"github.com/TIBCOSoftware/flogo-cli/util"
)

// Command represents a command that is executed within gateway project
// Derived from github.com/constabulary/gb
type Command interface {
	HasOptionInfo

	AddFlags(fs *flag.FlagSet)

	Exec(args []string) error
}

// CommandRegistry is a registry for commands
type CommandRegistry struct {
	commandsMu sync.Mutex
	commands   map[string]Command
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{commands: make(map[string]Command)}
}

// RegisterCommand registers a command
func (r *CommandRegistry) RegisterCommand(command Command) {

	r.commandsMu.Lock()
	defer r.commandsMu.Unlock()

	if command == nil {
		panic("CommandRegistry: command cannot be nil")
	}

	commandName := command.OptionInfo().Name

	if _, cmdExists := r.commands[commandName]; cmdExists {
		panic("CommandRegistry: command [" + commandName + "] already registered")
	}

	r.commands[commandName] = command
}

// Command gets the specified command
func (r *CommandRegistry) Command(commandName string) (cmd Command, exists bool) {

	r.commandsMu.Lock()
	defer r.commandsMu.Unlock()

	command, exists := r.commands[commandName]
	return command, exists
}

// Commands gets all the registered commands
func (r *CommandRegistry) Commands() []Command {

	r.commandsMu.Lock()
	defer r.commandsMu.Unlock()

	var cmds []Command
	for _, v := range r.commands {
		cmds = append(cmds, v)
	}

	return cmds
}

// CommandOptionInfos gets the OptionInfos for all registered commands
func (r *CommandRegistry) CommandOptionInfos() []*OptionInfo {

	r.commandsMu.Lock()
	defer r.commandsMu.Unlock()

	//return command options sorted by name
	var sortedKeys []string
	for k := range r.commands {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)
	var options []*OptionInfo
	for _, k := range sortedKeys {
		options = append(options, r.commands[k].OptionInfo())
	}

	return options
}

// ExecCommand executes the specified command
func ExecCommand(fs *flag.FlagSet, cmd Command, args []string) error {

	cmd.AddFlags(fs)

	if err := fs.Parse(args); err != nil {
		fs.Usage()
		os.Exit(1)
	}

	args = fs.Args()

	return cmd.Exec(args)
}

// CmdUsage prints the usage details of the specified Command and
// exits with error
func CmdUsage(toolName string, command Command) {
	PrintCmdUsage(os.Stderr, toolName, command)
	os.Exit(2)
}

// PrintCmdUsage prints the usage details of the specified Command
func PrintCmdUsage(w io.Writer, toolName string, command Command) {
	bw := bufio.NewWriter(w)

	data := struct {
		ToolName     string
		CmdUsageLine string
	}{
		toolName,
		command.OptionInfo().UsageLine,
	}

	fgutil.RenderTemplate(bw, tplCmdUsage, data)
	bw.Flush()
}

var tplCmdUsage = `Usage:

    mashling {{if .ToolName}}{{.ToolName}}{{end}} {{.CmdUsageLine}}

`

// PrintCmdHelp prints the help details of the specified Command
func PrintCmdHelp(toolName string, command Command) {
	bw := bufio.NewWriter(os.Stdout)

	data := struct {
		ToolName     string
		CmdUsageLine string
		CmdLong      string
	}{
		toolName,
		command.OptionInfo().UsageLine,
		command.OptionInfo().Long,
	}

	fgutil.RenderTemplate(bw, tplCmdHelp, data)
	bw.Flush()
}

var tplCmdHelp = `Usage:

    mashling {{if .ToolName}}{{.ToolName}}{{end}} {{.CmdUsageLine}}

{{.CmdLong | trim}}

`
