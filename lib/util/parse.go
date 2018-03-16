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
	// ErrorXMLMismatchedElements XML start and end elements dont' match
	ErrorXMLMismatchedElements = errors.New("mismatched elements")
	// ErrorXMLMapRequired a XML map is required
	ErrorXMLMapRequired = errors.New("XML map required")
	// ErrorXMLItemNotMap a XML item is not a map
	ErrorXMLItemNotMap = errors.New("item should be a map")
	// ErrorXMLBodyNotSlice a XML body is not a slice
	ErrorXMLBodyNotSlice = errors.New("body should be a slice")
	// ErrorXMLInvalidAttribute XML attribute is invalid
	ErrorXMLInvalidAttribute = errors.New("invalid attribute")
	// ErrorXMLAttributeNotString XML attribute is not a string
	ErrorXMLAttributeNotString = errors.New("attribute should be a string")
	// ErrorInvalidXML XML is invalid
	ErrorInvalidXML = errors.New("invalid XML")
)

const (
	// MIMEApplicationJSON the JSON MIME type
	MIMEApplicationJSON = "application/json"
	// MIMEApplicationJSONUTF8 the UTF8 JSON MIME type
	MIMEApplicationJSONUTF8 = "application/json; charset=UTF-8"
	// MIMETextXML a XML MIME type
	MIMETextXML = "text/xml"
	// MIMEApplicationXML a XML MIME type
	MIMEApplicationXML = "application/xml"
	// MIMEUnknown is an unknown MIME type
	MIMEUnknown = "application/octet-stream"
	// MIMEForm is a form MIME type
	MIMEForm = "application/x-www-form-urlencoded"

	// MetaMIME the meta MIME key
	MetaMIME = "___mime___"
	// MetaCopy the meta copy key
	MetaCopy = "___copy___"

	// XMLKeyType is the key for the XML type
	XMLKeyType = "_type"
	// XMLKeyTarget is the key for the processing target
	XMLKeyTarget = "_target"
	// XMLKeyInst is the key for the processing instruction
	XMLKeyInst = "_inst"
	// XMLKeyBody is the key for the XML body
	XMLKeyBody = "_body"
	// XMLKeySpace is the key for the XML namespace
	XMLKeySpace = "_space"
	// XMLKeyName is the key for the XML name
	XMLKeyName = "_name"

	// XMLTypeProcInst is a XML processing instruction
	XMLTypeProcInst = "ProcInst"
	// XMLTypeComment is a XML comment
	XMLTypeComment = "Comment"
	// XMLTypeElement is a XML element
	XMLTypeElement = "Element"
	// XMLTypeCharData is XML char data
	XMLTypeCharData = "CharData"

	// XMLAttributeSep XML attribute separator
	XMLAttributeSep = "___"
)

// XMLUnmarshal parases XML and stores it in a JSON like data structure
func XMLUnmarshal(data []byte, v interface{}) error {
	output := reflect.ValueOf(v)
	if output.Kind() != reflect.Ptr {
		return ErrorPointerRequired
	}

	buffer := bytes.NewReader(data)
	decoder, valid := xml.NewDecoder(buffer), false
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
				data[XMLKeyType] = XMLTypeProcInst
				data[XMLKeyTarget] = t.Target
				data[XMLKeyInst] = string(t.Inst)
				children := x[XMLKeyBody].([]interface{})
				x[XMLKeyBody] = append(children, data)
			case xml.Comment:
				data := make(map[string]interface{})
				data[XMLKeyType] = XMLTypeComment
				data[XMLKeyBody] = string(t)
				children := x[XMLKeyBody].([]interface{})
				x[XMLKeyBody] = append(children, data)
			case xml.StartElement:
				valid = true
				child := make(map[string]interface{})
				child[XMLKeyType] = XMLTypeElement
				child[XMLKeySpace] = t.Name.Space
				child[XMLKeyName] = t.Name.Local
				for _, attr := range t.Attr {
					name := attr.Name.Local
					if attr.Name.Space != "" {
						name = attr.Name.Space + XMLAttributeSep + name
					}
					child[name] = attr.Value
				}
				child[XMLKeyBody] = make([]interface{}, 0)
				err = parse(t.Name, child)
				if err != nil {
					return err
				}
				children := x[XMLKeyBody].([]interface{})
				x[XMLKeyBody] = append(children, child)
			case xml.CharData:
				data := make(map[string]interface{})
				data[XMLKeyType] = XMLTypeCharData
				data[XMLKeyBody] = string(t)
				children := x[XMLKeyBody].([]interface{})
				x[XMLKeyBody] = append(children, data)
			case xml.EndElement:
				if name != t.Name {
					return ErrorXMLMismatchedElements
				}
				return nil
			}
		}
	}

	x := make(map[string]interface{})
	x[XMLKeyBody] = make([]interface{}, 0)
	err := parse(xml.Name{}, x)
	if err != nil && err != io.EOF {
		return err
	}
	if !valid {
		return ErrorInvalidXML
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
		return nil, ErrorXMLMapRequired
	}

	output := bytes.Buffer{}
	var unparse func(v map[string]interface{}) error
	processBody := func(v map[string]interface{}) error {
		if x, ok := v[XMLKeyBody]; ok {
			if y, ok := x.([]interface{}); ok {
				for _, z := range y {
					if item, ok := z.(map[string]interface{}); ok {
						if err := unparse(item); err != nil {
							return err
						}
					} else {
						return ErrorXMLItemNotMap
					}
				}
			} else {
				return ErrorXMLBodyNotSlice
			}
		}
		return nil
	}
	unparse = func(v map[string]interface{}) error {
		typ, err := getString(XMLKeyType, v)
		if err != nil {
			return err
		}

		switch typ {
		case XMLTypeProcInst:
			var target, inst string
			target, err = getString(XMLKeyTarget, v)
			if err != nil {
				return err
			}
			inst, err = getString(XMLKeyInst, v)
			if err != nil {
				return err
			}
			output.WriteString("<?")
			output.WriteString(target)
			output.WriteString(" ")
			output.WriteString(inst)
			output.WriteString("?>")
		case XMLTypeComment:
			var body string
			body, err = getString(XMLKeyBody, v)
			if err != nil {
				return err
			}
			output.WriteString("<!--")
			output.WriteString(body)
			output.WriteString("-->")
		case XMLTypeElement:
			var space, name string
			space, err = getString(XMLKeySpace, v)
			if err != nil {
				return err
			}
			name, err = getString(XMLKeyName, v)
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
				case XMLKeyType, XMLKeyTarget, XMLKeyInst, XMLKeyBody, XMLKeySpace, XMLKeyName:
				default:
					if x, ok := value.(string); ok {
						parts := strings.Split(key, XMLAttributeSep)
						output.WriteString(" ")
						switch len(parts) {
						case 1:
							output.WriteString(parts[0])
						case 2:
							output.WriteString(parts[0])
							output.WriteString(":")
							output.WriteString(parts[1])
						default:
							return ErrorXMLInvalidAttribute
						}
						output.WriteString("=")
						output.WriteString("\"")
						output.WriteString(x)
						output.WriteString("\"")
					} else {
						return ErrorXMLAttributeNotString
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
		case XMLTypeCharData:
			var body string
			body, err = getString(XMLKeyBody, v)
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

// Unmarshal parses JSON or XML with a give MIME type into an interface
func Unmarshal(mime string, data []byte, v interface{}) error {
	output := reflect.ValueOf(v)
	if output.Kind() != reflect.Ptr {
		return ErrorPointerRequired
	}

	parsed := false
	switch mime {
	case MIMEApplicationJSON, MIMEApplicationJSONUTF8:
		if len(data) == 0 {
			return ErrorJSONRequired
		}
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.UseNumber()
		err := decoder.Decode(v)
		if err != nil {
			return err
		}
		parsed = true
	case MIMETextXML, MIMEApplicationXML:
		if len(data) == 0 {
			return ErrorXMLRequired
		}
		err := XMLUnmarshal(data, v)
		if err != nil {
			return err
		}
		parsed = true
	case MIMEForm, "":
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.UseNumber()
		err := decoder.Decode(v)
		if err == nil {
			mime = MIMEApplicationJSON
			parsed = true
			break
		}
		err = XMLUnmarshal(data, v)
		if err == nil {
			mime = MIMEApplicationXML
			parsed = true
			break
		}
		mime = MIMEUnknown
		fallthrough
	default:
		output.Elem().Set(reflect.ValueOf(make(map[string]interface{})))
	}

	if y, ok := output.Elem().Interface().(map[string]interface{}); ok {
		if mime != MIMEApplicationJSON && mime != MIMEApplicationJSONUTF8 {
			y[MetaMIME] = mime
		}
		if !parsed {
			y[MetaCopy] = string(data)
		}
	}

	return nil
}

// Clean removed meta keys from the map
func Clean(input map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	for key, value := range input {
		switch key {
		case MetaMIME, MetaCopy:
		default:
			output[key] = value
		}
	}
	return output
}

// Marshal generates a string from a map
func Marshal(v interface{}) ([]byte, error) {
	var input map[string]interface{}
	if x, ok := v.(map[string]interface{}); ok {
		input = x
	} else {
		return json.MarshalIndent(v, "", " ")
	}

	mime, err := getString(MetaMIME, input)
	if err != nil {
		return nil, err
	}
	switch mime {
	case MIMEApplicationJSON, MIMEApplicationJSONUTF8, "":
		return json.MarshalIndent(Clean(input), "", " ")
	case MIMETextXML, MIMEApplicationXML:
		return XMLMarshal(Clean(input))
	}

	cp, err := getString(MetaCopy, input)
	if err != nil {
		return nil, err
	}
	return []byte(cp), nil
}
