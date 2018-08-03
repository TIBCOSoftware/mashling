package flow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/instance"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	_ "github.com/TIBCOSoftware/flogo-contrib/action/flow/model/simple"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/support"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/tester"
	"github.com/TIBCOSoftware/flogo-lib/app/resource"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

const (
	FLOW_REF = "github.com/TIBCOSoftware/flogo-contrib/action/flow"

	ENV_FLOW_RECORD = "FLOGO_FLOW_RECORD"
)

type FlowAction struct {
	flowURI    string
	ioMetadata *data.IOMetadata
}

type ActionData struct {
	// The flow is a URI
	FlowURI string `json:"flowURI"`

	// The flow is embedded and uncompressed
	//DEPRECATED
	Flow json.RawMessage `json:"flow"`

	// The flow is a URI
	//DEPRECATED
	FlowCompressed json.RawMessage `json:"flowCompressed"`
}

var ep ExtensionProvider
var idGenerator *util.Generator
var record bool
var manager *support.FlowManager

//todo expose and support this properly
var maxStepCount = 1000000

//todo fix this
var metadata = &action.Metadata{ID: "github.com/TIBCOSoftware/flogo-contrib/action/flow", Async: true}

func init() {
	action.RegisterFactory(FLOW_REF, &ActionFactory{})
}

func SetExtensionProvider(provider ExtensionProvider) {
	ep = provider
}

type ActionFactory struct {
}

func (ff *ActionFactory) Init() error {

	if manager != nil {
		return nil
	}

	if ep == nil {
		testerEnabled := os.Getenv(tester.ENV_ENABLED)
		if strings.ToLower(testerEnabled) == "true" {
			ep = tester.NewExtensionProvider()

			sm := util.GetDefaultServiceManager()
			sm.RegisterService(ep.GetFlowTester())
			record = true
		} else {
			ep = NewDefaultExtensionProvider()
			record = recordFlows()
		}
	}

	definition.SetMapperFactory(ep.GetMapperFactory())
	definition.SetLinkExprManagerFactory(ep.GetLinkExprManagerFactory())

	if idGenerator == nil {
		idGenerator, _ = util.NewGenerator()
	}

	model.RegisterDefault(ep.GetDefaultFlowModel())
	manager = support.NewFlowManager(ep.GetFlowProvider())
	resource.RegisterManager(support.RESTYPE_FLOW, manager)

	return nil
}

func recordFlows() bool {
	recordFlows := os.Getenv(ENV_FLOW_RECORD)
	if len(recordFlows) == 0 {
		return false
	}
	b, _ := strconv.ParseBool(recordFlows)
	return b
}

func GetFlowManager() *support.FlowManager {
	return manager
}

func (ff *ActionFactory) New(config *action.Config) (action.Action, error) {

	flowAction := &FlowAction{}

	//temporary hack to support dynamic process running by tester
	if config.Data == nil {
		return flowAction, nil
	}

	var actionData ActionData
	err := json.Unmarshal(config.Data, &actionData)
	if err != nil {
		return nil, fmt.Errorf("faild to load flow action data '%s' error '%s'", config.Id, err.Error())
	}

	if len(actionData.FlowURI) > 0 {

		flowAction.flowURI = actionData.FlowURI
	} else {
		uri, err := createResource(&actionData)
		if err != nil {
			return nil, err
		}
		flowAction.flowURI = uri
	}

	if config.Metadata != nil {
		flowAction.ioMetadata = config.Metadata
	} else {
		//todo add flag to remove startup validation
		def, err := manager.GetFlow(flowAction.flowURI)
		if err != nil {
			return nil, err
		} else {
			if def == nil {
				return nil, errors.New("unable to resolve flow: " + flowAction.flowURI)
			}
		}

		flowAction.ioMetadata = def.Metadata()
	}

	return flowAction, nil
}

//Deprecated
func createResource(actionData *ActionData) (string, error) {

	manager := resource.GetManager(support.RESTYPE_FLOW)

	resourceCfg := &resource.Config{ID: "flow:" + strconv.Itoa(time.Now().Nanosecond())}

	if actionData.FlowCompressed != nil {
		resourceCfg.Compressed = true
		resourceCfg.Data = actionData.FlowCompressed
	} else if actionData.Flow != nil {
		resourceCfg.Data = actionData.Flow
	} else {
		return "", fmt.Errorf("flow not provided for FlowBehavior Action")
	}

	err := manager.LoadResource(resourceCfg)
	if err != nil {
		return "", err
	}

	return "res://" + resourceCfg.ID, nil
}

//Metadata get the Action's metadata
func (fa *FlowAction) Metadata() *action.Metadata {
	return metadata
}

func (fa *FlowAction) IOMetadata() *data.IOMetadata {
	return fa.ioMetadata
}

// Run implements action.Action.Run
//func (fa *FlowAction) Run(context context.Context, uri string, options interface{}, handler action.ResultHandler) error {
func (fa *FlowAction) Run(context context.Context, inputs map[string]*data.Attribute, handler action.ResultHandler) error {

	op := instance.OpStart
	retID := false
	var initialState *instance.IndependentInstance
	var flowURI string

	runOptions, exists := inputs["_run_options"]

	var execOptions *instance.ExecOptions

	if exists {
		ro, ok := runOptions.Value().(*instance.RunOptions)

		if ok {
			op = ro.Op
			retID = ro.ReturnID
			initialState = ro.InitialState
			flowURI = ro.FlowURI
			execOptions = ro.ExecOptions
		}
	}

	delete(inputs, "_run_options")

	if flowURI == "" {
		flowURI = fa.flowURI
	}

	logger.Infof("Running FlowAction for URI: '%s'", flowURI)

	//todo: catch panic
	//todo: consider switch to URI to dictate flow operation (ex. flow://blah/resume)

	var inst *instance.IndependentInstance

	switch op {
	case instance.OpStart:
		flowDef, err := manager.GetFlow(flowURI)
		if err != nil {
			return err
		}

		if flowDef == nil {
			return errors.New("flow not found for URI: " + flowURI)
		}

		instanceID := idGenerator.NextAsString()
		logger.Debug("Creating Flow Instance: ", instanceID)

		inst = instance.NewIndependentInstance(instanceID, flowURI, flowDef)
	case instance.OpResume:
		if initialState != nil {
			inst = initialState
			logger.Debug("Resuming Flow Instance: ", inst.ID())
		} else {
			return errors.New("unable to resume instance, initial state not provided")
		}
	case instance.OpRestart:
		if initialState != nil {
			inst = initialState
			instanceID := idGenerator.NextAsString()
			//flowDef, err := manager.GetFlow(flowURI)
			//if err != nil {
			//	return err
			//}

			//if flowDef.Metadata == nil {
			//	//flowDef.SetMetadata(fa.config.Metadata)
			//}
			err := inst.Restart(instanceID, manager)
			if err != nil {
				return err
			}

			logger.Debug("Restarting Flow Instance: ", instanceID)
		} else {
			return errors.New("unable to restart instance, initial state not provided")
		}
	}

	if execOptions != nil {
		logger.Debugf("Applying Exec Options to instance: %s", inst.ID())
		instance.ApplyExecOptions(inst, execOptions)
	}

	//todo how do we check if debug is enabled?
	//logInputs(inputs)

	logger.Debugf("Executing Flow Instance: %s", inst.ID())

	if op == instance.OpStart {

		inst.Start(inputs)
	} else {
		inst.UpdateAttrs(inputs)
	}

	stepCount := 0
	hasWork := true

	inst.SetResultHandler(handler)

	go func() {

		defer handler.Done()

		if !inst.FlowDefinition().ExplicitReply() || retID {

			idAttr, _ := data.NewAttribute("id", data.TypeString, inst.ID())
			results := map[string]*data.Attribute{
				"id": idAttr,
			}

			//todo remove
			//if old {
			//	dataAttr, _ := data.NewAttribute("data", data.OBJECT, &instance.IDResponse{ID: inst.ID()})
			//	results["data"] = dataAttr
			//	codeAttr, _ := data.NewAttribute("code", data.INTEGER, 200)
			//	results["code"] = codeAttr
			//}

			handler.HandleResult(results, nil)
		}

		for hasWork && inst.Status() < model.FlowStatusCompleted && stepCount < maxStepCount {
			stepCount++
			logger.Debugf("Step: %d", stepCount)
			hasWork = inst.DoStep()

			if record {
				ep.GetStateRecorder().RecordSnapshot(inst)
				ep.GetStateRecorder().RecordStep(inst)
			}
		}

		if inst.Status() == model.FlowStatusCompleted {
			returnData, err := inst.GetReturnData()
			handler.HandleResult(returnData, err)
		} else if inst.Status() == model.FlowStatusFailed {
			handler.HandleResult(nil, inst.GetError())
		}

		logger.Debugf("Done Executing flow instance [%s] - Status: %d", inst.ID(), inst.Status())

		if inst.Status() == model.FlowStatusCompleted {
			logger.Infof("Flow instance [%s] Completed Successfully", inst.ID())
		} else if inst.Status() == model.FlowStatusFailed {
			logger.Infof("Flow instance [%s] Failed", inst.ID())
		}
	}()

	return nil
}

func logInputs(attrs map[string]*data.Attribute) {
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

//func extractAttributes(inputs map[string]interface{}) []*data.Attribute {
//
//	size := len(inputs)
//
//	attrs := make([]*data.Attribute, 0, size)
//
//	//todo do special handling for complex_object metadata (merge or ref it)
//	for _, value := range inputs {
//
//		attr, _ := value.(*data.Attribute)
//		attrs = append(attrs, attr)
//	}
//
//	return attrs
//}
