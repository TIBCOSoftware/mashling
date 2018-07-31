package command

import (
	"log"

	"github.com/TIBCOSoftware/mashling/internal/pkg/grpcsupport"
	"github.com/spf13/cobra"
)

func init() {
	cleanUpCommand.Flags().StringVarP(&protoPath, "protoPath", "p", "", "grpc proto file path")
	cleanUpCommand.Flags().BoolVarP(&allFilesFlag, "all", "a", false, "clean all files flag")
	grpcCommand.AddCommand(cleanUpCommand)
}

var (
	protoPath    string
	allFilesFlag bool
)

var cleanUpCommand = &cobra.Command{
	Use:   "clean",
	Short: "Removes grpc Support Files",
	Long:  `Removes grpc Support Files from grpc trigger location`,
	Run:   cleanFiles,
}

func cleanFiles(command *cobra.Command, args []string) {

	if len(protoPath) == 0 && !allFilesFlag {
		log.Fatal("argument missing proto file path(-p path/to/proto/file) or (-a )is needed")
	} else {
		supStruct := grpcsupport.GrpcSupportData{
			ProtoPath: protoPath,
			AllFiles:  allFilesFlag,
		}
		err := grpcsupport.CleanSupportFiles(supStruct)
		if err != nil {
			log.Fatal(err)
		}
	}
}
