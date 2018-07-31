package grpcsupport

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	allFIlesFlag bool
)

//GrpcSupportData holds clean up support data
type GrpcSupportData struct {
	ProtoPath string
	AllFiles  bool
}

//CleanSupportFiles cleans all the generated code depends on the proto file name
func CleanSupportFiles(supStruct GrpcSupportData) error {
	assignValues()
	var fileDetails []string

	if len(supStruct.ProtoPath) != 0 {

		path, _ := filepath.Abs(supStruct.ProtoPath)
		_, err := os.Stat(path)
		if err != nil {
			fmt.Println("file path provided is invalid")
			return err
		}
		protoFileName = path[strings.LastIndex(path, string(filepath.Separator))+1:]
		protoFileName = strings.Split(protoFileName, ".")[0]

		err = filepath.Walk(filepath.Join(goPath, "src", grpcGenPath), func(path string, info os.FileInfo, err error) error {
			if !(strings.Compare(info.Name(), "import.go") == 0 || strings.Compare(info.Name(), "grpc") == 0 || strings.Compare(info.Name(), "client") == 0 || strings.Compare(info.Name(), "server") == 0) {
				if (strings.Contains(info.Name(), protoFileName+".") && (strings.Contains(info.Name(), ".pb.go") || strings.Contains(info.Name(), ".grpcservice.go"))) || strings.Compare(info.Name(), protoFileName) == 0 {
					fileDetails = append(fileDetails, path)
				}
			}
			return nil
		})
		if err != nil {
			return err
		}

		return deleteFiles(fileDetails)
	}

	// if proto file not given deleting all the support files
	err := filepath.Walk(filepath.Join(goPath, "src", grpcGenPath), func(path string, info os.FileInfo, err error) error {
		if !(strings.Compare(info.Name(), "import.go") == 0 || strings.Compare(info.Name(), "grpc") == 0 || strings.Compare(info.Name(), "client") == 0 || strings.Compare(info.Name(), "server") == 0) {
			fileDetails = append(fileDetails, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return deleteFiles(fileDetails)
}

func deleteFiles(files []string) error {
	for _, fileName := range files {
		if strings.Compare(filepath.Join(goPath, "src", grpcGenPath), fileName) != 0 {
			os.RemoveAll(fileName)
		}
	}
	return nil
}
