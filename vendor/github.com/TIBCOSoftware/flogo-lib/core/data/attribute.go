package data

import (
	"encoding/json"
	"errors"
)

// Attribute is a simple structure used to define a data Attribute/property
type Attribute struct {
	name     string
	dataType Type
	value    interface{}
}

// NewAttribute constructs a new attribute
func NewAttribute(name string, dataType Type, value interface{}) (*Attribute, error) {

	var attr Attribute
	attr.name = name
	attr.dataType = dataType
	attr.value = value

	var err error
	attr.value, err = CoerceToValue(value, dataType)

	return &attr, err
}

// NewZeroAttribute constructs a new attribute
func NewZeroAttribute(name string, dataType Type) *Attribute {

	var attr Attribute
	attr.name = name
	attr.dataType = dataType
	attr.value, _ = CoerceToValue(nil, dataType)

	return &attr
}

// CloneAttribute clones the given attribute assigning a new name
func CloneAttribute(name string, oldAttr *Attribute) *Attribute {
	var attr Attribute
	attr.name = name
	attr.dataType = oldAttr.dataType
	attr.value = oldAttr.value

	return &attr
}

func (a *Attribute) Name() string {
	return a.name
}

func (a *Attribute) Type() Type {
	return a.dataType
}

func (a *Attribute) Value() interface{} {
	return a.value
}

func (a *Attribute) SetValue(value interface{}) (err error) {
	a.value, err = CoerceToValue(value, a.dataType)
	return err
}

// MarshalJSON implements json.Marshaler.MarshalJSON
func (a *Attribute) MarshalJSON() ([]byte, error) {

	return json.Marshal(&struct {
		Name  string      `json:"name"`
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
	}{
		Name:  a.name,
		Type:  a.dataType.String(),
		Value: a.value,
	})
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON
func (a *Attribute) UnmarshalJSON(data []byte) error {

	ser := &struct {
		Name  string      `json:"name"`
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
	}{}

	if err := json.Unmarshal(data, ser); err != nil {
		return err
	}

	a.name = ser.Name
	dt, exists := ToTypeEnum(ser.Type)

	if !exists {
		return errors.New("unknown data type: " + ser.Type)
	}
	a.dataType = dt

	val, err := CoerceToValue(ser.Value, a.dataType)

	if err != nil {
		return err
	} else {
		a.value = val
	}

	return nil
}

// ComplexObject is the value that is used when using a "COMPLEX_OBJECT" type
type ComplexObject struct {
	Metadata string      `json:"metadata"`
	Value    interface{} `json:"value"`
}

type IOMetadata struct {
	Input  map[string]*Attribute
	Output map[string]*Attribute
}

func (md *IOMetadata) UnmarshalJSON(b []byte) error {

	ser := &struct {
		Input  []*Attribute `json:"input"`
		Output []*Attribute `json:"output"`
	}{}

	if err := json.Unmarshal(b, ser); err != nil {
		return err
	}

	md.Input = make(map[string]*Attribute, len(ser.Input))
	md.Output = make(map[string]*Attribute, len(ser.Output))

	for _, attr := range ser.Input {
		md.Input[attr.Name()] = attr
	}

	for _, attr := range ser.Output {
		md.Output[attr.Name()] = attr
	}

	return nil
}

func (md *IOMetadata) MarshalJSON() ([]byte, error) {
	var mdInputs []*Attribute
	var mdOutputs []*Attribute

	for _, v := range md.Input {
		mdInputs = append(mdInputs, v)
	}
	for _, v := range md.Output {
		mdOutputs = append(mdOutputs, v)
	}

	return json.Marshal(&struct {
		Input  []*Attribute `json:"input"`
		Output []*Attribute `json:"output"`
	}{
		Input:  mdInputs,
		Output: mdOutputs,
	})
}
