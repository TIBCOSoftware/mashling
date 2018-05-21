package command

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/TIBCOSoftware/mashling/internal/pkg/mashery"
	"github.com/spf13/cobra"
)

func init() {
	masheryCommand.Flags().StringVarP(&apiKey, "apiKey", "k", "", "the API key")
	masheryCommand.Flags().StringVarP(&apiSecret, "secretKey", "s", "", "the secret key")
	masheryCommand.Flags().StringVarP(&username, "username", "u", "", "username")
	masheryCommand.Flags().StringVarP(&password, "password", "p", "", "password")
	masheryCommand.Flags().StringVarP(&areaDomain, "areaDomain", "d", "", "the public domain of the Mashery gateway")
	masheryCommand.Flags().StringVarP(&areaID, "areaID", "i", "", "the Mashery area id")
	masheryCommand.Flags().StringVarP(&mHost, "host", "H", "", "the publicly available hostname where this mashling will be deployed (e.g. hostip:port)")
	masheryCommand.Flags().BoolVarP(&iodocs, "iodocs", "I", false, "true to create iodocs")
	masheryCommand.Flags().BoolVarP(&testplan, "testplan", "t", false, "true to create package, plan and test app/key")
	masheryCommand.Flags().BoolVarP(&mock, "mock", "m", false, "true to mock, where it will simply display the transformed swagger doc; false to actually publish to Mashery")
	masheryCommand.Flags().StringVarP(&apiTemplate, "apiTemplate", "T", "", "json file that contains defaults for api/endpoint settings in mashery")
	masheryCommand.MarkFlagRequired("apiKey")
	masheryCommand.MarkFlagRequired("secretKey")
	masheryCommand.MarkFlagRequired("username")
	masheryCommand.MarkFlagRequired("password")
	masheryCommand.MarkFlagRequired("areaDomain")
	masheryCommand.MarkFlagRequired("areaID")
	masheryCommand.MarkFlagRequired("host")
	publishCommand.AddCommand(masheryCommand)
}

var (
	apiKey      string
	apiSecret   string
	username    string
	password    string
	areaID      string
	areaDomain  string
	mock        bool
	mHost       string
	iodocs      bool
	testplan    bool
	apiTemplate string
)

var masheryCommand = &cobra.Command{
	Use:   "mashery",
	Short: "Publishes to Mashery",
	Long:  `Publishes the details of the mashling.json configuration file Mashery`,
	Run:   masheryPublish,
}

func masheryPublish(command *cobra.Command, args []string) {
	err := loadGateway()
	if err != nil {
		log.Fatal(err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	user := mashery.ApiUser{Username: username, Password: password, ApiKey: apiKey, ApiSecretKey: apiSecret, Uuid: areaID, Portal: areaDomain, Noop: false}
	var apiTemplateJSON []byte
	if apiTemplate != "" {
		apiTemplateJSON, err = ioutil.ReadFile(apiTemplate)
		if err != nil {
			log.Fatal(err)
		}
	}
	docs, err := gateway.Swagger(mHost, "")
	if err != nil {
		log.Fatal(err)
	}
	err = mashery.PublishToMashery(&user, pwd, docs, mHost, mock, iodocs, testplan, apiTemplateJSON)
	if err != nil {
		log.Fatal(err)
	}
}
