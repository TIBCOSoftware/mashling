package app

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/app/resource"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

func CreateTriggers(tConfigs []*trigger.Config, runner action.Runner) (map[string]trigger.Trigger, error) {

	triggers := make(map[string]trigger.Trigger)
	for _, tConfig := range tConfigs {

		_, exists := triggers[tConfig.Id]
		if exists {
			return nil, fmt.Errorf("Trigger with id '%s' already registered, trigger ids have to be unique", tConfig.Id)
		}

		triggerFactory := trigger.GetFactory(tConfig.Ref)

		if triggerFactory == nil {
			return nil, fmt.Errorf("Trigger Factory '%s' not registered", tConfig.Ref)
		}

		trg := triggerFactory.New(tConfig)

		if trg == nil {
			return nil, fmt.Errorf("cannot create Trigger nil for id '%s'", tConfig.Id)
		}

		tConfig.FixUp(trg.Metadata())

		initCtx := &initContext{handlers: make([]*trigger.Handler, 0, len(tConfig.Handlers))}

		var legacyRunner *trigger.LegacyRunner

		newTrg, isNew := trg.(trigger.Initializable)

		if !isNew {
			legacyRunner = trigger.NewLegacyRunner(runner, trg.Metadata())
		}

		//create handlers for that trigger and init
		for _, hConfig := range tConfig.Handlers {

			//create the action
			actionFactory := action.GetFactory(hConfig.Action.Ref)
			if actionFactory == nil {
				return nil, fmt.Errorf("Action Factory '%s' not registered", hConfig.Action.Ref)
			}

			act, err := actionFactory.New(hConfig.Action)
			if err != nil {
				return nil, err
			}

			handler := trigger.NewHandler(hConfig, act, trg.Metadata().Output, trg.Metadata().Reply, runner)
			initCtx.handlers = append(initCtx.handlers, handler)

			if !isNew {
				action.Register(hConfig.ActionId, act)
				legacyRunner.RegisterHandler(handler)
			}
		}

		if isNew {
			err := newTrg.Initialize(initCtx)
			if err != nil {
				return nil, err
			}
		} else {
			oldTrg, isOld := trg.(trigger.InitOld)
			if isOld {
				oldTrg.Init(legacyRunner)
			}
		}

		triggers[tConfig.Id] = trg
	}

	return triggers, nil
}

func RegisterResources(rConfigs []*resource.Config) error {

	if len(rConfigs) == 0 {
		return nil
	}

	for _, rConfig := range rConfigs {
		err := resource.Load(rConfig)
		if err != nil {
			return err
		}

	}

	return nil
}

type initContext struct {
	handlers []*trigger.Handler
}

func (ctx *initContext) GetHandlers() []*trigger.Handler {
	return ctx.handlers
}
