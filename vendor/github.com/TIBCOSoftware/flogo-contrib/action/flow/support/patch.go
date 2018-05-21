package support

import "github.com/TIBCOSoftware/flogo-lib/core/data"

// Patch contains a set of task patches for a Flow Patch, this
// can be used to override the default data and mappings of a Flow
type Patch struct {
	TaskPatches []*TaskPatch `json:"tasks"` //put in mapper object

	taskPatchMap map[string]*TaskPatch
}

// Init initializes the FlowPatch, usually called after deserialization
func (pp *Patch) Init() {

	numAttrs := len(pp.TaskPatches)
	if numAttrs > 0 {

		pp.taskPatchMap = make(map[string]*TaskPatch, numAttrs)

		for _, patch := range pp.TaskPatches {
			pp.taskPatchMap[patch.ID] = patch
		}
	}
}

// GetPatch returns the Task Patch for the specified task (referred to by ID)
func (pp *Patch) GetPatch(taskID string) *TaskPatch {
	return pp.taskPatchMap[taskID]
}

// GetInputMapper returns the InputMapper for the specified task (referred to by ID)
func (pp *Patch) GetInputMapper(taskID string) data.Mapper {
	taskPatch, exists := pp.taskPatchMap[taskID]

	if exists {
		return taskPatch.InputMapper()
	}

	return nil
}

// GetOutputMapper returns the OutputMapper for the specified task (referred to by ID)
func (pp *Patch) GetOutputMapper(taskID string) data.Mapper {
	taskPatch, exists := pp.taskPatchMap[taskID]

	if exists {
		return taskPatch.OutputMapper()
	}

	return nil
}

// TaskPatch contains patching information for a Task, such has attributes,
// input mappings, output mappings.  This is used to override the corresponding
// settings for a Task in the Process
type TaskPatch struct {
	ID             string             `json:"id"`
	Attributes     []*data.Attribute  `json:"attributes"`
	InputMappings  []*data.MappingDef `json:"inputMappings"` //put in mapper object
	OutputMappings []*data.MappingDef `json:"ouputMappings"` //put in mapper object

	Attrs        map[string]*data.Attribute
	inputMapper  data.Mapper
	outputMapper data.Mapper
}

// InputMapper returns the overriding InputMapper
func (tp *TaskPatch) InputMapper() data.Mapper {
	return tp.inputMapper
}

// OutputMapper returns the overriding OutputMapper
func (tp *TaskPatch) OutputMapper() data.Mapper {
	return tp.outputMapper
}
