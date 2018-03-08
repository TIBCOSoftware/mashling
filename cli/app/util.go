/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/TIBCOSoftware/mashling/cli/cli"
	"github.com/TIBCOSoftware/mashling/cli/env"
	"github.com/TIBCOSoftware/mashling/lib/types"
	"github.com/TIBCOSoftware/mashling/lib/util"
)

var (
	CommandRegistry = cli.NewCommandRegistry()
)

func SetupNewProjectEnv() env.Project {
	return env.NewMashlingProject()
}

func SetupExistingProjectEnv(appDir string) env.Project {

	env := env.NewMashlingProject()

	if err := env.Init(appDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing mashling app project: %s\n\n", err.Error())
		os.Exit(2)
	}

	if err := env.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "Error opening mashling app project: %s\n\n", err.Error())
		os.Exit(2)
	}

	return env
}

func GetGatewayJSON(fileName string) (string, string, error) {
	var gatewayJson string
	var gatewayName string
	var err error

	if fgutil.IsRemote(fileName) {

		gatewayJson, err = fgutil.LoadRemoteFile(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Error loading app file '%s' - %s\n\n", fileName, err.Error())
			os.Exit(2)
		}
	} else {
		gatewayJson, err = fgutil.LoadLocalFile(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Error loading app file '%s' - %s\n\n", fileName, err.Error())
			os.Exit(2)
		}
	}

	return gatewayJson, gatewayName, err
}

func splitVersion(t string) (path string, version string) {

	idx := strings.LastIndex(t, "@")

	version = ""
	path = t

	if idx > -1 {
		v := t[idx+1:]

		if isValidVersion(v) {
			version = v
			path = t[0:idx]
		}
	}

	return path, version
}

//todo validate that "s" a valid semver
func isValidVersion(s string) bool {

	if s == "" {
		//assume latest version
		return true
	}

	if s[0] == 'v' && len(s) > 1 && isNumeric(string(s[1])) {
		return true
	}

	if isNumeric(string(s[0])) {
		return true
	}

	return false
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func delayMilli(amount int) {
	time.Sleep(time.Duration(amount) * time.Millisecond)
}

func shortDelay() {
	delayMilli(500)
}

func Usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func cmdUsage(command cli.Command) {
	cli.CmdUsage("", command)
}

func printUsage(w io.Writer) {
	bw := bufio.NewWriter(w)

	options := CommandRegistry.CommandOptionInfos()
	options = append(options, cli.GetToolOptionInfos()...)

	fgutil.RenderTemplate(bw, usageTpl, options)
	bw.Flush()
}

var usageTpl = `Usage:

    mashling <command> [arguments]

Commands:
{{range .}}
    {{.Name | printf "%-12s"}} {{.Short}}{{end}}
`

// getLocalIP gets the public ip address of the system
func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

/*
CreateMashlingPingModel creates rest based ping mashling json used for ping functionality
*/
func CreateMashlingPingModel(pingPort string) (types.Microgateway, error) {

	microGateway := types.Microgateway{
		MashlingSchema: "0.2",
		Gateway: types.Gateway{
			Name:         "GatewayPingApp",
			Version:      "1.0.0",
			DisplayName:  "Gateway Ping Application",
			DisplayImage: "GatewayPingIcon.jpg",
			Description:  "This is the first microgateway ping app",
			Configurations: []types.Config{
				{
					Name:        "ping_config",
					Type:        "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
					Description: "The trigger for ping functionality",
					Settings: json.RawMessage(`{
						"port": "` + pingPort + `"
						}`),
				},
			},
			Triggers: []types.Trigger{
				{
					Name:        util.Mashling_Ping_Trigger_Name,
					Description: "The trigger for ping functionality",
					Type:        "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
					Settings: json.RawMessage(`{
						"config": "${configurations.ping_config}",
						"method": "GET",
						"path": "/ping/",
						"optimize": "true"
					}`),
				}, {
					Name:        util.Mashling_Ping_Detail_Trigger_Name,
					Description: "The trigger for detailed ping functionality",
					Type:        "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
					Settings: json.RawMessage(`{
						"config": "${configurations.ping_config}",
						"method": "GET",
						"path": "/ping/details/",
						"optimize": "true"
					}`),
				},
			},
			EventHandlers: []types.EventHandler{
				{
					Name:        "ping_handler",
					Description: "Handle Ping get call",
					Reference:   "github.com/TIBCOSoftware/mashling/lib/flow/pingflow.json",
				},
				{
					Name:        "ping_handler_detail",
					Description: "Handle Ping detailed get call",
					Reference:   "github.com/TIBCOSoftware/mashling/lib/flow/pingflowdetailed.json",
				},
			},
			EventLinks: []types.EventLink{
				{
					Triggers: []string{
						util.Mashling_Ping_Trigger_Name,
					},
					Dispatches: []types.Dispatch{
						{
							Path: types.Path{
								Handler: "ping_handler",
							},
						},
					},
				},
				{
					Triggers: []string{
						util.Mashling_Ping_Detail_Trigger_Name,
					},
					Dispatches: []types.Dispatch{
						{
							Path: types.Path{
								Handler: "ping_handler_detail",
							},
						},
					},
				},
			},
		},
	}

	return microGateway, nil
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}
