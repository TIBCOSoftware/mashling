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
	"io/ioutil"
	"github.com/stretchr/testify/assert"
)

func TestPublishWithInvalidFile(t *testing.T) {
	defer func() {
        if r := recover(); r == nil {
            t.Errorf("Should have paniced.")
        }
    }()

	cmd, exists := CommandRegistry.Command("publish")

	if !exists {
		t.Error("Publish command should be registered.")
	}

	incompleteArgs := []string{"-creds", "fake_file.conf"}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	if err := cli.ExecCommand(fs, cmd, incompleteArgs); err == nil {
		t.Error("The mashery creds should not exist.")
	}
}

const dataPlain string = `
ApiKey=xxxyyyzzz
ApiSecret=aaabbbccc
Username=someuser
Password=somepassword
AreaDomain=somedomain.example.com
AreaId=xxxyyyzzz
PublicHost=somewhere.example.com
IoDocs=true
TestPlan=true
`

func TestLoadDotFile(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(dataPlain)); err != nil {
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(err)
	}

	path := file.Name()

	cfg, err := parseConfigFile(path)
	if err != nil {
		t.Error("The mashery creds file should load.")
	}

	mashery := new(masheryCredFileStruct)
	err = cfg.MapTo(mashery)

	assert.Equal(t, "xxxyyyzzz", mashery.ApiKey)
	assert.Equal(t, "aaabbbccc", mashery.ApiSecret)
	assert.Equal(t, "someuser", mashery.Username)
	assert.Equal(t, "somepassword", mashery.Password)
	assert.Equal(t, "somedomain.example.com", mashery.AreaDomain)
	assert.Equal(t, "xxxyyyzzz", mashery.AreaId)
	assert.Equal(t, true, mashery.IoDocs)
	assert.Equal(t, "somewhere.example.com", mashery.PublicHost)
	assert.Equal(t, true, mashery.TestPlan)
}
