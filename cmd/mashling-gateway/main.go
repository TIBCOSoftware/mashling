package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/TIBCOSoftware/mashling/internal/app/gateway/assets"
	"github.com/TIBCOSoftware/mashling/internal/app/gateway/flogo"
	model "github.com/TIBCOSoftware/mashling/internal/pkg/model"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/cache"
	gwerrors "github.com/TIBCOSoftware/mashling/internal/pkg/model/errors"
	"github.com/fsnotify/fsnotify"
)

var (
	flogoLibRevision = "52e50da7cdbe38eada4b1449992a296c48ae2349"
	mashlingRevision = "8cc2d8fe27b7252afb5744156096d2d77f792f5c"
	// Version is grabbed from the git repo at build time.
	Version = "dev"
	// BuildDate is set with the current timestamp at build time.
	BuildDate          = "unknown"
	gateway            model.Gateway
	watcher            *fsnotify.Watcher
	config             string
	envVarName         string
	loadFromEnv        bool
	dev                bool
	configCache        string
	configCacheEnabled bool
)

func init() {
	flag.StringVar(&config, "config", "mashling.json", "mashling gateway configuration")
	flag.StringVar(&envVarName, "env-var-name", "MASHLING_CONFIG", "name of the environment variable that contain sthe base64 encoded mashling gateway configuration")
	flag.BoolVar(&loadFromEnv, "load-from-env", false, "load the mashling gateway configuration from an environment variable")
	flag.BoolVar(&dev, "dev", false, "run mashling in dev mode")
	flag.StringVar(&configCache, "config-cache", ".cache", "location of the configuration artifacts cache")
	flag.BoolVar(&configCacheEnabled, "config-cache-enabled", true, "cache post-processed configuration artifacts locally")
}

func main() {
	// Output friendly Mashling mascot with some build details.
	bannerTxt, err := assets.Asset("banner.txt")
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	log.Println("\n", string(bannerTxt))
	log.Println("[mashling] Gateway Version: ", Version)
	log.Println("[mashling] Build Date: ", BuildDate)
	flag.Parse()

	// Setup configuration artifacts cache.
	if configCacheEnabled {
		err = cache.Initialize("file", configCache)
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}
	}

	// Load the appropriate Gateway instance from the configuration.
	if loadFromEnv {
		// Loading base64 encoded configuration from the env.
		gateway, err = model.LoadFromEnv(envVarName)
	} else {
		// Loading the configuration JSON from the specified file.
		gateway, err = model.LoadFromFile(config)
	}
	if err != nil {
		log.Println(err)
		// Attempt to give insight into known  potentials errors.
		if gateway == nil {
			os.Exit(1)
		}
		for _, errd := range gateway.Errors() {
			switch e := errd.(type) {
			case *gwerrors.MissingDependency:
				log.Println("Missing dependencies found: ", strings.Join(e.MissingDependencies, " "))
			default:
				fmt.Printf("Do not know how to handle error type %T!\n", e)
				os.Exit(1)
			}
		}
		// For now we still exit no matter what. The CLI handles missing errors.
		os.Exit(1)
	}

	if dev && !loadFromEnv {
		// Configuring file watcher to automatically reload configuration file contents in dev mode.
		watcher, err = fsnotify.NewWatcher()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer watcher.Close()

		go handleFileNotifications()

		err = watcher.Add(config)
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}
	}

	log.Println("[mashling] Schema Version: ", gateway.Version())
	log.Println("[mashling] App Version: ", gateway.AppVersion())
	log.Println("[mashling] flogo-lib revision: ", flogoLibRevision)
	log.Println("[mashling] mashling revision: ", mashlingRevision)
	log.Println("[mashling] App Description: ", gateway.Description())

	// Startup the configured gateway instance.
	gateway.Init()
	gateway.Start()

	exitChan := setupSignalHandling()

	code := <-exitChan

	// Try to gracefully shutdown
	gateway.Stop()

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
				exitChan <- 1
			}
		}
	}()

	return exitChan
}

// handleFileNotifications intercepts file system notifications and triggers a reload of the configured gateway if the file has changed.
func handleFileNotifications() {
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("[mashling] modified configuration file:", event.Name)
				reloadGatewayFromConfigurationFile()
			}
		case ferr := <-watcher.Errors:
			log.Println("[mashling] error:", ferr)
		}
	}
}

// reloadGatewayFromConfigurationFile attempts to shutdown the current gateway instance and then re-start the gateway with the new configuration loaded.
func reloadGatewayFromConfigurationFile() {
	err := gateway.Stop()
	if err != nil {
		log.Println("[mashling] error stopping gateway:", err)
	}
	flogo.ResetGlobalContext()
	gateway, err = model.LoadFromFile(config)
	if err != nil {
		log.Println("[mashling] error re-loading gateway from file:", err)
	}
	err = gateway.Init()
	if err != nil {
		log.Println("[mashling] error re-initializing gateway:", err)
	}
	err = gateway.Start()
	if err != nil {
		log.Println("[mashling] error re-starting gateway:", err)
	}
}
