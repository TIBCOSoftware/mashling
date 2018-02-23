/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package util

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
)

var (
	// ErrorJSONRequired json body required error
	ErrorJSONRequired = errors.New("json body required")
	// ErrorXMLRequired xml body required error
	ErrorXMLRequired = errors.New("xml body required")
)

// XMLUnmarshal parases XML and stores it in a JSON like data structure
func XMLUnmarshal(data []byte, v interface{}) error {
	output := reflect.ValueOf(v)
	if output.Kind() != reflect.Ptr {
		return errors.New("v must be a pointer")
	}

	buffer := bytes.NewBuffer(data)
	decoder := xml.NewDecoder(buffer)
	var parse func(name xml.Name, x map[string]interface{}) error
	parse = func(name xml.Name, x map[string]interface{}) error {
		for {
			token, err := decoder.RawToken()
			if err != nil {
				return err
			}
			switch t := token.(type) {
			case xml.ProcInst:
				data := make(map[string]interface{})
				data["_type"] = "ProcInst"
				data["_target"] = t.Target
				data["_inst"] = string(t.Inst)
				children := x["_body"].([]interface{})
				x["_body"] = append(children, data)
			case xml.Comment:
				data := make(map[string]interface{})
				data["_type"] = "Comment"
				data["_body"] = string(t)
				children := x["_body"].([]interface{})
				x["_body"] = append(children, data)
			case xml.StartElement:
				child := make(map[string]interface{})
				child["_type"] = "Element"
				child["_space"] = t.Name.Space
				child["_name"] = t.Name.Local
				for _, attr := range t.Attr {
					name := attr.Name.Local
					if attr.Name.Space != "" {
						name = attr.Name.Space + "___" + name
					}
					child[name] = attr.Value
				}
				child["_body"] = make([]interface{}, 0)
				err = parse(t.Name, child)
				if err != nil {
					return err
				}
				children := x["_body"].([]interface{})
				x["_body"] = append(children, child)
			case xml.CharData:
				data := make(map[string]interface{})
				data["_type"] = "CharData"
				data["_body"] = string(t)
				children := x["_body"].([]interface{})
				x["_body"] = append(children, data)
			case xml.EndElement:
				if name != t.Name {
					return errors.New("mismatched elements")
				}
				return nil
			default:
				fmt.Println(t)
			}
		}
	}

	x := make(map[string]interface{})
	x["_body"] = make([]interface{}, 0)
	err := parse(xml.Name{}, x)
	if err != nil && err != io.EOF {
		return err
	}

	output.Elem().Set(reflect.ValueOf(x))

	return nil
}

// Parse parses JSON or XML with a give MIME type into an interface
func Parse(mime string, data []byte, v interface{}) error {
	output := reflect.ValueOf(v)
	if output.Kind() != reflect.Ptr {
		return errors.New("v must be a pointer")
	}

	switch mime {
	case "application/json":
		if len(data) == 0 {
			return ErrorJSONRequired
		}
		err := json.Unmarshal(data, v)
		if err != nil {
			return err
		}
	case "text/xml", "application/xml":
		if len(data) == 0 {
			return ErrorXMLRequired
		}
		err := XMLUnmarshal(data, v)
		if err != nil {
			return err
		}
	case "":
		err := json.Unmarshal(data, v)
		if err == nil {
			break
		}
		XMLUnmarshal(data, v)
	}

	x := output.Elem().Interface()
	if x == nil {
		x = make(map[string]interface{})
	}
	if y, ok := x.(map[string]interface{}); ok {
		if mime != "" {
			y["___mime___"] = mime
		}
		if len(data) > 0 {
			y["___orig___"] = string(data)
		}
	}
	output.Elem().Set(reflect.ValueOf(x))

	return nil
}
