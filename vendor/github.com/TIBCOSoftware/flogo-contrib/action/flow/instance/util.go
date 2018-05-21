package instance

import (
	"errors"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

func applyInputMapper(taskInst *TaskInst) error {

	// get the input mapper
	inputMapper := taskInst.task.ActivityConfig().InputMapper()

	master := taskInst.flowInst.master

	if master.patch != nil {
		// check if the patch has a overriding mapper
		mapper := master.patch.GetInputMapper(taskInst.task.ID())
		if mapper != nil {
			inputMapper = mapper
		}
	}

	if inputMapper != nil {
		logger.Debug("Applying InputMapper")

		var inputScope data.Scope
		inputScope = taskInst.flowInst

		if taskInst.workingData != nil {
			inputScope = NewWorkingDataScope(taskInst.flowInst, taskInst.workingData)
		}

		err := inputMapper.Apply(inputScope, taskInst.InputScope())

		if err != nil {
			return err
		}
	}

	return nil
}

func applyInputInterceptor(taskInst *TaskInst) bool {

	master := taskInst.flowInst.master

	if master.interceptor != nil {

		// check if this task as an interceptor
		taskInterceptor := master.interceptor.GetTaskInterceptor(taskInst.task.ID())

		if taskInterceptor != nil {

			logger.Debug("Applying Interceptor")

			if len(taskInterceptor.Inputs) > 0 {
				// override input attributes
				for _, attribute := range taskInterceptor.Inputs {

					logger.Debugf("Overriding Attr: %s = %s", attribute.Name(), attribute.Value())

					//todo: validation
					taskInst.InputScope().SetAttrValue(attribute.Name(), attribute.Value())
				}
			}

			// check if we should not evaluate the task
			return !taskInterceptor.Skip
		}
	}

	return true
}

func applyOutputInterceptor(taskInst *TaskInst) {

	master := taskInst.flowInst.master

	if master.interceptor != nil {

		// check if this task as an interceptor and overrides ouputs
		taskInterceptor := master.interceptor.GetTaskInterceptor(taskInst.task.ID())
		if taskInterceptor != nil && len(taskInterceptor.Outputs) > 0 {
			// override output attributes
			for _, attribute := range taskInterceptor.Outputs {

				//todo: validation
				taskInst.OutputScope().SetAttrValue(attribute.Name(), attribute.Value())
			}
		}
	}
}

// applyOutputMapper applies the output mapper, returns flag indicating if
// there was an output mapper
func applyOutputMapper(taskInst *TaskInst) (bool, error) {

	// get the Output Mapper for the TaskOld if one exists
	outputMapper := taskInst.task.ActivityConfig().OutputMapper()

	master := taskInst.flowInst.master

	if master.patch != nil {
		// check if the patch overrides the Output Mapper
		mapper := master.patch.GetOutputMapper(taskInst.task.ID())
		if mapper != nil {
			outputMapper = mapper
		}
	}

	if outputMapper != nil {
		logger.Debug("Applying OutputMapper")
		err := outputMapper.Apply(taskInst.OutputScope(), taskInst.flowInst)

		return true, err
	}

	return false, nil
}

func GetFlowIOMetadata(flowURI string) (*data.IOMetadata, error) {
	manager := support.GetFlowManager()
	def, err := manager.GetFlow(flowURI)

	if err != nil {
		return nil, err
	}

	if def == nil {
		return nil, errors.New("unable to resolve flow: " + flowURI)
	}

	return def.Metadata(), nil
}

func StartSubFlow(ctx activity.Context, flowURI string, inputs map[string]*data.Attribute) error {

	taskInst, ok := ctx.(*TaskInst)

	if !ok {
		return errors.New("unable to create subFlow using this context")
	}

	manager := support.GetFlowManager()
	def, err := manager.GetFlow(flowURI)

	if err != nil {
		return err
	}

	if def == nil {
		return errors.New("unable to resolve subflow: " + flowURI)
	}

	//todo make sure that there is only one subFlow per taskinst
	flowInst := taskInst.flowInst.master.newEmbeddedInstance(taskInst, flowURI, def)

	if err != nil {
		return err
	}

	logger.Debugf("starting embedded subflow `%s`", flowInst.Name())

	err = taskInst.flowInst.master.startEmbedded(flowInst, inputs)
	if err != nil {
		return err
	}

	return nil
}
