package command

import (
	"log"

	"github.com/TIBCOSoftware/mashling/internal/pkg/grpcsupport"
	"github.com/spf13/cobra"
)

func init() {
	generateCommand.Flags().StringVarP(&grpcProtoPath, "protoPath", "p", "", "grpc proto file path")
	grpcCommand.AddCommand(generateCommand)
}

var (
	grpcProtoPath string
)

var generateCommand = &cobra.Command{
	Use:   "generate",
	Short: "Generates Support Files",
	Long:  `Generates Support Files for Mashling Use`,
	Run:   generateFiles,
}

func generateFiles(command *cobra.Command, args []string) {

	if len(grpcProtoPath) == 0 {
		log.Fatal("argument missing proto file path(-p path/to/proto/file) is needed")
	} else {
		err := grpcsupport.GenerateSupportFiles(grpcProtoPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}
