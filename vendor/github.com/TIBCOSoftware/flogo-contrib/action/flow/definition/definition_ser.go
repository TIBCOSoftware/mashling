package definition

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	flowutil "github.com/TIBCOSoftware/flogo-contrib/action/flow/util"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

// DefinitionRep is a serializable representation of a flow Definition
type DefinitionRep struct {
	ExplicitReply bool   `json:"explicitReply"`
	Name          string `json:"name"`
	ModelID       string `json:"model"`

	Metadata   *data.IOMetadata  `json:"metadata"`
	Attributes []*data.Attribute `json:"attributes,omitempty"`

	Tasks []*TaskRep `json:"tasks"`
	Links []*LinkRep `json:"links"`

	ErrorHandler *ErrorHandlerRep `json:"errorHandler"`

	//deprecated
	RootTask         *TaskRepOld `json:"rootTask"`
	ErrorHandlerTask *TaskRepOld `json:"errorHandlerTask"`
}

// ErrorHandlerRep is a serializable representation of the error flow
type ErrorHandlerRep struct {
	Tasks []*TaskRep `json:"tasks"`
	Links []*LinkRep `json:"links"`
}

// TaskRep is a serializable representation of a flow task
type TaskRep struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Name     string                 `json:"name"`
	Settings map[string]interface{} `json:"settings"`

	ActivityCfgRep *ActivityConfigRep `json:"activity"`
}

// ActivityConfigRep is a serializable representation of an activity configuration
type ActivityConfigRep struct {
	Ref         string                 `json:"ref"`
	Mappings    *Mappings              `json:"mappings,omitempty"`
	Settings    map[string]interface{} `json:"settings"`
	InputAttrs  map[string]interface{} `json:"input,omitempty"`
	OutputAttrs map[string]interface{} `json:"output,omitempty"`
}

// LinkRep is a serializable representation of a flow LinkOld
type LinkRep struct {
	Type string `json:"type"`

	Name   string `json:"name"`
	ToID   string `json:"to"`
	FromID string `json:"from"`
	Value  string `json:"value"`
}

// Mappings is a collection of input & output mappings
type Mappings struct {
	Input  []*data.MappingDef `json:"input,omitempty"`
	Output []*data.MappingDef `json:"output,omitempty"`
}

// NewDefinition creates a flow Definition from a serializable
// definition representation
func NewDefinition(rep *DefinitionRep) (def *Definition, err error) {

	defer util.HandlePanic("NewDefinition", &err)

	if rep.RootTask != nil {
		return definitionFromOldRep(rep)
	}

	def = &Definition{}
	def.name = rep.Name
	def.modelID = rep.ModelID
	def.metadata = rep.Metadata
	def.explicitReply = rep.ExplicitReply
	if len(rep.Attributes) > 0 {
		def.attrs = make(map[string]*data.Attribute, len(rep.Attributes))

		for _, value := range rep.Attributes {
			def.attrs[value.Name()] = value
		}
	}

	def.tasks = make(map[string]*Task)
	def.links = make(map[int]*Link)

	if len(rep.Tasks) != 0 {

		for _, taskRep := range rep.Tasks {

			task, err := createTask(def, taskRep)

			if err != nil {
				return nil, err
			}
			def.tasks[task.id] = task
		}
	}

	if len(rep.Links) != 0 {

		for id, linkRep := range rep.Links {

			link, err := createLink(def.tasks, linkRep, id)
			if err != nil {
				return nil, err
			}

			def.links[link.id] = link
		}
	}

	if rep.ErrorHandler != nil {

		errorHandler := &ErrorHandler{}
		errorHandler.tasks = make(map[string]*Task)
		errorHandler.links = make(map[int]*Link)
		def.errorHandler = errorHandler

		if len(rep.ErrorHandler.Tasks) != 0 {

			for _, taskRep := range rep.ErrorHandler.Tasks {

				task, err := createTask(def, taskRep)

				if err != nil {
					return nil, err
				}
				errorHandler.tasks[task.id] = task
			}
		}

		if len(rep.ErrorHandler.Links) != 0 {

			idOffset := len(rep.Links)

			for id, linkRep := range rep.ErrorHandler.Links {

				link, err := createLink(errorHandler.tasks, linkRep, id+idOffset)
				if err != nil {
					return nil, err
				}
				errorHandler.links[link.id] = link
			}
		}

	}

	return def, nil
}

func createTask(def *Definition, rep *TaskRep) (*Task, error) {
	task := &Task{}
	task.id = rep.ID
	task.name = rep.Name
	task.definition = def

	if rep.Type != "" {
		if !flowutil.IsValidTaskType(def.modelID, rep.Type) {
			return nil, errors.New("Unsupported task type: " + rep.Type)
		}
		task.typeID = rep.Type
	}

	if len(rep.Settings) > 0 {
		task.settings = make(map[string]interface{}, len(rep.Settings))

		for name, value := range rep.Settings {
			task.settings[name] = resolveSettingValue(name, value)
		}
	}

	if rep.ActivityCfgRep != nil {

		actCfg, err := createActivityConfig(task, rep.ActivityCfgRep)

		if err != nil {
			return nil, err
		}

		if actCfg.Activity.Metadata().ProducesResult || def.explicitReply {
			def.explicitReply = true
		}

		task.activityCfg = actCfg
	}

	return task, nil
}

func createActivityConfig(task *Task, rep *ActivityConfigRep) (*ActivityConfig, error) {

	if rep.Ref == "" {
		return nil, errors.New("Activity Not Specified for Task :" + task.ID())
	}

	act := activity.Get(rep.Ref)
	if act == nil {
		return nil, errors.New("Unsupported Activity:" + rep.Ref)
	}

	activityCfg := &ActivityConfig{}
	activityCfg.Activity = act

	//todo need to fix this
	task.activityCfg = activityCfg

	if len(rep.Settings) > 0 {
		activityCfg.settings = make(map[string]*data.Attribute, len(rep.Settings))

		for name, value := range rep.Settings {

			attr := act.Metadata().Settings[name]

			if attr != nil {
				//var err error
				//todo handle error
				activityCfg.settings[name], _ = data.NewAttribute(name, attr.Type(), resolveSettingValue(name, value))
			}
		}
	}

	inputAttrs := rep.InputAttrs

	if len(inputAttrs) > 0 {
		activityCfg.inputAttrs = make(map[string]*data.Attribute, len(inputAttrs))

		for name, value := range inputAttrs {

			attr := act.Metadata().Input[name]

			if attr != nil {
				//var err error
				//todo handle error
				activityCfg.inputAttrs[name], _ = data.NewAttribute(name, attr.Type(), value)
			}
		}
	}

	outputAttrs := rep.OutputAttrs

	if len(outputAttrs) > 0 {

		activityCfg.outputAttrs = make(map[string]*data.Attribute, len(outputAttrs))

		for name, value := range outputAttrs {

			attr := act.Metadata().Output[name]

			if attr != nil {
				//var err error
				//todo handle error
				activityCfg.outputAttrs[name], _ = data.NewAttribute(name, attr.Type(), value)
			}
		}
	}

	// create mappers
	if rep.Mappings != nil {
		if rep.Mappings.Input != nil {
			activityCfg.inputMapper = GetMapperFactory().NewActivityInputMapper(task, &data.MapperDef{Mappings: rep.Mappings.Input})
		}
		if rep.Mappings.Output != nil {
			activityCfg.outputMapper = GetMapperFactory().NewActivityOutputMapper(task, &data.MapperDef{Mappings: rep.Mappings.Output})
		} else {
			activityCfg.outputMapper = GetMapperFactory().GetDefaultActivityOutputMapper(task)
		}
	}

	//If outmapper still empty set to default
	if activityCfg.outputMapper == nil {
		activityCfg.outputMapper = GetMapperFactory().GetDefaultActivityOutputMapper(task)
	}

	return activityCfg, nil
}

func resolveSettingValue(setting string, value interface{}) interface{} {

	strVal, ok := value.(string)

	if ok && len(strVal) > 0 && strVal[0] == '$' {
		v, err := data.GetBasicResolver().Resolve(strVal, nil)

		if err == nil {

			logger.Debugf("Resolved setting [%s: %s] to : %v", setting, value, v)
			return v
		}
	}

	return value
}

func createLink(tasks map[string]*Task, linkRep *LinkRep, id int) (*Link, error) {

	link := &Link{}
	link.id = id
	link.linkType = LtDependency

	if len(linkRep.Type) > 0 {
		switch linkRep.Type {
		case "default", "dependency", "0":
			link.linkType = LtDependency
		case "expression", "1":
			link.linkType = LtExpression
		case "label", "2":
			link.linkType = LtLabel
		case "error", "3":
			link.linkType = LtError
		default:
			logger.Warnf("Unsupported link type '%s', using default link")
		}
	}

	link.value = linkRep.Value
	link.fromTask = tasks[linkRep.FromID]
	link.toTask = tasks[linkRep.ToID]

	if link.toTask == nil {
		strId := strconv.Itoa(link.ID())
		return nil, errors.New("Link[" + strId + "]: ToTask '" + linkRep.ToID + "' not found")
	}

	if link.fromTask == nil {
		strId := strconv.Itoa(link.ID())
		return nil, errors.New("Link[" + strId + "]: FromTask '" + linkRep.FromID + "' not found")
	}

	// add this link as predecessor "fromLink" to the "toTask"
	link.toTask.fromLinks = append(link.toTask.fromLinks, link)

	// add this link as successor "toLink" to the "fromTask"
	link.fromTask.toLinks = append(link.fromTask.toLinks, link)

	return link, nil
}

///////////////////////////
// DEPRECATED

func definitionFromOldRep(rep *DefinitionRep) (def *Definition, err error) {

	def = &Definition{}
	def.name = rep.Name
	def.modelID = rep.ModelID
	def.metadata = rep.Metadata
	def.explicitReply = rep.ExplicitReply
	if len(rep.Attributes) > 0 {
		def.attrs = make(map[string]*data.Attribute, len(rep.Attributes))

		for _, value := range rep.Attributes {
			def.attrs[value.Name()] = value
		}
	}

	def.tasks = make(map[string]*Task)
	def.links = make(map[int]*Link)

	// support for deprecated flow format
	err = addTasksOld(def, rep.RootTask, false)
	if err != nil {
		return nil, err
	}

	err = addLinksOld(def, rep.RootTask, false)

	if err != nil {
		return nil, err
	}

	if rep.ErrorHandlerTask != nil {

		errorHandler := &ErrorHandler{}
		errorHandler.tasks = make(map[string]*Task)
		errorHandler.links = make(map[int]*Link)
		def.errorHandler = errorHandler

		err = addTasksOld(def, rep.ErrorHandlerTask, true)
		if err != nil {
			return nil, err
		}

		addLinksOld(def, rep.ErrorHandlerTask, true)
	}

	return def, nil
}

// TaskRepOld is a serializable representation of a flow TaskOld
//Deprecated
type TaskRepOld struct {
	// Using interface{} type to support backward compatibility changes since Id was
	// int before, change to string once BC is removed
	ID          interface{} `json:"id"`
	TypeID      int         `json:"type"`
	ActivityRef string      `json:"activityRef"`
	Name        string      `json:"name"`

	Tasks []*TaskRepOld `json:"tasks,omitempty"`
	Links []*LinkRepOld `json:"links,omitempty"`

	Mappings    *Mappings              `json:"mappings,omitempty"`
	InputAttrs  map[string]interface{} `json:"input,omitempty"`
	OutputAttrs map[string]interface{} `json:"output,omitempty"`

	Settings map[string]interface{} `json:"settings"`

	//keep temporarily for backwards compatibility
	InputAttrsOld  map[string]interface{} `json:"inputs,omitempty"`
	OutputAttrsOld map[string]interface{} `json:"outputs,omitempty"`
	InputMappings  []*data.MappingDef     `json:"inputMappings,omitempty"`
	OutputMappings []*data.MappingDef     `json:"outputMappings,omitempty"`
	Attributes     []*data.Attribute      `json:"attributes,omitempty"`
	ActivityType   string                 `json:"activityType"`
}

// LinkRepOld is a serializable representation of a flow LinkOld
//Deprecated
type LinkRepOld struct {
	ID   int    `json:"id"`
	Type int    `json:"type"`
	Name string `json:"name"`
	// Using interface{} type to support backward compatibility changes since Id was
	// int before, change to string once BC is removed
	ToID interface{} `json:"to"`
	// Using interface{} type to support backward compatibility changes since Id was
	// int before, change to string once BC is removed
	FromID interface{} `json:"from"`
	Value  string      `json:"value"`
}

//Deprecated
func addTasksOld(def *Definition, rootTaskRep *TaskRepOld, eh bool) error {
	// flow  tasks
	if len(rootTaskRep.Tasks) > 0 {

		for _, childTaskRep := range rootTaskRep.Tasks {

			task, err := createTaskFromOld(def, childTaskRep)

			if err != nil {
				return err
			}
			if eh {
				def.errorHandler.tasks[task.id] = task
			} else {
				def.tasks[task.id] = task
			}
		}
	}

	return nil
}

//Deprecated
func createTaskFromOld(def *Definition, rep *TaskRepOld) (*Task, error) {

	task := &Task{}

	// Workaround to support Backwards compatibility
	// Remove once rep.ID is string
	task.id = convertInterfaceToString(rep.ID)
	// Workaround to support Backwards compatibility
	// Remove once rep.ID is string

	//temporary hack
	if rep.TypeID == 2 {
		task.typeID = "iterator"
	} else {
		task.typeID = ""
	}

	task.name = rep.Name
	task.definition = def

	actCfg, err := createActivityConfigFromOld(task, rep)

	if err != nil {
		return nil, err
	}

	if actCfg != nil && (actCfg.Activity.Metadata().ProducesResult || def.explicitReply) {
		def.explicitReply = true
	}

	task.activityCfg = actCfg

	return task, nil
}

//Deprecated
func createActivityConfigFromOld(task *Task, rep *TaskRepOld) (*ActivityConfig, error) {

	if rep.ActivityRef == "" {
		return nil, nil
	}

	act := activity.Get(rep.ActivityRef)
	if act == nil {
		return nil, errors.New("Unsupported Activity:" + rep.ActivityRef)
	}

	activityCfg := &ActivityConfig{}
	activityCfg.Activity = act

	//todo need to fix this
	task.activityCfg = activityCfg

	// Keep for now, DEPRECATE "attributes" section from flogo.json
	if len(rep.Settings) > 0 {
		task.settings = make(map[string]interface{}, len(rep.Settings))

		for name, value := range rep.Settings {
			task.settings[name] = resolveSettingValue(name, value)
		}
	}

	// create mappers
	if rep.Mappings != nil {
		if rep.Mappings.Input != nil {
			fixupMappings(rep.Mappings.Input)

			activityCfg.inputMapper = GetMapperFactory().NewActivityInputMapper(task, &data.MapperDef{Mappings: rep.Mappings.Input})
		}
		if rep.Mappings.Output != nil {
			activityCfg.outputMapper = GetMapperFactory().NewActivityOutputMapper(task, &data.MapperDef{Mappings: rep.Mappings.Output})
		} else {
			activityCfg.outputMapper = GetMapperFactory().GetDefaultActivityOutputMapper(task)
		}
	}

	// create mappers
	if rep.Mappings != nil {
		if rep.Mappings.Input != nil {
			fixupMappings(rep.Mappings.Input)
			activityCfg.inputMapper = GetMapperFactory().NewActivityInputMapper(task, &data.MapperDef{Mappings: rep.Mappings.Input})
		}
		if rep.Mappings.Output != nil {
			activityCfg.outputMapper = GetMapperFactory().NewActivityOutputMapper(task, &data.MapperDef{Mappings: rep.Mappings.Output})
		}
	} else {
		//temporary support for old configuration
		if rep.InputMappings != nil {
			fixupMappings(rep.InputMappings)
			activityCfg.inputMapper = GetMapperFactory().NewActivityInputMapper(task, &data.MapperDef{Mappings: rep.InputMappings})
		}
		if rep.OutputMappings != nil {
			activityCfg.outputMapper = GetMapperFactory().NewActivityOutputMapper(task, &data.MapperDef{Mappings: rep.OutputMappings})
		}
	}

	if activityCfg.outputMapper == nil {
		activityCfg.outputMapper = GetMapperFactory().GetDefaultActivityOutputMapper(task)
	}

	inputAttrs := rep.InputAttrs

	//for backwards compatibility
	if len(inputAttrs) == 0 {
		inputAttrs = rep.InputAttrsOld
	}

	if len(inputAttrs) > 0 {
		activityCfg.inputAttrs = make(map[string]*data.Attribute, len(inputAttrs))

		for name, value := range inputAttrs {

			attr := act.Metadata().Input[name]

			if attr != nil {
				//var err error
				//todo handle error
				activityCfg.inputAttrs[name], _ = data.NewAttribute(name, attr.Type(), value)
			}
		}
	} else if len(rep.Attributes) > 0 {

		activityCfg.inputAttrs = make(map[string]*data.Attribute, len(inputAttrs))

		for _, attr := range rep.Attributes {

			if attr != nil {
				//var err error
				//todo handle error
				activityCfg.inputAttrs[attr.Name()] = attr
			}
		}
	}

	outputAttrs := rep.OutputAttrs

	//for backwards compatibility
	if len(outputAttrs) == 0 {
		outputAttrs = rep.OutputAttrsOld
	}

	if len(outputAttrs) > 0 {

		activityCfg.outputAttrs = make(map[string]*data.Attribute, len(outputAttrs))

		for name, value := range outputAttrs {

			attr := act.Metadata().Output[name]

			if attr != nil {
				//var err error
				//todo handle error
				activityCfg.outputAttrs[name], _ = data.NewAttribute(name, attr.Type(), value)
			}
		}
	}

	return activityCfg, nil
}

//Deprecated
func addLinksOld(def *Definition, rep *TaskRepOld, eh bool) error {

	numLinks := len(rep.Links)

	if numLinks > 0 {

		for _, linkRep := range rep.Links {

			link := &Link{}
			link.id = linkRep.ID
			//link.Parent = task
			//link.Definition = pd
			link.linkType = LinkType(linkRep.Type)
			link.value = linkRep.Value

			if eh {
				link.fromTask = def.errorHandler.tasks[convertInterfaceToString(linkRep.FromID)]
				link.toTask = def.errorHandler.tasks[convertInterfaceToString(linkRep.ToID)]

			} else {
				link.fromTask = def.tasks[convertInterfaceToString(linkRep.FromID)]
				link.toTask = def.tasks[convertInterfaceToString(linkRep.ToID)]
			}

			if link.toTask == nil {
				strId := strconv.Itoa(link.ID())
				return errors.New("Link[" + strId + "]: ToTask '" + convertInterfaceToString(linkRep.ToID) + "' not found")
			}

			if link.fromTask == nil {
				strId := strconv.Itoa(link.ID())
				return errors.New("Link[" + strId + "]: FromTask '" + convertInterfaceToString(linkRep.FromID) + "' not found")
			}

			// add this link as predecessor "fromLink" to the "toTask"
			link.toTask.fromLinks = append(link.toTask.fromLinks, link)

			// add this link as successor "toLink" to the "fromTask"
			link.fromTask.toLinks = append(link.fromTask.toLinks, link)

			if eh {
				def.errorHandler.links[link.id] = link
			} else {
				def.links[link.id] = link
			}
		}
	}

	return nil
}

//convertInterfaceToString will identify whether the interface is int or string and return a string in any case
//Deprecated
func convertInterfaceToString(m interface{}) string {
	if m == nil {
		panic("Invalid nil activity id found")
	}
	switch m.(type) {
	case string:
		return m.(string)
	case float64:
		return strconv.Itoa(int(m.(float64)))
	default:
		panic(fmt.Sprintf("Error parsing TaskOld with Id '%v', invalid type '%T'", m, m))
	}
}

//fixupMappings updates old mappings to new syntax
//Deprecated
func fixupMappings(mappings []*data.MappingDef) {
	for _, md := range mappings {
		if md.Type == data.MtAssign {

			val, ok := md.Value.(string)

			if ok {
				if strings.HasPrefix(val, "{T.") {
					md.Value = strings.Replace(val, "{T.", "${trigger.", 1)
				} else if strings.HasPrefix(val, "{TriggerData.") {
					md.Value = strings.Replace(val, "{TriggerData.", "${trigger.", 1)
				} else if strings.HasPrefix(val, "{A") {
					md.Value = strings.Replace(val, "{A", "${activity.", 1)
				}
			}
		}
	}
}
