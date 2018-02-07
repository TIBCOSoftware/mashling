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
	"net"
	"os"
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
	return env.NewGbProjectEnv()
}

func SetupExistingProjectEnv(appDir string) env.Project {

	env := env.NewGbProjectEnv()

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
