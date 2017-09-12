package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
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
    -f       specify the mashling json
    -h       the hostname where this mashling will be deployed
 `,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdSwagger{option: optSwagger})
}

type cmdSwagger struct {
	option   *cli.OptionInfo
	fileName string
	host     string
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdSwagger) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdSwagger) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&(c.fileName), "f", "", "filename")
	fs.StringVar(&(c.host), "h", "", "host")
}

// Exec implementation of cli.Command.Exec
func (c *cmdSwagger) Exec(args []string) error {

	if c.host == "" {
		fmt.Fprint(os.Stderr, "Error: host is required\n\n")
		os.Exit(2)
	}

	_, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Not able read current directory. \n\n")
		return err
	}

	gatewayJSON, _, err := GetGatewayJSON(c.fileName)

	docs, err := generate_swagger(c.host, gatewayJSON)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s\n", string(docs))
	return nil
}

type swaggerEndpoint struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Method      string `json:"method"`
}

func generate_swagger(host string, gatewayJSON string) ([]byte, error) {
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
				parameters := swagger_parameters(endpoint.Path, begin_delim, end_delim)
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
				paths[endpoint.Path] = path
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

func swagger_parameters(path string, begin_delim rune, end_delim rune) []interface{} {
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
			parameter := map[string]interface{}{
				"name":     key.String(),
				"in":       "path",
				"required": true,
				"type":     "string",
			}
			parameters = append(parameters, parameter)
		}
	}
	return parameters
}
