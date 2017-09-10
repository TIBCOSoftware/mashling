package cli

import (
	"sort"
	"sync"
)

var (
	toolsMu sync.Mutex
	tools   = make(map[string]*Tool)
)

// RegisterTool registers the provided tool
func RegisterTool(tool *Tool) {

	toolsMu.Lock()
	defer toolsMu.Unlock()

	if tool == nil {
		panic("RegisterTool: tool cannot be nil")
	}

	name := tool.OptionInfo().Name

	if _, exists := tools[name]; exists {
		panic("RegisterTool: tool [" + name + "] already registered")
	}

	tools[name] = tool
}

// GetTool returns the specified tool
func GetTool(toolName string) (tool *Tool, exists bool) {

	toolsMu.Lock()
	defer toolsMu.Unlock()

	tool, exists = tools[toolName]
	return tool, exists
}

// GetTools gets all the registered tools
func GetTools() []*Tool {

	toolsMu.Lock()
	defer toolsMu.Unlock()

	var ts []*Tool
	for _, v := range tools {
		ts = append(ts, v)
	}

	return ts
}

// GetToolOptionInfos gets all the OptionInfos for the
// registered tools
func GetToolOptionInfos() []*OptionInfo {

	toolsMu.Lock()
	defer toolsMu.Unlock()

	//return tools sorted by name
	var sortedKeys []string
	for k := range tools {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)
	var optionInfos []*OptionInfo
	for _, k := range sortedKeys {
		optionInfos = append(optionInfos, tools[k].optionInfo)
	}

	return optionInfos
}
