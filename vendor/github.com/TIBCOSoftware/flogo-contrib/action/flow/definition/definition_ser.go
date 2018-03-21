package definition

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

// DefinitionRep is a serializable representation of a flow Definition
type DefinitionRep struct {
	ExplicitReply    bool              `json:"explicitReply"`
	Name             string            `json:"name"`
	ModelID          string            `json:"model"`
	Attributes       []*data.Attribute `json:"attributes,omitempty"`
	RootTask         *TaskRep          `json:"rootTask"`
	ErrorHandlerTask *TaskRep          `json:"errorHandlerTask"`
}

// TaskRep is a serializable representation of a flow Task
type TaskRep struct {
	// Using interface{} type to support backward compatibility changes since Id was
	// int before, change to string once BC is removed
	ID          interface{} `json:"id"`
	TypeID      int         `json:"type"`
	ActivityRef string      `json:"activityRef"`
	Name        string      `json:"name"`

	Tasks []*TaskRep `json:"tasks,omitempty"`
	Links []*LinkRep `json:"links,omitempty"`

	Mappings    *Mappings              `json:"mappings,omitempty"`
	InputAttrs  map[string]interface{} `json:"input,omitempty"`
	OutputAttrs map[string]interface{} `json:"output,omitempty"`
	Settings    map[string]interface{} `json:"settings"`

	//keep temporarily for backwards compatibility
	InputAttrsOld  map[string]interface{} `json:"inputs,omitempty"`
	OutputAttrsOld map[string]interface{} `json:"outputs,omitempty"`
	InputMappings  []*data.MappingDef     `json:"inputMappings,omitempty"`
	OutputMappings []*data.MappingDef     `json:"outputMappings,omitempty"`
	Attributes     []*data.Attribute      `json:"attributes,omitempty"`
	ActivityType   string                 `json:"activityType"`
}

type Mappings struct {
	Input  []*data.MappingDef `json:"input,omitempty"`
	Output []*data.MappingDef `json:"output,omitempty"`
}

// LinkRep is a serializable representation of a flow Link
type LinkRep struct {
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

// NewDefinition creates a flow Definition from a serializable
// definition representation
func NewDefinition(rep *DefinitionRep) (def *Definition, err error) {

	defer util.HandlePanic("NewDefinition", &err)

	def = &Definition{}
	def.name = rep.Name
	def.modelID = rep.ModelID
	def.explicitReply = rep.ExplicitReply

	if len(rep.Attributes) > 0 {
		def.attrs = make(map[string]*data.Attribute, len(rep.Attributes))

		for _, value := range rep.Attributes {
			def.attrs[value.Name()] = value
		}
	}

	def.rootTask = &Task{}

	def.tasks = make(map[string]*Task)
	def.links = make(map[int]*Link)

	addTask(def, def.rootTask, rep.RootTask)
	addLinks(def, def.rootTask, rep.RootTask)

	if rep.ErrorHandlerTask != nil {
		def.ehTask = &Task{}

		addTask(def, def.ehTask, rep.ErrorHandlerTask)
		addLinks(def, def.ehTask, rep.ErrorHandlerTask)
	}

	return def, nil
}

func addTask(def *Definition, task *Task, rep *TaskRep) {
	// Workaround to support Backwards compatibility
	// Remove once rep.ID is string
	task.id = convertInterfaceToString(rep.ID)
	task.activityRef = rep.ActivityRef
	task.typeID = rep.TypeID
	task.name = rep.Name
	task.definition = def

	//temporary support for old configuration
	task.activityType = rep.ActivityType


	// Keep for now, DEPRECATE "attributes" section from flogo.json
	if len(rep.Settings) > 0 {
		task.settings = make(map[string]interface{}, len(rep.Settings))

		for name, value := range rep.Settings {
			task.settings[name] = value
		}
	}

	// create mappers
	if rep.Mappings != nil {
		if rep.Mappings.Input != nil {
			fixupMappings(rep.Mappings.Input)
			task.inputMapper = GetMapperFactory().NewTaskInputMapper(task, &MapperDef{Mappings: rep.Mappings.Input})
		}
		if rep.Mappings.Output != nil {
			task.outputMapper = GetMapperFactory().NewTaskOutputMapper(task, &MapperDef{Mappings: rep.Mappings.Output})
		}
	} else {
		//temporary support for old configuration
		if rep.InputMappings != nil {
			fixupMappings(rep.InputMappings)
			task.inputMapper = GetMapperFactory().NewTaskInputMapper(task, &MapperDef{Mappings: rep.InputMappings})
		}
		if rep.OutputMappings != nil {
			task.outputMapper = GetMapperFactory().NewTaskOutputMapper(task, &MapperDef{Mappings: rep.OutputMappings})
		}
	}

	if task.outputMapper == nil {
		task.outputMapper = GetMapperFactory().GetDefaultTaskOutputMapper(task)
	}

	// Keep for now, DEPRECATE "attributes" section from flogo.json
	if len(rep.Attributes) > 0 {
		task.inputAttrs = make(map[string]*data.Attribute, len(rep.Attributes))

		for _, value := range rep.Attributes {
			task.inputAttrs[value.Name()] = value
		}
	}

	act := activity.Get(task.activityRef)

	//todo report error if activity not registered

	if act != nil {

		inputAttrs := rep.InputAttrs

		if act.Metadata().ProducesResult || def.explicitReply {
			def.explicitReply = true
		}

		//for backwards compatibility
		if len(inputAttrs) == 0 {
			inputAttrs = rep.InputAttrsOld
		}

		if len(inputAttrs) > 0 {
			task.inputAttrs = make(map[string]*data.Attribute, len(inputAttrs))

			for name, value := range inputAttrs {

				attr := act.Metadata().Input[name]

				if attr != nil {
					//var err error
					//todo handle error
					task.inputAttrs[name], _ = data.NewAttribute(name, attr.Type(), value)
				}
			}
		}

		outputAttrs := rep.OutputAttrs

		//for backwards compatibility
		if len(outputAttrs) == 0 {
			outputAttrs = rep.OutputAttrsOld
		}

		if len(outputAttrs) > 0 {

			task.outputAttrs = make(map[string]*data.Attribute, len(outputAttrs))

			for name, value := range outputAttrs {

				attr := act.Metadata().Output[name]

				if attr != nil {
					//var err error
					//todo handle error
					task.outputAttrs[name], _ = data.NewAttribute(name, attr.Type(), value)
				}
			}
		}
	}

	def.tasks[task.id] = task
	numTasks := len(rep.Tasks)

	// flow child tasks
	if numTasks > 0 {

		for _, childTaskRep := range rep.Tasks {

			childTask := &Task{}
			childTask.parent = task
			task.tasks = append(task.tasks, childTask)
			addTask(def, childTask, childTaskRep)
		}
	}
}

//convertInterfaceToString will identify whether the interface is int or string and return a string in any case
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
		panic(fmt.Sprintf("Error parsing Task with Id '%v', invalid type '%T'", m, m))
	}
}

func addLinks(def *Definition, task *Task, rep *TaskRep) {

	numLinks := len(rep.Links)

	if numLinks > 0 {

		task.links = make([]*Link, numLinks)

		for i, linkRep := range rep.Links {

			link := &Link{}
			link.id = linkRep.ID
			//link.Parent = task
			//link.Definition = pd
			link.linkType = LinkType(linkRep.Type)
			link.value = linkRep.Value
			link.fromTask = def.tasks[convertInterfaceToString(linkRep.FromID)]
			link.toTask = def.tasks[convertInterfaceToString(linkRep.ToID)]

			// add this link as predecessor "fromLink" to the "toTask"
			link.toTask.fromLinks = append(link.toTask.fromLinks, link)

			// add this link as successor "toLink" to the "fromTask"
			link.fromTask.toLinks = append(link.fromTask.toLinks, link)

			task.links[i] = link
			def.links[link.id] = link
		}
	}
}

//fixupMappings updates old mappings to new syntax
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
