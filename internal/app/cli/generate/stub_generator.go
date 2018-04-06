// This program generates a stub gateway project with no activities or triggers pre-defined.

// +build ignore

package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/TIBCOSoftware/mashling/pkg/files"
)

func main() {
	_, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	destDir, err := ioutil.TempDir("", "mashling-gateway")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(destDir)

	directories := []string{"cmd/mashling-gateway", "internal/app/gateway", "internal/app/version", "internal/app/assets", "internal/pkg", "pkg", "scripts", "vendor", "dockerfiles", "ext", "lib", "cli"}
	topLevelFiles := []string{"Gopkg.lock", "Gopkg.toml", "Makefile"}

	// Copy directories into temp directory.
	for _, dir := range directories {
		err = files.CopyDir(filepath.Join("..", "..", "..", "..", dir), filepath.Join(destDir, dir))
		if err != nil {
			log.Fatal(err)
		}
	}
	// Copy top level files into temp directory.
	for _, file := range topLevelFiles {
		err = files.CopyFile(filepath.Join("..", "..", "..", "..", file), filepath.Join(destDir, file))
		if err != nil {
			log.Fatal(err)
		}
	}
	// Compress the temporary directory and write the archive to the assets directory.
	err = files.Pack(destDir, filepath.Join("..", "assets", "stub.zip"), true)
	if err != nil {
		log.Fatal(err)
	}
}
