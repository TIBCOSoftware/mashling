// Command, OptionInfo, help and command execution pattern derived from
// github.com/constabulary/gb, released under MIT license
// https://github.com/constabulary/gb/blob/master/LICENSE

package cli

import (
	"bufio"
	"io"
	"os"
	"sync"

	"github.com/TIBCOSoftware/flogo-cli/util"
)

// Tool is a
type Tool struct {
	commandsMu  sync.Mutex
	optionInfo  *OptionInfo
	registry    *CommandRegistry
	TplUsage    string
	TplCmdUsage string
	TplCmdHelp  string
}

// NewTool creates a new tool
func NewTool(optionInfo *OptionInfo) *Tool {
	return &Tool{
		optionInfo:  optionInfo,
		registry:    NewCommandRegistry(),
		TplUsage:    tplToolUsage,
		TplCmdUsage: tplCmdUsage,
		TplCmdHelp:  tplCmdHelp,
	}
}

// OptionInfo implements HasOptionInfo
func (t *Tool) OptionInfo() *OptionInfo {
	return t.optionInfo
}

// CommandRegistry gets the command registry for the tool
func (t *Tool) CommandRegistry() *CommandRegistry {
	return t.registry
}

// Usage prints the usage details of the tool and exits with error
func (t *Tool) Usage() {
	t.PrintUsage(os.Stderr)
	os.Exit(2)
}

var tplToolUsage = `Usage:

    gateway {{.Name}} <command> [arguments]

Commands:
{{range .OptionInfos}}
    {{.Name | printf "%-20s"}} {{.Short}}{{end}}

`

// PrintUsage prints the usage details of the tool
func (t *Tool) PrintUsage(w io.Writer) {
	bw := bufio.NewWriter(w)

	data := struct {
		Name        string
		OptionInfos []*OptionInfo
	}{
		t.optionInfo.Name,
		t.registry.CommandOptionInfos(),
	}

	fgutil.RenderTemplate(bw, t.TplUsage, data)
	bw.Flush()
}

// CmdUsage prints the usage details of the specified Command and
// exits with error
func (t *Tool) CmdUsage(command Command) {
	CmdUsage(t.optionInfo.Name, command)
}

// PrintCmdUsage prints the usage details of the specified Command
func (t *Tool) PrintCmdUsage(w io.Writer, command Command) {

	PrintCmdUsage(w, t.optionInfo.Name, command)
}

// PrintCmdHelp prints the help details of the specified Command
func (t *Tool) PrintCmdHelp(command Command) {

	PrintCmdHelp(t.optionInfo.Name, command)
}
