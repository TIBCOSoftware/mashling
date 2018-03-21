/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/TIBCOSoftware/mashling/cli/cli"
	"github.com/TIBCOSoftware/mashling/lib/model"
)

var optSwagger = &cli.OptionInfo{
	Name:      "swagger",
	UsageLine: "swagger",
	Short:     "Generate Swagger docs",
	Long: `Generate a Swagger doc representation of HTTP triggers.

Options:
-f	specify the mashling json
-h	the hostname where this mashling will be deployed (default is localhost)
-t	the trigger name to target (default is all)
-o 	the output file to write the swagger.json to (default is stdout)
 `,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdSwagger{option: optSwagger})
}

type cmdSwagger struct {
	option     *cli.OptionInfo
	fileName   string
	host       string
	trigger    string
	outputFile string
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdSwagger) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdSwagger) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&(c.fileName), "f", "mashling.json", "filename")
	fs.StringVar(&(c.host), "h", "localhost", "host")
	fs.StringVar(&(c.trigger), "t", "", "trigger")
	fs.StringVar(&(c.outputFile), "o", "", "output file")
}

// Exec implementation of cli.Command.Exec
func (c *cmdSwagger) Exec(args []string) error {
	if c.host == "" {
		fmt.Fprint(os.Stderr, "Error: host is required. \n\n")
		os.Exit(2)
	}

	_, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Not able read current directory. \n\n")
		return err
	}

	gatewayJSON, _, err := GetGatewayJSON(c.fileName)

	docs, err := generateSwagger(c.host, c.trigger, gatewayJSON)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to generate Swagger representation. \n\n")
		return err
	}
	if c.outputFile == "" {
		fmt.Fprintf(os.Stdout, "%s\n", string(docs))
	} else {
		err := ioutil.WriteFile(c.outputFile, docs, 0644)
		if err != nil {
			fmt.Fprint(os.Stderr, "Error: Not able write Swagger to output file. \n\n")
			return err
		}
	}
	return nil
}

type swaggerEndpoint struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Method      string `json:"method"`
}

func generateSwagger(host string, triggerName string, gatewayJSON string) ([]byte, error) {
	var docs []byte = nil
	descriptor, err := model.ParseGatewayDescriptor(gatewayJSON)
	if err != nil {
		return nil, err
	}
	if descriptor.Gateway.Triggers != nil {
		paths := map[string]interface{}{}
		swagger := map[string]interface{}{
			"swagger": "2.0",
			"info": map[string]interface{}{
				"version":     descriptor.Gateway.Version,
				"title":       descriptor.Gateway.Name,
				"description": descriptor.Gateway.Description,
			},
			"host":  host,
			"paths": paths,
		}

		for _, trigger := range descriptor.Gateway.Triggers {
			if triggerName == "" || triggerName == trigger.Name {
				if trigger.Type == "github.com/TIBCOSoftware/flogo-contrib/trigger/rest" || trigger.Type == "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger" {
					var endpoint swaggerEndpoint
					endpoint.Name = trigger.Name
					endpoint.Description = trigger.Description
					err := json.Unmarshal(trigger.Settings, &endpoint)
					if err != nil {
						return nil, err
					}
					path := map[string]interface{}{}
					var begin_delim, end_delim rune
					switch trigger.Type {
					case "github.com/TIBCOSoftware/flogo-contrib/trigger/rest":
						begin_delim = ':'
						end_delim = '/'
					case "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger":
						begin_delim = '{'
						end_delim = '}'
					default:
						begin_delim = '{'
						end_delim = '}'
					}
					parameters, scrubbedPath := swaggerParametersExtractor(endpoint.Path, begin_delim, end_delim)
					ok := map[string]interface{}{
						"description": endpoint.Description,
					}

					path[strings.ToLower(endpoint.Method)] = map[string]interface{}{
						"description": endpoint.Description,
						"tags":        []interface{}{endpoint.Name},
						"parameters":  parameters,
						"responses": map[string]interface{}{
							"200": ok,
							"default": map[string]interface{}{
								"description": "error",
							},
						},
					}
					paths[scrubbedPath] = path
				}
			}
		}
		docs, err = json.MarshalIndent(&swagger, "", "    ")
		if err != nil {
			return nil, err
		}
	} else {
		err = errors.New("No triggers defined.")
	}

	return docs, err
}

func swaggerParametersExtractor(path string, begin_delim rune, end_delim rune) ([]interface{}, string) {
	parameters := []interface{}{}
	routePath := []rune(path)
	for i := 0; i < len(routePath); i++ {
		if routePath[i] == begin_delim {
			key := bytes.Buffer{}
			for i++; i < len(routePath) && routePath[i] != end_delim; i++ {
				if routePath[i] != ' ' && routePath[i] != '\t' {
					key.WriteRune(routePath[i])
				}
			}
			if begin_delim == ':' {
				path = strings.Replace(path, fmt.Sprintf(":%s", key.String()), fmt.Sprintf("{%s}", key.String()), 1)
			}
			parameter := map[string]interface{}{
				"name":     key.String(),
				"in":       "path",
				"required": true,
				"type":     "string",
			}
			parameters = append(parameters, parameter)
		}
	}
	return parameters, path
}
