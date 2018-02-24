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
	"strings"
)

var (
	// ErrorJSONRequired json body required error
	ErrorJSONRequired = errors.New("json body required")
	// ErrorXMLRequired xml body required error
	ErrorXMLRequired = errors.New("xml body required")
	// ErrorPointerRequired produced when v is not a pointer
	ErrorPointerRequired = errors.New("v must be a pointer")
)

// XMLUnmarshal parases XML and stores it in a JSON like data structure
func XMLUnmarshal(data []byte, v interface{}) error {
	output := reflect.ValueOf(v)
	if output.Kind() != reflect.Ptr {
		return ErrorPointerRequired
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

func getString(key string, v map[string]interface{}) (string, error) {
	if x, ok := v[key]; ok {
		if y, ok := x.(string); ok {
			return y, nil
		}
		return "", fmt.Errorf("%s should be a string", key)
	}
	return "", nil
}

// XMLMarshal take a XML map and generates XML
func XMLMarshal(v interface{}) ([]byte, error) {
	var input map[string]interface{}
	if x, ok := v.(map[string]interface{}); ok {
		input = x
	} else {
		return nil, errors.New("XML map required")
	}

	output := bytes.Buffer{}
	var unparse func(v map[string]interface{}) error
	processBody := func(v map[string]interface{}) error {
		if x, ok := v["_body"]; ok {
			if y, ok := x.([]interface{}); ok {
				for _, z := range y {
					if item, ok := z.(map[string]interface{}); ok {
						if err := unparse(item); err != nil {
							return err
						}
					} else {
						return errors.New("item should be a map")
					}
				}
			} else {
				return errors.New("body should be a slice")
			}
		}
		return nil
	}
	unparse = func(v map[string]interface{}) error {
		typ, err := getString("_type", v)
		if err != nil {
			return err
		}

		switch typ {
		case "ProcInst":
			var target, inst string
			target, err = getString("_target", v)
			if err != nil {
				return err
			}
			inst, err = getString("_inst", v)
			if err != nil {
				return err
			}
			output.WriteString("<?")
			output.WriteString(target)
			output.WriteString(" ")
			output.WriteString(inst)
			output.WriteString("?>")
		case "Comment":
			var body string
			body, err = getString("_body", v)
			if err != nil {
				return err
			}
			output.WriteString("<!--")
			output.WriteString(body)
			output.WriteString("-->")
		case "Element":
			var space, name string
			space, err = getString("_space", v)
			if err != nil {
				return err
			}
			name, err = getString("_name", v)
			if err != nil {
				return err
			}
			output.WriteString("<")
			if space != "" {
				output.WriteString(space)
				output.WriteString(":")
			}
			output.WriteString(name)
			for key, value := range v {
				switch key {
				case "_type", "_target", "_inst", "_body", "_space", "_name":
				default:
					if x, ok := value.(string); ok {
						parts := strings.Split(key, "___")
						output.WriteString(" ")
						switch len(parts) {
						case 1:
							output.WriteString(parts[0])
						case 2:
							output.WriteString(parts[0])
							output.WriteString(":")
							output.WriteString(parts[1])
						default:
							return errors.New("invalid attribute")
						}
						output.WriteString("=")
						output.WriteString("\"")
						output.WriteString(x)
						output.WriteString("\"")
					} else {
						return errors.New("attribute should be a string")
					}
				}
			}
			output.WriteString(">")
			err = processBody(v)
			if err != nil {
				return err
			}
			output.WriteString("</")
			if space != "" {
				output.WriteString(space)
				output.WriteString(":")
			}
			output.WriteString(name)
			output.WriteString(">")
		case "CharData":
			var body string
			body, err = getString("_body", v)
			if err != nil {
				return err
			}
			output.WriteString(body)
		default:
			err = processBody(v)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := unparse(input)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

// Parse parses JSON or XML with a give MIME type into an interface
func Parse(mime string, data []byte, v interface{}) error {
	output := reflect.ValueOf(v)
	if output.Kind() != reflect.Ptr {
		return ErrorPointerRequired
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
			mime = "application/json"
			break
		}
		err = XMLUnmarshal(data, v)
		if err == nil {
			mime = "application/xml"
		}
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
			y["___copy___"] = string(data)
		}
	}
	output.Elem().Set(reflect.ValueOf(x))

	return nil
}

// Clean removed meta keys from the map
func Clean(input map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	for key, value := range input {
		switch key {
		case "___mime___", "___copy___":
		default:
			output[key] = value
		}
	}
	return output
}

var jsonMarshal = json.Marshal

// Unparse generates a string from a map
func Unparse(v interface{}) ([]byte, error) {
	var input map[string]interface{}
	if x, ok := v.(map[string]interface{}); ok {
		input = x
	} else {
		return jsonMarshal(v)
	}

	mime, err := getString("___mime___", input)
	if err != nil {
		return nil, err
	}
	switch mime {
	case "application/json":
		return jsonMarshal(Clean(input))
	case "text/xml", "application/xml":
		return XMLMarshal(Clean(input))
	}

	cp, err := getString("___copy___", input)
	if err != nil {
		return nil, err
	}
	if cp == "" {
		return jsonMarshal(Clean(input))
	}
	return []byte(cp), nil
}
