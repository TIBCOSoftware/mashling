package mapper

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

func NewAssignExpr(assignTo string, value interface{}) data.Expr {

	attrName, attrPath, _ := data.PathDeconstruct(assignTo)
	return &assignExpr{assignAttrName: attrName, assignAttrPath: attrPath, value: value}
}

type assignExpr struct {
	assignAttrName string
	assignAttrPath string
	value          interface{}
}

func (e *assignExpr) Eval(scope data.Scope) (interface{}, error) {

	var err error

	if e.assignAttrPath == "" {
		//simple assignment
		err = scope.SetAttrValue(e.assignAttrName, e.value)
		return nil, err
	}

	attr, exists := scope.GetAttr(e.assignAttrName)

	if !exists {
		return nil, fmt.Errorf("Attribute '%s' does not exists\n", e.assignAttrName)
	}

	//temporary hack
	if attr.Value() == nil {
		switch attr.Type() {
		case data.TypeObject:
			attr.SetValue(make(map[string]interface{}))
		case data.TypeParams:
			attr.SetValue(make(map[string]string))
		}
	}

	err = data.PathSetValue(attr.Value(), e.assignAttrPath, e.value)
	return nil, err
}

func NewMapperDefFromAnyArray(mappings []interface{}) (*data.MapperDef, error) {

	var mappingDefs []*data.MappingDef

	for _, mapping := range mappings {

		mappingObject := mapping.(map[string]interface{})

		mappingType, err := data.ConvertMappingType(mappingObject["type"])

		if err != nil {
			return nil, err
		}

		value := mappingObject["value"]
		mapTo := mappingObject["mapTo"].(string)

		mappingDef := &data.MappingDef{Type: data.MappingType(mappingType), MapTo: mapTo, Value: value}
		mappingDefs = append(mappingDefs, mappingDef)
	}

	return &data.MapperDef{Mappings: mappingDefs}, nil
}
