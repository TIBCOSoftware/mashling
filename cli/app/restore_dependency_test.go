/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"flag"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/TIBCOSoftware/mashling/cli/cli"
)

const restoreGatewayManifest string = `{
	"version": 0,
	"dependencies": [
		{
			"importpath": "github.com/Sirupsen/logrus",
			"repository": "https://github.com/Sirupsen/logrus",
			"revision": "89742aefa4b206dcf400792f3bd35b542998eb3b",
			"branch": "master"
		},
		{
			"importpath": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
			"repository": "https://github.com/TIBCOSoftware/flogo-contrib",
			"revision": "d49b1a9a060e0d6faa52ab41ef46b8f131a23abb",
			"branch": "master",
			"path": "/action/flow"
		},
		{
			"importpath": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
			"repository": "https://github.com/TIBCOSoftware/flogo-contrib",
			"revision": "d49b1a9a060e0d6faa52ab41ef46b8f131a23abb",
			"branch": "master",
			"path": "/activity/log"
		},
		{
			"importpath": "github.com/TIBCOSoftware/flogo-contrib/activity/reply",
			"repository": "https://github.com/TIBCOSoftware/flogo-contrib",
			"revision": "d49b1a9a060e0d6faa52ab41ef46b8f131a23abb",
			"branch": "master",
			"path": "/activity/reply"
		},
		{
			"importpath": "github.com/TIBCOSoftware/flogo-contrib/activity/rest",
			"repository": "https://github.com/TIBCOSoftware/flogo-contrib",
			"revision": "d49b1a9a060e0d6faa52ab41ef46b8f131a23abb",
			"branch": "master",
			"path": "/activity/rest"
		},
		{
			"importpath": "github.com/TIBCOSoftware/flogo-contrib/model/simple",
			"repository": "https://github.com/TIBCOSoftware/flogo-contrib",
			"revision": "d49b1a9a060e0d6faa52ab41ef46b8f131a23abb",
			"branch": "master",
			"path": "/model/simple"
		},
		{
			"importpath": "github.com/TIBCOSoftware/flogo-contrib/trigger/rest",
			"repository": "https://github.com/TIBCOSoftware/flogo-contrib",
			"revision": "d49b1a9a060e0d6faa52ab41ef46b8f131a23abb",
			"branch": "master",
			"path": "/trigger/rest"
		},
		{
			"importpath": "github.com/TIBCOSoftware/flogo-lib",
			"repository": "https://github.com/TIBCOSoftware/flogo-lib",
			"revision": "52e50da7cdbe38eada4b1449992a296c48ae2349",
			"branch": "master"
		},
		{
			"importpath": "github.com/japm/goScript",
			"repository": "https://github.com/japm/goScript",
			"revision": "caab90145b05376536bbab332312f4812123e197",
			"branch": "master"
		},
		{
			"importpath": "github.com/julienschmidt/httprouter",
			"repository": "https://github.com/julienschmidt/httprouter",
			"revision": "975b5c4c7c21c0e3d2764200bf2aa8e34657ae6e",
			"branch": "master"
		},
		{
			"importpath": "golang.org/x/crypto/ssh/terminal",
			"repository": "https://go.googlesource.com/crypto",
			"revision": "541b9d50ad47e36efd8fb423e938e59ff1691f68",
			"branch": "master",
			"path": "/ssh/terminal"
		},
		{
			"importpath": "golang.org/x/sys/unix",
			"repository": "https://go.googlesource.com/sys",
			"revision": "8dbc5d05d6edcc104950cc299a1ce6641235bc86",
			"branch": "master",
			"path": "/unix"
		}
	]
}
`

const restoreDepTestGatewayJSON string = `{
	"mashling_schema": "0.2",
	"gateway": {
		"name": "demo",
		"version": "1.0.0",
		"display_name":"Reference Gateway",
		"display_image":"displayImage.svg",
		"description": "This is the first microgateway app",
		"configurations": [],
		"triggers": [
			{
				"name": "rest_trigger",
				"description": "The trigger on 'pets' endpoint",
				"type": "github.com/TIBCOSoftware/flogo-contrib/trigger/rest",
				"settings": {
					"port": "9096",
					"method": "GET",
					"path": "/pets/:petId"
				}
			}
		],
		"event_handlers": [
			{
				"name": "get_pet_handler",
				"description": "Handle the user access",
				"reference": "github.com/TIBCOSoftware/mashling/lib/flow/flogo.json",
				"params": {
					"uri": "petstore.swagger.io/v2/pet/3"
				}
			}
		],
		"event_links": [
			{
				"triggers": [
					"rest_trigger"
				],
				"dispatches": [
					{
						"handler": "get_pet_handler"
					}
				]
			}
		]
	}
}
`

func TestRestoreDependency(t *testing.T) {
	cmd, _ := CommandRegistry.Command("create")

	mf, err := os.Create("manifest")
	if err != nil {
		t.Error("Unable to create manifest. Restore dependecny test failed.")
	}
	defer mf.Close()
	rf, err := os.Create("restoreDepTestGateway.json")
	if err != nil {
		t.Error("Unable to create gateway json. Restore dependecny test failed.")
	}
	defer rf.Close()
	_, err = io.Copy(mf, strings.NewReader(restoreGatewayManifest))
	_, err = io.Copy(rf, strings.NewReader(restoreDepTestGatewayJSON))

	cmdArgs := []string{"-f", "restoreDepTestGateway.json", "gw"}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	_ = cli.ExecCommand(fs, cmd, cmdArgs)
	defer os.RemoveAll("gw")
	defer os.Remove("restoreDepTestGateway.json")
	defer os.Remove("manifest")
	if _, err := os.Stat("gw/bin/gw"); os.IsNotExist(err) {
		t.Error("Restore dependecny test failed.")
	}
}
