package extension

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/instance"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/provider"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/tester"
	"github.com/TIBCOSoftware/flogo-contrib/model/simple"
)

//Provider is the extension provider for the flow action
type Provider struct {
	flowProvider provider.Provider
	flowModel    *model.FlowModel
}

func New() *Provider {
	return &Provider{}
}

func (fp *Provider) GetFlowProvider() provider.Provider {

	if fp.flowProvider == nil {
		fp.flowProvider = provider.NewRemoteFlowProvider()
	}

	return fp.flowProvider
}

func (fp *Provider) GetFlowModel() *model.FlowModel {

	if fp.flowModel == nil {
		fp.flowModel = simple.New()
	}

	return fp.flowModel
}

func (fp *Provider) GetStateRecorder() instance.StateRecorder {
	return nil
}

func (fp *Provider) GetMapperFactory() definition.MapperFactory {
	return nil
}

func (fp *Provider) GetLinkExprManagerFactory() definition.LinkExprManagerFactory {
	return nil
}

//todo make FlowTester an interface
func (fp *Provider) GetFlowTester() *tester.RestEngineTester {
	return nil
}
