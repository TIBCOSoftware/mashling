/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/TIBCOSoftware/mashling/lib/model"
)

const (
	registerURI   = "/v1/agent/service/register"
	deRegisterURI = "/v1/agent/service/deregister/"
)

func generateConsulDef(gatewayJSON string) ([]byte, error) {

	var registerPayload []byte
	descriptor, err := model.ParseGatewayDescriptor(gatewayJSON)
	if err != nil {
		return nil, err
	}

	content := map[string]interface{}{
		"Name":    descriptor.Gateway.Name,
		"Address": "127.0.0.1",
		"Port":    9096,
	}

	registerPayload, err = json.MarshalIndent(&content, "", "    ")
	if err != nil {
		return nil, err
	}
	return registerPayload, nil
}

//RegisterWithConsul registers suplied gateway json with consul
func RegisterWithConsul(gatewayJSON string, consulAddress string, consulToken string, consulDefDir string) error {

	consulPaylod, err := generateConsulDef(gatewayJSON)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to generate consul payload \n\n")
		return err
	}

	fullURI := consulAddress + registerURI

	statusCode, err := callConsulService(fullURI, []byte(consulPaylod))

	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("registration failed : status code %v", statusCode)
	}

	return nil
}

//DeregisterFromConsul removes suplied gateway json from consul
func DeregisterFromConsul(gatewayJSON string, consulAddress string, consulToken string, consulDefDir string) error {

	descriptor, err := model.ParseGatewayDescriptor(gatewayJSON)
	if err != nil {
		return err
	}

	fullURI := consulAddress + deRegisterURI + descriptor.Gateway.Name

	statusCode, err := callConsulService(fullURI, []byte(""))

	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("deregistration failed : status code %v", statusCode)
	}

	return nil
}

func callConsulService(uri string, payload []byte) (int, error) {

	client := &http.Client{}
	r, _ := http.NewRequest("PUT", uri, bytes.NewReader([]byte(payload)))
	r.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, err
}
