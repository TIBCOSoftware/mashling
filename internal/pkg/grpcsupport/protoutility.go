package grpcsupport

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes/duration"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// ConvertValues converts given string values to appropriate protobuf data types
func convertAndApplyValues(v reflect.Value, value string) {
	switch fmt.Sprintf("%s", v.Type()) {
	case "string":
		v.SetString(value)
	case "bool":
		bol, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatal("error in converting bool from value ", value, err)
		}
		v.SetBool(bol)
	case "float64":
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Fatal("error in converting float64 from value ", value, err)
		} else {
			v.SetFloat(val)
		}
	case "float32":
		val, err := strconv.ParseFloat(value, 32)
		if err != nil {
			log.Fatal("error in converting Float32 from value ", value, err)
		} else {
			v.Set(reflect.ValueOf(float32(val)))
		}
	case "int64":
		val, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			log.Fatal("error in converting Int64 from value ", value, err)
		} else {
			v.SetInt(val)
		}
	case "int32":
		val, err := strconv.ParseInt(value, 0, 32)
		if err != nil {
			log.Fatal("error in converting Int32 from value ", value, err)
		} else {
			v.Set(reflect.ValueOf(int32(val)))
		}
	case "uint64":
		val, err := strconv.ParseUint(value, 0, 64)
		if err != nil {
			log.Fatal("error in converting Uint64 from value ", value, err)
		} else {
			v.SetUint(val)
		}
	case "uint32":
		val, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			log.Fatal("error in converting Uint32 from value ", value, err)
		} else {
			v.Set(reflect.ValueOf(uint32(val)))
		}
	case "[]byte":
		val, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			log.Fatal("error in converting Bytes from value ", value, err)
		} else {
			v.SetBytes(val)
		}
	case "*timestamp.Timestamp":
		var tmStmp *timestamp.Timestamp
		err := jsonpb.UnmarshalString(value, tmStmp)
		if err != nil {
			log.Fatal("error in converting Timestamp from value", value, err)
		} else {
			v.Set(reflect.ValueOf(tmStmp))
		}
	case "*duration.Duration":
		var tmpDur *duration.Duration
		err := jsonpb.UnmarshalString(value, tmpDur)
		if err != nil {
			log.Fatal("error in converting *duration.Duration from value", value, err)
		} else {
			v.Set(reflect.ValueOf(tmpDur))
		}
	}
}

// AssignStructValues assigns given struct to values map
func AssignStructValues(structVal interface{}, values interface{}) {
	for k, v := range values.(map[string]interface{}) {
		switch k {
		case "PathParams", "Params", "QueryParams":
			mapStructWithValues(structVal, v.(map[string]string))
		case "Content":
			mapStructWithContent(structVal, v)
		}
	}
}

// mapStructWithValues assigns given struct to values map
func mapStructWithValues(structVal interface{}, values map[string]string) {
	for field, value := range values {
		temp := field[:1]
		field = strings.ToUpper(temp) + field[1:]
		v := reflect.ValueOf(structVal).Elem().FieldByName(field)
		if v.IsValid() {
			log.Printf("\n%T:%v", v.Type(), v.Type())
			convertAndApplyValues(v, value)
		}
	}
}

// mapStructWithContent assigns given struct to content data
func mapStructWithContent(structVal interface{}, content interface{}) {

	bytes, err := json.Marshal(content)
	if err != nil {
		log.Fatal("error in mapStructWithContent: ", err)
	}

	err = json.Unmarshal(bytes, structVal)
	if err != nil {
		log.Fatal("error in mapStructWithContent: ", err)
	}

}
