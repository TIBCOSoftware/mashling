package support

import (
	//"bytes"
	//"compress/gzip"
	//"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

//const (
//	uriSchemeFile = "file://"
//	uriSchemeHttp = "http://"
//	uriSchemeRes  = "res://"
//)

// FlowManagerOld is a simple manager for flows
type FlowManagerOld struct {
	mu    *sync.Mutex // protects the flow maps
	flows map[string]*FlowEntry
}

// FlowEntry will contain either a compressed flow, an uncompressed flow or a flow uri
type FlowEntry struct {
	compressed   string
	uncompressed []byte
	uri          string
}

// NewFlowManager creates a new FlowManagerOld
func NewFlowManagerOld() *FlowManagerOld {
	return &FlowManagerOld{mu: &sync.Mutex{}, flows: make(map[string]*FlowEntry)}
}

// AddCompressed adds a compressed flow to the map of flow entries
func (mgr *FlowManagerOld) AddCompressed(id string, newFlow string) error {
	if len(newFlow) < 3 {
		return fmt.Errorf("Empty Flow with id '%s'", id)
	}
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	_, ok := mgr.flows[id]
	if ok {
		return fmt.Errorf("Flow with id '%s' already exists", id)
	}
	// Add the flow
	mgr.flows[id] = &FlowEntry{compressed: newFlow}
	logger.Debugf("Compressed flow with id '%s' added", id)
	return nil
}

// TODO add schema validation for flow
// AddUncompressed adds an uncompressed flow to the map of flow entries
func (mgr *FlowManagerOld) AddUncompressed(id string, newFlow []byte) error {
	if len(newFlow) < 3 {
		return fmt.Errorf("Empty Flow with id '%s'", id)
	}
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	_, ok := mgr.flows[id]
	if ok {
		return fmt.Errorf("Flow with id '%s' already exists", id)
	}
	// Add the flow
	mgr.flows[id] = &FlowEntry{uncompressed: newFlow}
	logger.Debugf("Uncompressed flow with id '%s' added", id)
	return nil
}

// AddURI adds a uri flow to the map of flow entries
func (mgr *FlowManagerOld) AddURI(id string, newUri string) error {
	if len(newUri) < 3 {
		return fmt.Errorf("Empty Flow URI with id '%s'", id)
	}
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	_, ok := mgr.flows[id]
	if ok {
		return fmt.Errorf("Flow with id '%s' already exists", id)
	}
	// Add the flow
	mgr.flows[id] = &FlowEntry{uri: newUri}
	logger.Debugf("URI flow with id '%s' added", id)
	return nil
}

// GetFlow gets the specified embedded flow
func (mgr *FlowManagerOld) GetFlow(id string) (*definition.DefinitionRep, error) {

	entry, ok := mgr.flows[id]

	if !ok {
		//temporary fix for tester (dynamic uri)
		if strings.HasPrefix(id, uriSchemeHttp) {
			entry = &FlowEntry{uri: id}
			mgr.flows[id] = entry
		} else {
			err := fmt.Errorf("No flow found for id '%s'", id)
			logger.Errorf(err.Error())
			return nil, err
		}
	}

	var flowDefBytes []byte

	// Uncompressed Flow condition
	if len(entry.uncompressed) > 0 {
		// Uncompressed flow
		flowDefBytes = entry.uncompressed
	}

	// Compressed Flow condition
	if len(entry.compressed) > 0 {

		decodedBytes, err := decodeAndUnzip(entry.compressed)
		if err != nil {
			decodeErr := fmt.Errorf("Error decoding compressed flow with id '%s', %s", id, err.Error())
			logger.Errorf(decodeErr.Error())
			return nil, decodeErr
		}
		flowDefBytes = decodedBytes
	}

	// URI Flow condition
	if len(entry.uri) > 0 {
		if strings.HasPrefix(entry.uri, uriSchemeFile) {
			// File URI
			logger.Infof("Loading Local Flow: %s\n", entry.uri)
			flowFilePath, _ := util.URLStringToFilePath(entry.uri)

			readBytes, err := ioutil.ReadFile(flowFilePath)
			if err != nil {
				readErr := fmt.Errorf("Error reading flow file with id '%s', %s", id, err.Error())
				logger.Errorf(readErr.Error())
				return nil, readErr
			}
			flowDefBytes = readBytes
		} else {
			// URI
			req, err := http.NewRequest("GET", entry.uri, nil)
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				getErr := fmt.Errorf("Error getting flow uri with id '%s', %s", id, err.Error())
				logger.Errorf(getErr.Error())
				return nil, getErr
			}
			defer resp.Body.Close()

			logger.Infof("response Status:", resp.Status)

			if resp.StatusCode >= 300 {
				//not found
				getErr := fmt.Errorf("Error getting flow uri with id '%s', status code %d", id, resp.StatusCode)
				logger.Errorf(getErr.Error())
				return nil, getErr
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				readErr := fmt.Errorf("Error reading flow uri response body with id '%s', %s", id, err.Error())
				logger.Errorf(readErr.Error())
				return nil, readErr
			}
			flowDefBytes = body
		}
	}

	var flow *definition.DefinitionRep
	err := json.Unmarshal(flowDefBytes, &flow)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, fmt.Errorf("Error marshalling flow with id '%s', %s", id, err.Error())
	}
	return flow, nil
}
