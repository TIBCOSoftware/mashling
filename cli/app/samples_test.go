/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/stretchr/testify/assert"
)

func TestSampleGateways(t *testing.T) {
	resetDir, err := os.Getwd()
	defer os.Chdir(resetDir)
	assert.NoError(t, err, "Unable to access the current directory %v", err)

	tmp, err := ioutil.TempDir("", "mashlingApps")
	if err != nil {
		t.Error("Temp Directory creation failed")
	}
	defer os.RemoveAll(tmp)

	testDir := filepath.Join(tmp, "sample_mashling_gateways")

	err = os.Mkdir(testDir, 0755)
	assert.NoError(t, err, "Unable to create the tests directory %v", err)

	samplesDir, err := filepath.Abs("../samples")
	assert.NoError(t, err, "Unable to access the samples directory %v", samplesDir)

	//change into the test directory
	os.Chdir(testDir)

	fileList := []string{}

	filepath.Walk(samplesDir, func(fpath string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(info.Name()) == ".json" {
			fileList = append(fileList, fpath)
		}
		return nil
	})

	for index, file := range fileList {
		// Skip secure sample for now until env vars are handled at test time.
		if filepath.Base(file) == "secure-rest-conditional-gateway.json" {
			fmt.Fprintf(os.Stdout, "Skipping: '%v'\n", file)
			continue
		}
		fmt.Fprintf(os.Stdout, "Loading: '%v'\n", file)

		gatewayJson, err := fgutil.LoadLocalFile(file)
		assert.NoError(t, err, "Error: Error loading sample app file '%s' - %s\n\n", file, err)

		currentDir, err := os.Getwd()
		assert.NoError(t, err, "Error: Error getting working dir '%v'", err)

		CopyFile(strings.Replace(file, filepath.Base(file), "Gopkg.lock", 1), filepath.Join(currentDir, "Gopkg.lock"))
		CopyFile(strings.Replace(file, filepath.Base(file), "Gopkg.toml", 1), filepath.Join(currentDir, "Gopkg.toml"))

		gatewayName := "Sample" + strconv.Itoa(index)
		appDir := filepath.Join(currentDir, gatewayName)

		err = CreateMashling(gatewayJson, false, appDir, gatewayName, "9090")
		assert.NoError(t, err, "Error: Error creating mashling app '%v' - %v", gatewayName, err)

		sample := filepath.Join(testDir, gatewayName)
		if _, err := os.Stat(sample); os.IsNotExist(err) {
			fmt.Sprintf("File [%v] generated Samples dir [%v]", file, sample)
		}
		assert.NoError(t, err, "Error: Error getting the sample app dir '%v' %v", sample, err)

	}

	err = os.RemoveAll(testDir)
	assert.NoError(t, err, "Error: Error cleaning up working dir '%v' %v", testDir, err)
}
