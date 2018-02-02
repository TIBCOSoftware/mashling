package app

import (
	"os"
	"strings"

	"path"

	api "github.com/TIBCOSoftware/flogo-cli/app"
	"github.com/TIBCOSoftware/flogo-cli/util"
)

const (
	fileDescriptor    string = "flogo.json"
	fileMainGo        string = "main.go"
	fileImportsGo     string = "imports.go"
	fileEmbeddedAppGo string = "embeddedapp.go"
	makeFile          string = "Makefile"
	fileShimGo        string = "shim.go"
	fileShimSupportGo string = "shim_support.go"

	dirShim      string = "shim"
	pathFlogoLib string = "github.com/TIBCOSoftware/flogo-lib"
)

func CreateMainGoFile(codeSourcePath string, flogoJSON string) {

	data := struct {
		FlogoJSON string
	}{
		flogoJSON,
	}

	f, _ := os.Create(path.Join(codeSourcePath, fileMainGo))
	fgutil.RenderTemplate(f, tplNewMainGoFile, &data)
	f.Close()
}

func removeMainGoFile(codeSourcePath string) {
	os.Remove(path.Join(codeSourcePath, fileMainGo))
}

var tplNewMainGoFile = `package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	` + getImports() + `
)

var (
	cp app.ConfigProvider
)

func main() {

	if cp == nil {
		// Use default config provider
		cp = app.DefaultConfigProvider()
	}

	app, err := cp.GetApp()
	if err != nil {
        	fmt.Println(err.Error())
        	os.Exit(1)
    	}

    	e, err := engine.New(app)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	e.Start()

	exitChan := setupSignalHandling()

	code := <-exitChan

	e.Stop()

	os.Exit(code)
}

func setupSignalHandling() chan int {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exitChan := make(chan int)
	go func() {
		for {
			s := <-signalChan
			switch s {
			// kill -SIGHUP
			case syscall.SIGHUP:
				exitChan <- 0
			// kill -SIGINT/Ctrl+c
			case syscall.SIGINT:
				exitChan <- 0
			// kill -SIGTERM
			case syscall.SIGTERM:
				exitChan <- 0
			// kill -SIGQUIT
			case syscall.SIGQUIT:
				exitChan <- 0
			default:
				logger.Debug("Unknown signal.")
				exitChan <- 1
			}
		}
	}()

	return exitChan
}
`

func getImports() string {
	if strings.Compare(os.Getenv("PING_ENABLE"), "FALSE") != 0 {
		return "\"github.com/nareshkumarthota/sampleflows/pingreply\""
	}
	return ""
}

func CreateImportsGoFile(codeSourcePath string, deps []*api.Dependency) error {
	f, err := os.Create(path.Join(codeSourcePath, fileImportsGo))

	if err != nil {
		return err
	}

	fgutil.RenderTemplate(f, tplNewImportsGoFile, deps)
	f.Close()

	return nil
}

var tplNewImportsGoFile = `package main

import (

{{range $i, $dep := .}}	_ "{{ $dep.Ref }}"
{{end}}
)
`

func createEmbeddedAppGoFile(codeSourcePath string, flogoJSON string) {

	data := struct {
		FlogoJSON string
	}{
		flogoJSON,
	}

	f, _ := os.Create(path.Join(codeSourcePath, fileEmbeddedAppGo))
	fgutil.RenderTemplate(f, tplEmbeddedAppGoFile, &data)
	f.Close()
}

func removeEmbeddedAppGoFile(codeSourcePath string) {
	os.Remove(path.Join(codeSourcePath, fileEmbeddedAppGo))
}

var tplEmbeddedAppGoFile = `// Do not change this file, it has been generated using flogo-cli
// If you change it and rebuild the application your changes might get lost
package main

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/app"
)

// embedded flogo app descriptor file
const flogoJSON string = ` + "`{{.FlogoJSON}}`" + `

func init () {
	cp = EmbeddedProvider()
}

// embeddedConfigProvider implementation of ConfigProvider
type embeddedProvider struct {
}

//EmbeddedProvider returns an app config from a compiled json file
func EmbeddedProvider() (app.ConfigProvider){
	return &embeddedProvider{}
}

// GetApp returns the app configuration
func (d *embeddedProvider) GetApp() (*app.Config, error){

	app := &app.Config{}
	err := json.Unmarshal([]byte(flogoJSON), app)
	if err != nil {
		return nil, err
	}
	return app, nil
}
`

func createShimSupportGoFile(codeSourcePath string, flogoJSON string, embeddedConfig bool) {

	configJson := ""

	if embeddedConfig {
		configJson = flogoJSON
	}

	data := struct {
		FlogoJSON string
	}{
		configJson,
	}

	f, _ := os.Create(path.Join(codeSourcePath, fileShimSupportGo))
	fgutil.RenderTemplate(f, tplShimSupportGoFile, &data)
	f.Close()
}

func removeShimGoFiles(codeSourcePath string) {
	os.Remove(path.Join(codeSourcePath, fileShimGo))
	os.Remove(path.Join(codeSourcePath, fileShimSupportGo))
}

var tplShimSupportGoFile = `// Do not change this file, it has been generated using flogo-cli
// If you change it and rebuild the application your changes might get lost
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/config"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/flogo-lib/logger"

)

// embedded flogo app descriptor file
const flogoJSON string = ` + "`{{.FlogoJSON}}`" + `

func init() {
	config.SetDefaultLogLevel("ERROR")
	logger.SetLogLevel(logger.ErrorLevel)

	var cp app.ConfigProvider

	if flogoJSON != "" {
		cp = EmbeddedProvider()
	} else {
		cp = app.DefaultConfigProvider()
	}

	app, err := cp.GetApp()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	e, err := engine.New(app)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	e.Init(true)
}

// embeddedConfigProvider implementation of ConfigProvider
type embeddedProvider struct {
}

//EmbeddedProvider returns an app config from a compiled json file
func EmbeddedProvider() (app.ConfigProvider){
	return &embeddedProvider{}
}

// GetApp returns the app configuration
func (d *embeddedProvider) GetApp() (*app.Config, error){

	appCfg := &app.Config{}
	err := json.Unmarshal([]byte(flogoJSON), appCfg)
	if err != nil {
		return nil, err
	}
	return appCfg, nil
}
`
