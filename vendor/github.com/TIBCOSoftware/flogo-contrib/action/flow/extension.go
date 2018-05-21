package flow

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/instance"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model/simple"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/support"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/tester"
)

// Provides the different extension points to the FlowBehavior Action
type ExtensionProvider interface {
	GetStateRecorder() instance.StateRecorder
	GetFlowTester() *tester.RestEngineTester

	GetDefaultFlowModel() *model.FlowModel
	GetFlowProvider() definition.Provider
	GetMapperFactory() definition.MapperFactory
	GetLinkExprManagerFactory() definition.LinkExprManagerFactory
}

//ExtensionProvider is the extension provider for the flow action
type DefaultExtensionProvider struct {
	flowProvider definition.Provider
	flowModel    *model.FlowModel
}

func NewDefaultExtensionProvider() *DefaultExtensionProvider {
	return &DefaultExtensionProvider{}
}

func (fp *DefaultExtensionProvider) GetFlowProvider() definition.Provider {

	if fp.flowProvider == nil {
		fp.flowProvider = &support.BasicRemoteFlowProvider{}
	}

	return fp.flowProvider
}

func (fp *DefaultExtensionProvider) GetDefaultFlowModel() *model.FlowModel {

	if fp.flowModel == nil {
		fp.flowModel = simple.New()
	}

	return fp.flowModel
}

func (fp *DefaultExtensionProvider) GetStateRecorder() instance.StateRecorder {
	return nil
}

func (fp *DefaultExtensionProvider) GetMapperFactory() definition.MapperFactory {
	return nil
}

func (fp *DefaultExtensionProvider) GetLinkExprManagerFactory() definition.LinkExprManagerFactory {
	return nil
}

//todo make FlowTester an interface
func (fp *DefaultExtensionProvider) GetFlowTester() *tester.RestEngineTester {
	return nil
}
