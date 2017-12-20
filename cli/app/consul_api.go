/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"encoding/json"

	"github.com/TIBCOSoftware/mashling/lib/model"
)

func generateRegisterPayload(gatewayJSON string) ([]byte, error) {

	var registerPayload []byte = nil
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
