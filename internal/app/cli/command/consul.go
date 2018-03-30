package command

import (
	"log"

	"github.com/TIBCOSoftware/mashling/internal/pkg/consul"
	"github.com/spf13/cobra"
)

func init() {
	consulCommand.Flags().StringVarP(&consulToken, "consulToken", "t", "", "consul agent security token")
	consulCommand.Flags().StringVarP(&consulDefDir, "consulDefDir", "D", "", "service definition folder")
	consulCommand.Flags().StringVarP(&cHost, "host", "H", "", "the hostname where consul is running (e.g. hostip:port)")
	consulCommand.Flags().BoolVarP(&consulRegister, "consulRegister", "r", true, "register services with consul (required -a & -r mutually exclusive)")
	consulCommand.Flags().BoolVarP(&consulDeRegister, "consulDeRegister", "d", false, "de-register services with consul (required -a & -r mutually exclusive)")
	consulCommand.MarkFlagRequired("consulToken")
	publishCommand.AddCommand(consulCommand)
}

var (
	consulToken      string
	consulRegister   bool
	consulDeRegister bool
	consulDefDir     string
	cHost            string
)

var consulCommand = &cobra.Command{
	Use:   "consul",
	Short: "Publishes to Consul",
	Long:  `Publishes the details of the mashling.json configuration file Consul`,
	Run:   consulReg,
}

func consulReg(command *cobra.Command, args []string) {
	err := loadGateway()
	if err != nil {
		log.Fatal(err)
	}
	if !consulRegister && !consulDeRegister {
		log.Fatal("use register or de-register flag")
	}
	if consulRegister && consulDeRegister {
		log.Fatal("cannot use register and de-register together")
	}
	if consulDefDir == "" && cHost == "" {
		log.Fatal("argument missing consul agent address(-h ip:port) is needed")
	}

	if consulRegister {
		consulServiceDefinitions, cErr := gateway.ConsulServiceDefinition()
		if cErr != nil {
			log.Fatal(cErr)
		}
		err = consul.RegisterWithConsul(consulServiceDefinitions, consulToken, consulDefDir, cHost)
	} else {
		consulServiceDefinitions, cErr := gateway.ConsulServiceDefinition()
		if cErr != nil {
			log.Fatal(cErr)
		}
		err = consul.DeregisterFromConsul(consulServiceDefinitions, consulToken, consulDefDir, cHost)
	}
	if err != nil {
		log.Fatal(err)
	}
}
