package command

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/TIBCOSoftware/mashling/internal/app/assets"
	"github.com/TIBCOSoftware/mashling/internal/app/gateway/flogo"
	"github.com/TIBCOSoftware/mashling/internal/app/version"
	"github.com/TIBCOSoftware/mashling/internal/pkg/logger"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/cache"
	gwerrors "github.com/TIBCOSoftware/mashling/internal/pkg/model/errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

func init() {
	gatewayCommand.PersistentFlags().StringVarP(&config, "config", "c", "mashling.json", "mashling gateway configuration")
	gatewayCommand.PersistentFlags().StringVarP(&envVarName, "env-var-name", "e", "MASHLING_CONFIG", "name of the environment variable that contains the base64 encoded mashling gateway configuration")
	gatewayCommand.PersistentFlags().BoolVarP(&loadFromEnv, "load-from-env", "l", false, "load the mashling gateway configuration from an environment variable")
	gatewayCommand.PersistentFlags().BoolVarP(&dev, "dev", "d", false, "run mashling in dev mode")
	gatewayCommand.PersistentFlags().StringVarP(&configCache, "config-cache", "C", ".cache", "location of the configuration artifacts cache")
	gatewayCommand.PersistentFlags().BoolVarP(&configCacheEnabled, "config-cache-enabled", "E", true, "cache post-processed configuration artifacts locally")
	gatewayCommand.PersistentFlags().BoolVarP(&pingEnabled, "ping-enabled", "p", true, "enable gateway ping service")
	gatewayCommand.PersistentFlags().StringVarP(&pingPort, "ping-port", "P", "9090", "configure mashling gateway ping service port")
}

var (
	gateway            model.Gateway
	watcher            *fsnotify.Watcher
	config             string
	envVarName         string
	loadFromEnv        bool
	dev                bool
	configCache        string
	configCacheEnabled bool
	pingEnabled        bool
	pingPort           string
)

var gatewayCommand = &cobra.Command{
	Use:   "mashling-gateway",
	Short: "mashling-gateway is a tool that serves up mashling instances",
	Long: "A static binary that executes Mashling gateway logic defined in a mashling.json configuration file. Complete documentation is available at https://github.com/TIBCOSoftware/mashling\n\n" +
		"Version: " + version.Version + "\nBuild Date: " + version.BuildDate + "\n",
	Run: run,
}

// Execute executes registered commands.
func Execute() {
	// Setup Mashling Logger
	if err := gatewayCommand.Execute(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func loadGateway() (err error) {
	// Load the appropriate Gateway instance from the configuration.
	if loadFromEnv {
		// Loading base64 encoded configuration from the env.
		gateway, err = model.LoadFromEnv(envVarName)
	} else {
		// Loading the configuration JSON from the specified file.
		gateway, err = model.LoadFromFile(config)
	}
	return err
}

func run(command *cobra.Command, args []string) {
	// Output friendly Mashling mascot with some build details.
	bannerTxt, err := assets.Asset("banner.txt")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	fmt.Println("\n", string(bannerTxt))
	logger.Info("Gateway Version: ", version.Version)
	logger.Info("Build Date: ", version.BuildDate)
	// Setup configuration artifacts cache.
	if configCacheEnabled {
		err = cache.Initialize("file", configCache)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}

	err = loadGateway()

	if err != nil {
		// Attempt to give insight into known  potentials errors.
		if gateway == nil {
			logger.Error(err)
			os.Exit(1)
		} else {
			logger.Error(err)
		}
		for _, errd := range gateway.Errors() {
			switch e := errd.(type) {
			case *gwerrors.UndefinedReference:
				logger.Errorf("%s: %s", e.Type(), e.Details())
			case *gwerrors.MissingDependency:
				logger.Error("Missing dependencies found: ", strings.Join(e.MissingDependencies, " "))
			default:
				logger.Errorf("Do not know how to handle error type %T!\n", e)
			}
		}
		// For now we still exit no matter what. The CLI handles missing errors.
		os.Exit(1)
	}

	if dev && !loadFromEnv {
		// Configuring file watcher to automatically reload configuration file contents in dev mode.
		watcher, err = fsnotify.NewWatcher()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		defer watcher.Close()

		go handleFileNotifications()

		err = watcher.Add(config)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}

	logger.Info("Schema Version: ", gateway.Version())
	logger.Info("App Version: ", gateway.AppVersion())
	logger.Info("App Description: ", gateway.Description())

	// Startup the configured gateway instance.
	gateway.Init(pingEnabled, pingPort)
	err = gateway.Start()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

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
				logger.Info("modified configuration file:", event.Name)
				reloadGatewayFromConfigurationFile()
			}
		case ferr := <-watcher.Errors:
			logger.Error(ferr)
		}
	}
}

// reloadGatewayFromConfigurationFile attempts to shutdown the current gateway instance and then re-start the gateway with the new configuration loaded.
func reloadGatewayFromConfigurationFile() {
	err := gateway.Stop()
	if err != nil {
		logger.Error("error stopping gateway:", err)
	}

	flogo.ResetGlobalContext()
	gateway, err = model.LoadFromFile(config)
	if err != nil {
		logger.Error("error re-loading gateway from file:", err)
	}
	err = gateway.Init(pingEnabled, pingPort)
	if err != nil {
		logger.Error("error re-initializing gateway:", err)
	}

	err = gateway.Start()
	if err != nil {
		logger.Error("error re-starting gateway:", err)
	}
}
