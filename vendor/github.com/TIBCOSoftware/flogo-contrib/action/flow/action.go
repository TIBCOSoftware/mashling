package flow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/extension"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/instance"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/provider"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/tester"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

const (
	FLOW_REF = "github.com/TIBCOSoftware/flogo-contrib/action/flow"
)

// ActionOptions are the options for the FlowAction
type ActionOptions struct {
	MaxStepCount int
	Record       bool
}

type FlowAction struct {
	idGenerator   *util.Generator
	actionOptions *ActionOptions
	config        *action.Config
}

// Provides the different extension points to the Flow Action
type ExtensionProvider interface {
	GetFlowProvider() provider.Provider
	GetFlowModel() *model.FlowModel
	GetStateRecorder() instance.StateRecorder
	GetMapperFactory() definition.MapperFactory
	GetLinkExprManagerFactory() definition.LinkExprManagerFactory
	GetFlowTester() *tester.RestEngineTester
}

var actionMu sync.Mutex
var ep ExtensionProvider
//var flowAction *FlowAction
var idGenerator *util.Generator
var record bool

func init() {
	action.RegisterFactory(FLOW_REF, &FlowFactory{})
}

func SetExtensionProvider(provider ExtensionProvider) {
	actionMu.Lock()
	defer actionMu.Unlock()

	ep = provider
}

type FlowFactory struct{}

func (ff *FlowFactory) New(config *action.Config) action.Action {

	actionMu.Lock()
	defer actionMu.Unlock()

	var flowAction *FlowAction

	if flowAction == nil {
		options := &ActionOptions{Record: record}

		if ep == nil {
			testerEnabled := os.Getenv(tester.ENV_ENABLED)
			if strings.ToLower(testerEnabled) == "true" {
				ep = tester.NewExtensionProvider()

				sm := util.GetDefaultServiceManager()
				sm.RegisterService(ep.GetFlowTester())
				record = true
				options.Record = true
			} else {
				ep = extension.New()
			}

			definition.SetMapperFactory(ep.GetMapperFactory())
			definition.SetLinkExprManagerFactory(ep.GetLinkExprManagerFactory())
		}

		if idGenerator == nil {
			idGenerator, _ = util.NewGenerator()
		}

		if options.MaxStepCount < 1 {
			options.MaxStepCount = int(^uint16(0))
		}

		flowAction = &FlowAction{config: config}

		flowAction.actionOptions = options
		flowAction.idGenerator = idGenerator
	}

	//temporary hack to support dynamic process running by tester
	if config.Data == nil {
		return flowAction
	}

	var flavor Flavor
	err := json.Unmarshal(config.Data, &flavor)
	if err != nil {
		errorMsg := fmt.Sprintf("Error while loading flow '%s' error '%s'", config.Id, err.Error())
		logger.Errorf(errorMsg)
		panic(errorMsg)
	}

	if len(flavor.Flow) > 0 {
		// It is an uncompressed and embedded flow
		err := ep.GetFlowProvider().AddUncompressedFlow(config.Id, flavor.Flow)
		if err != nil {
			errorMsg := fmt.Sprintf("Error while loading uncompressed flow '%s' error '%s'", config.Id, err.Error())
			logger.Errorf(errorMsg)
			panic(errorMsg)
		}
		return flowAction
	}

	if len(flavor.FlowCompressed) > 0 {
		// It is a compressed and embedded flow
		err := ep.GetFlowProvider().AddCompressedFlow(config.Id, string(flavor.FlowCompressed[:]))
		if err != nil {
			errorMsg := fmt.Sprintf("Error while loading compressed flow '%s' error '%s'", config.Id, err.Error())
			logger.Errorf(errorMsg)
			panic(errorMsg)
		}
		return flowAction
	}

	if len(flavor.FlowURI) > 0 {
		// It is a URI flow
		err := ep.GetFlowProvider().AddFlowURI(config.Id, string(flavor.FlowURI[:]))
		if err != nil {
			errorMsg := fmt.Sprintf("Error while loading flow URI '%s' error '%s'", config.Id, err.Error())
			logger.Errorf(errorMsg)
			panic(errorMsg)
		}
		return flowAction
	}

	errorMsg := fmt.Sprintf("No flow found in action data for id '%s'", config.Id)
	logger.Errorf(errorMsg)
	panic(errorMsg)

	return flowAction
}

//Config get the Action's config
func (fa *FlowAction) Config() *action.Config {
	return fa.config
}

//Metadata get the Action's metadata
func (fa *FlowAction) Metadata() *action.Metadata {
	return nil
}

// Run implements action.Action.Run
//func (fa *FlowAction) Run(context context.Context, uri string, options interface{}, handler action.ResultHandler) error {
func (fa *FlowAction) Run(context context.Context, inputs []*data.Attribute, options map[string]interface{}, handler action.ResultHandler) error {

	op := instance.OpStart
	retID := false
	var initialState *instance.Instance
	var flowURI string
	mh := data.GetMapHelper()

	oldOptions, old := options["deprecated_options"]
	var execOptions *instance.ExecOptions

	if old {
		ro, ok := oldOptions.(*instance.RunOptions)

		if ok {
			op = ro.Op
			retID = ro.ReturnID
			initialState = ro.InitialState
			flowURI = ro.FlowURI
			execOptions = ro.ExecOptions
		}

	} else {
		//todo enable support for ExecOption when using new action options

		if v, ok := mh.GetInt(options, "op"); ok {
			op = v
		}
		if v, ok := mh.GetBool(options, "returnId"); ok {
			retID = v
		}
		if v, ok := mh.GetString(options, "flowURI"); ok {
			flowURI = v
		}
		if v, ok := options["initialState"]; ok {
			if v, ok := v.(*instance.Instance); ok {
				initialState = v
			}
		}
	}

	if flowURI == "" {
		flowURI = fa.Config().Id
	}

	logger.Infof("In Flow Run uri: '%s'", flowURI)

	//todo: catch panic
	//todo: consider switch to URI to dictate flow operation (ex. flow://blah/resume)

	var inst *instance.Instance

	switch op {
	case instance.OpStart:
		flowDef, err := ep.GetFlowProvider().GetFlow(flowURI)
		if err != nil {
			return err
		}

		instanceID := fa.idGenerator.NextAsString()
		logger.Debug("Creating Instance: ", instanceID)

		inst = instance.New(instanceID, flowURI, flowDef, ep.GetFlowModel())
	case instance.OpResume:
		if initialState != nil {
			inst = initialState
			logger.Debug("Resuming Instance: ", inst.ID())
		} else {
			return errors.New("unable to resume instance, initial state not provided")
		}
	case instance.OpRestart:
		if initialState != nil {
			inst = initialState
			instanceID := fa.idGenerator.NextAsString()
			inst.Restart(instanceID, ep.GetFlowProvider())

			logger.Debug("Restarting Instance: ", instanceID)
		} else {
			return errors.New("unable to restart instance, initial state not provided")
		}
	}

	if execOptions != nil {
		logger.Debugf("Applying Exec Options to instance: %s\n", inst.ID())
		instance.ApplyExecOptions(inst, execOptions)
	}

	//todo how do we check if debug is enabled?
	logInputs(inputs)

	if op == instance.OpStart {
		inst.Start(inputs)
	} else {
		inst.UpdateAttrs(inputs)
	}

	logger.Debugf("Executing instance: %s\n", inst.ID())

	stepCount := 0
	hasWork := true

	inst.InitActionContext(fa.config, handler)

	go func() {

		defer handler.Done()

		if !inst.Flow.ExplicitReply() || retID {

			idAttr, _ := data.NewAttribute("id", data.STRING, inst.ID())
			results := map[string]*data.Attribute{
				"id": idAttr,
			}

			if old {
				dataAttr, _ := data.NewAttribute("data", data.OBJECT, &instance.IDResponse{ID: inst.ID()})
				results["data"] = dataAttr
				codeAttr, _ := data.NewAttribute("code", data.INTEGER, 200)
				results["code"] = codeAttr
			}

			handler.HandleResult(results, nil)
		}

		for hasWork && inst.Status() < instance.StatusCompleted && stepCount < fa.actionOptions.MaxStepCount {
			stepCount++
			logger.Debugf("Step: %d\n", stepCount)
			hasWork = inst.DoStep()

			if fa.actionOptions.Record {
				ep.GetStateRecorder().RecordSnapshot(inst)
				ep.GetStateRecorder().RecordStep(inst)
			}
		}

		if inst.Status() == instance.StatusCompleted {
			returnData, err := inst.GetReturnData()
			handler.HandleResult(returnData, err)
		}

		//if retID {
		//
		//	resp := map[string]*data.Attribute {
		//		"id": data.NewAttribute("id", data.STRING, inst.ID()),
		//	}
		//
		//	if old {
		//		resp["data"] = data.NewAttribute("data", data.OBJECT, &instance.IDResponse{ID: inst.ID()})
		//		resp["code"] = data.NewAttribute("code", data.INTEGER, 200)
		//	}
		//
		//	handler.HandleResult(resp, nil)
		//}

		logger.Debugf("Done Executing A.instance [%s] - Status: %d\n", inst.ID(), inst.Status())

		if inst.Status() == instance.StatusCompleted {
			logger.Infof("Flow [%s] Completed", inst.ID())
		}
	}()

	return nil
}

func logInputs(attrs []*data.Attribute) {
	if len(attrs) > 0 {
		logger.Debug("Input Attributes:")
		for _, attr := range attrs {

			if attr == nil {
				logger.Error("Nil Attribute passed as input")
			} else {
				logger.Debugf(" Attr:%s, Type:%s, Value:%v", attr.Name(), attr.Type().String(), attr.Value())
			}
		}
	}
}

func extractAttributes(inputs map[string]interface{}) []*data.Attribute {

	size := len(inputs)

	attrs := make([]*data.Attribute, 0, size)

	//todo do special handling for complex_object metadata (merge or ref it)
	for _, value := range inputs {

		attr, _ := value.(*data.Attribute)
		attrs = append(attrs, attr)
	}

	return attrs
}
