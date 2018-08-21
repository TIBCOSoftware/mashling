package engine

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/config"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/TIBCOSoftware/flogo-lib/util/managed"
	"sync"
	"github.com/TIBCOSoftware/flogo-lib/engine/channels"
)

var managedServices []managed.Managed
var lock = &sync.Mutex{}

// Interface for the engine behaviour
type Engine interface {
	// Init initialize the engine
	Init(directRunner bool) error

	// Start starts the engine
	Start() error

	// Stop stop the engine
	Stop() error

	// TriggerInfos get info for the triggers
	TriggerInfos() []*managed.Info
}

func LifeCycle(managedEntity managed.Managed)  {
	defer lock.Unlock()
	lock.Lock()
	managedServices = append(managedServices, managedEntity)
}


// engineImpl is the type for the Default Engine Implementation
type engineImpl struct {
	app            *app.Config
	initialized    bool
	logLevel       string
	actionRunner   action.Runner
	serviceManager *util.ServiceManager

	triggers     map[string]trigger.Trigger
	triggerInfos map[string]*managed.Info
}

// New creates a new Engine
func New(appCfg *app.Config) (Engine, error) {
	// App is required
	if appCfg == nil {
		return nil, errors.New("no App configuration provided")
	}
	// Name is required
	if len(appCfg.Name) == 0 {
		return nil, errors.New("no App name provided")
	}
	// Version is required
	if len(appCfg.Version) == 0 {
		return nil, errors.New("no App version provided")
	}

	//fix up app configuration if it is older
	//app.FixUpApp(appCfg)

	logLevel := config.GetLogLevel()

	return &engineImpl{app: appCfg, serviceManager: util.GetDefaultServiceManager(), logLevel: logLevel}, nil
}

func (e *engineImpl) Init(directRunner bool) error {

	if !e.initialized {
		e.initialized = true

		if directRunner {
			e.actionRunner = runner.NewDirect()
		} else {
			e.actionRunner = runner.NewPooled(NewPooledRunnerConfig())
		}

		propProvider := app.GetPropertyProvider()
		// Initialize the properties
		props, err := app.GetProperties(e.app.Properties)
		if err != nil {
			return err
		}
		propProvider.SetProperties(props)
		data.SetPropertyProvider(propProvider)

		actionFactories := action.Factories()
		for _, factory := range actionFactories {
			if initializable, ok := factory.(managed.Initializable); ok {

				if err := initializable.Init(); err != nil {
					return err
				}
			}
		}

		//add engine channels
		channelNames := e.app.Channels
		if len(channelNames) > 0 {
			for _, channelName := range channelNames {

				logger.Debugf("Creating Engine Channel '%s'", channelName)
				channels.Add(channelName)
			}
		}

		err = app.RegisterResources(e.app.Resources)
		if err != nil {
			return err
		}

		actions, err := app.CreateSharedActions(e.app.Actions)
		if err != nil {
			errorMsg := fmt.Sprintf("Error creating shared action instances - %s", err.Error())
			logger.Error(errorMsg)
			panic(errorMsg)
		}

		//todo add all actions to engine (will make cleanup easier)

		triggers, err := app.CreateTriggers(e.app.Triggers, actions, e.actionRunner)
		e.triggerInfos = make(map[string]*managed.Info)

		if err != nil {
			errorMsg := fmt.Sprintf("Error Creating trigger instances - %s", err.Error())
			logger.Error(errorMsg)
			panic(errorMsg)
		}

		e.triggers = triggers
	}

	return nil
}

//Start initializes and starts the Triggers and initializes the Actions
func (e *engineImpl) Start() error {

	logger.SetDefaultLogger("engine")

	logger.Debugf("Starting app [ %s ] with version [ %s ]", e.app.Name, e.app.Version)
	logger.Info("Engine Starting...")

	// Todo document RunnerType for engine configuration
	runnerType := GetRunnerType()
	err := e.Init(runnerType == "DIRECT")
	if err != nil {
		return err
	}

	logger.Info("Starting Services...")

	actionRunner := e.actionRunner.(interface{})

	if managedRunner, ok := actionRunner.(managed.Managed); ok {
		managed.Start("ActionRunner Service", managedRunner)
	}

	err = e.serviceManager.Start()

	if err != nil {
		logger.Error("Error Starting Services - " + err.Error())
	} else {
		logger.Info("Started Services")
	}

	if len(managedServices) > 0 {
		for _, mService := range managedServices {
			err = mService.Start()
			if err != nil {
				logger.Error("Error Starting Services - " + err.Error())
				//TODO Should we exit here?
			}
		}
	}

	// Start the triggers
	logger.Info("Starting Triggers...")

	var failed []string

	for key, value := range e.triggers {
		triggerInfo := &managed.Info{Name: key}
		err := managed.Start(fmt.Sprintf("Trigger [ %s ]", key), value)
		if err != nil {
			logger.Infof("Trigger [%s] failed to start due to error [%s]", key, err.Error())
			triggerInfo.Status = managed.StatusFailed
			triggerInfo.Error = err
			logger.Debugf("StackTrace: %s", debug.Stack())
			if config.StopEngineOnError() {
				logger.Debugf("{%s=true}. Stopping engine", config.ENV_STOP_ENGINE_ON_ERROR_KEY)
				logger.Info("Stopped")
				os.Exit(1)
			}
			failed = append(failed, key)
		} else {
			triggerInfo.Status = managed.StatusStarted
			logger.Infof("Trigger [ %s ]: Started", key)
			logger.Debugf("Trigger [ %s ] has ref [ %s ] and version [ %s ]", key, value.Metadata().ID, value.Metadata().Version)
		}

		e.triggerInfos[key] = triggerInfo
	}

	if len(failed) > 0 {
		//remove failed trigger, we have no use for them
		for _, triggerId := range failed {
			delete(e.triggers, triggerId)
		}
	}

	logger.Info("Triggers Started")

	logger.Info("Engine Started")

	return nil
}

func (e *engineImpl) Stop() error {
	logger.Info("Engine Stopping...")

	if channels.Count() > 0 {
		logger.Info("Closing Engine Channels...")
		channels.Close()
	}

	logger.Info("Stopping Triggers...")

	// Stop Triggers
	for trgId, tgr := range e.triggers {
		managed.Stop("Trigger [ "+trgId+" ]", tgr)
		e.triggerInfos[trgId].Status = managed.StatusStopped
	}

	logger.Info("Triggers Stopped")

	//TODO temporarily add services
	logger.Info("Stopping Services...")

	actionRunner := e.actionRunner.(interface{})

	if managedRunner, ok := actionRunner.(managed.Managed); ok {
		managed.Stop("ActionRunner", managedRunner)
	}

	err := e.serviceManager.Stop()

	if err != nil {
		logger.Error("Error Stopping Services - " + err.Error())
	} else {
		logger.Info("Stopped Services")
	}

	if len(managedServices) > 0 {
		for _, mService := range managedServices {
			err = mService.Stop()
			if err != nil {
				logger.Error("Error Stopping Services - " + err.Error())
			}
		}
	}

	logger.Info("Engine Stopped")
	return nil
}

func (e *engineImpl) TriggerInfos() []*managed.Info {

	infos := make([]*managed.Info, 0, len(e.triggerInfos))

	for _, info := range e.triggerInfos {
		infos = append(infos, info)
	}

	return infos
}

func RunEngine(e Engine) {

	err := e.Start()

	if err != nil {
		fmt.Println("Error starting engine", err.Error())
		os.Exit(1)
	}

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
