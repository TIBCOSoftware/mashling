package main

import (
	//used to load init func in generated code
	_ "github.com/TIBCOSoftware/mashling/gen/grpc/client"
	_ "github.com/TIBCOSoftware/mashling/gen/grpc/server"

	"github.com/TIBCOSoftware/mashling/internal/app/gateway/command"
)

func main() {
	command.Execute()
}
