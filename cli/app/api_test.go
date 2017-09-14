package app

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const gatewayJSON string = ``
const expectedFlogoJSON string = ``

func TestGetGatewayDetails(t *testing.T) {
	//Need to create gateway project under temp folder and pass the same for testing.
	_, err := GetGatewayDetails(SetupNewProjectEnv(), ALL)
	if err != nil {
		t.Error("Error while getting gateway details")
	}
}

func TestTranslateGatewayJSON2FlogoJSON(t *testing.T) {

	flogoJSON, err := TranslateGatewayJSON2FlogoJSON(gatewayJSON)

	if err != nil {
		t.Error("Error in TranslateGatewayJSON2FlogoJSON function. err: ", err)
	}
	isEqualJSON, err := AreEqualJSON(expectedFlogoJSON, flogoJSON)
	assert.NoError(t, err, "Error: Error comparing expected and actual Flogo JSON %v", err)
	assert.True(t, isEqualJSON, "Error: Expected and actual Flogo JSON contents are not equal.")
}

func AreEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}
