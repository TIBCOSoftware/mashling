package v1

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/mashling/internal/app/gateway/flogo/registry"
	"github.com/mashling/commons/lib/types"
	"github.com/TIBCOSoftware/mashling/pkg/strings"
)

// LoadFlogoFlow loads a flogo flow reference (remote or local file) and returns the raw JSON
func LoadFlogoFlow(reference string) (action json.RawMessage, err error) {
	var fileData []byte
	// Check for which supported flow location
	if strings.HasPrefix(reference, "github.com") || strings.HasPrefix(reference, "http://") || strings.HasPrefix(reference, "https://") {
		// Load remote JSON. Special handling for Github
		if strings.HasPrefix(reference, "github.com") {
			// Special handling for # specifying a branch, tag, or revision to use.
			var ref string
			if strings.Contains(reference, "#") {
				split := strings.Split(reference, "#")
				reference = split[0]
				commit := split[1]
				ref = util.ReplaceNth(reference, "/", "/"+commit+"/", 3)
			} else {
				ref = util.ReplaceNth(reference, "/", "/master/", 3)
			}

			reference = strings.Replace(ref, "github.com", githubRawContent, 1)
		}
		response, rerr := http.Get(reference)
		if rerr != nil {
			return action, rerr
		}
		fileData, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return action, err
		}
	} else if strings.HasPrefix(reference, "file://") {
		// Local file reference, load from disk, if relative path load relative to binary
		reference = strings.Replace(reference, "file://", "", 1)
		fileData, err = ioutil.ReadFile(filepath.FromSlash(reference))
		if err != nil {
			return action, err
		}
	} else {
		if err != nil {
			return action, errors.New("unsupported flogo flow reference format")
		}
	}
	var flogoFlowDef *app.Config
	err = json.Unmarshal(json.RawMessage(string(fileData)), &flogoFlowDef)
	if err != nil {
		return action, err
	}

	actions := flogoFlowDef.Actions
	if len(actions) != 1 {
		return action, errors.New("Please make sure each referenced flow has only one action")
	}

	action, err = json.Marshal(actions[0])
	if err != nil {
		return action, err
	}
	return action, err
}

//IdentifyMissingFlogoDependencies returns all missing flogo refs from the action.
func IdentifyMissingFlogoDependencies(actionDefinition json.RawMessage) (deps []string, err error) {
	action := &types.FlogoAction{}
	err = json.Unmarshal(actionDefinition, action)
	if err != nil {
		return deps, err
	}
	// Only support action/flow definitions right now.
	flowDefRep := &definition.DefinitionRep{}
	rawActionData, err := json.Marshal(action.Data)
	if err != nil {
		return deps, err
	}
	extract := make(map[string]interface{})
	err = json.Unmarshal(rawActionData, &extract)
	if err != nil {
		return deps, err
	}
	if _, ok := extract["flow"]; !ok {
		return deps, errors.New("no flow element in embedded flogo flow")
	}
	flow, ok := extract["flow"].(map[string]interface{})
	if !ok {
		return deps, errors.New("flow of embedded flogo flow is not JSON")
	}
	rawFlowData, err := json.Marshal(flow)
	if err != nil {
		return deps, err
	}
	err = json.Unmarshal(rawFlowData, flowDefRep)
	if err != nil {
		return deps, err
	}
	for _, task := range flowDefRep.RootTask.Tasks {
		if _, exists := registry.SupportedImports[task.ActivityRef]; !exists {
			deps = append(deps, task.ActivityRef)
		}
	}
	return deps, err
}
