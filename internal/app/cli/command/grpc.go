package command

import (
	"github.com/spf13/cobra"
)

func init() {
	cliCommand.AddCommand(grpcCommand)
}

var grpcCommand = &cobra.Command{
	Use:   "grpc",
	Short: "grpc supported platforms",
	Long:  `gRPC supported platform used to generated support files for gRPC`,
}
