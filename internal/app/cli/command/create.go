//go:generate go run ../generate/stub_generator.go

package command

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/TIBCOSoftware/mashling/internal/app/cli/assets"
	"github.com/TIBCOSoftware/mashling/pkg/files"
	"github.com/TIBCOSoftware/mashling/pkg/strings"
)

// Create builds a custom mashling-gateway project directory populated with
// dependencies listed in the provided Mashling config file.
func Create(name string, deps []string) (err error) {
	if _, err = os.Stat(name); os.IsNotExist(err) {
		err = os.MkdirAll(name, 0755)
		if err != nil {
			return err
		}
	}
	stub, err := assets.Asset("stub.zip")
	if err != nil {
		return err
	}
	err = files.UnpackBytes(stub, name)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	// Run dep to install base dependencies
	// cmd = exec.Command("make", "dep")
	// cmd.Dir = name
	// _, err = cmd.Output()
	// if err != nil {
	// 	log.Println("Failed make dep.")
	// 	return err
	// }
	// Run dep add for all identified new dependencies
	if len(deps) > 0 {
		// Turn deps into a string
		log.Println("Installing missing dependencies.")
		var buffer bytes.Buffer
		buffer.WriteString("NEWDEPS=\"")
		buffer.WriteString(strings.Join(util.UniqueStrings(deps), " "))
		buffer.WriteString("\"")
		cmd = exec.Command("make", "depadd", buffer.String())
		cmd.Dir = name
		_, err = cmd.Output()
		if err != nil {
			return err
		}
	}
	// Run make all target to generate appropriate code and build
	log.Println("Building customized Mashling binary.")
	cmd = exec.Command("make", "all")
	cmd.Dir = name
	_, err = cmd.Output()
	if err != nil {
		return err
	}
	return nil
}
