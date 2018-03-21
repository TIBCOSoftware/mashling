/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const referenceGatewaySwaggerJSON string = `{
    "host": "localhost",
    "info": {
        "description": "This is the first microgateway app",
        "title": "demo",
        "version": "1.0.0"
    },
    "paths": {
        "/pets/{petId}": {
            "get": {
                "description": "The trigger on 'pets' endpoint",
                "parameters": [
                    {
                        "in": "path",
                        "name": "petId",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "The trigger on 'pets' endpoint"
                    },
                    "default": {
                        "description": "error"
                    }
                },
                "tags": [
                    "rest_trigger"
                ]
            }
        }
    },
    "swagger": "2.0"
}`
const conditionatlGatewaySwaggerJSON string = `{
    "host": "localhost",
    "info": {
        "description": "This is the rest based microgateway app",
        "title": "demoRestGw",
        "version": "1.0.0"
    },
    "paths": {
        "/pets": {
            "put": {
                "description": "Animals rest trigger - PUT animal details",
                "parameters": [],
                "responses": {
                    "200": {
                        "description": "Animals rest trigger - PUT animal details"
                    },
                    "default": {
                        "description": "error"
                    }
                },
                "tags": [
                    "animals_rest_trigger"
                ]
            }
        },
        "/pets/{petId}": {
            "get": {
                "description": "Animals rest trigger - get animal details",
                "parameters": [
                    {
                        "in": "path",
                        "name": "petId",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Animals rest trigger - get animal details"
                    },
                    "default": {
                        "description": "error"
                    }
                },
                "tags": [
                    "get_animals_rest_trigger"
                ]
            }
        }
    },
    "swagger": "2.0"
}`

const emptySwaggerJSON string = `{
    "host": "localhost",
    "info": {
        "description": "This is the kafka based microgateway app",
        "title": "kafka",
        "version": "1.0.0"
    },
    "paths": {},
    "swagger": "2.0"
}`

func TestSwaggerGeneration(t *testing.T) {
	// Samples directory.
	dir, err := filepath.Abs("../samples/")
	assert.NoError(t, err, "Error: Error getting absolute path for samples directory '%s' %v", "../samples/", err)
	samples := make(map[string]string)
	// Map sample file names to expected Swagger 2.0 JSON output.
	samples["reference-gateway/reference-gateway.json"] = referenceGatewaySwaggerJSON
	samples["rest-conditional-gateway/rest-conditional-gateway.json"] = conditionatlGatewaySwaggerJSON
	samples["kafka-reference-gateway/kafka-reference-gateway.json"] = emptySwaggerJSON

	// Compare generated Swagger JSON to expected constants.
	for filename, expected := range samples {
		sample, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, filename))
		assert.NoError(t, err, "Error: Error getting the sample file '%s' %v", filename, err)
		swagger, err := generateSwagger("localhost", "", string(sample))
		assert.NoError(t, err, "Error: Error generating swagger for '%s' %v", filename, err)
		assert.Equal(t, string(swagger), expected)
	}

}
