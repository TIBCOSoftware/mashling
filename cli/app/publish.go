package app

import (
	"errors"
	"flag"
	"fmt"
	"github.com/TIBCOSoftware/mashling/cli/cli"
	"os"
	"strconv"
)

var optPublish = &cli.OptionInfo{
	Name:      "publish",
	UsageLine: "publish",
	Short:     "Publish to mashery",
	Long: `Publish http triggers to mashery.

Options:
    -f       specify the mashling json
    -k       the api key (required)
    -s       the api secret key (required)
    -u       username (required)
    -p       password (required)
    -portal  the portal (required)
    -uuid    the proxy uuid (required)
 `,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdPublish{option: optPublish})
}

type cmdPublish struct {
	option    *cli.OptionInfo
	fileName  string
	apiKey    string
	apiSecret string
	username  string
	password  string
	uuid      string
	portal    string
	mock      string
	host      string
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdPublish) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdPublish) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&(c.apiKey), "k", "", "api key")
	fs.StringVar(&(c.apiSecret), "s", "", "api secret")
	fs.StringVar(&(c.username), "u", "", "username")
	fs.StringVar(&(c.password), "p", "", "password")
	fs.StringVar(&(c.uuid), "uuid", "", "uuid")
	fs.StringVar(&(c.portal), "portal", "", "portal")
	fs.StringVar(&(c.fileName), "f", "mashling.json", "gateway app file")
	fs.StringVar(&(c.mock), "mock", "false", "mock")
	fs.StringVar(&(c.host), "h", "localhost", "the hostname where this mashling will be deployed (default is localhost)")

}

// Exec implementation of cli.Command.Exec
func (c *cmdPublish) Exec(args []string) error {
	if c.apiKey == "" || c.apiSecret == "" || c.username == "" || c.password == "" ||
		c.uuid == "" || c.portal == "" {
		return errors.New("Error: api key and api secret keys are required")
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Not able read current directory. \n\n")
		return err
	}

	gatewayJSON, _, err := GetGatewayJSON(c.fileName)

	user := ApiUser{c.username, c.password, c.apiKey, c.apiSecret, c.uuid, c.portal}
	b, err := strconv.ParseBool(c.mock)
	if err != nil {
		panic("Invalid option for -mock")
	}
	return PublishToMashery(&user, currentDir, gatewayJSON, c.host, b)
}
