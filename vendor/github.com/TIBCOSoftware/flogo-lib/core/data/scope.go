package data

import (
	"errors"
	"sync"
	"fmt"
)

//func init() {
//	SetResolver(RES_SCOPE, Resolve)
//}

// Resolve will resolve a value in the given scope
//func Resolve(scope Scope, value string) (interface{}, bool) {
//	attr, ok := scope.GetAttr(value)
//	if !ok {
//		return nil, false
//	}
//
//	return attr.Value, true
//}

// Scope is a set of attributes that are accessible
type Scope interface {
	// GetAttr gets the specified attribute
	GetAttr(name string) (attr *Attribute, exists bool)

	// SetAttrValue sets the value of the specified attribute
	SetAttrValue(name string, value interface{}) error
}

// MutableScope is a scope that new attributes can be added
type MutableScope interface {
	Scope

	//AddAttr adds an attribute to the scope
	AddAttr(name string, valueType Type, value interface{}) *Attribute
}

// SimpleScope is a basic implementation of a scope
type SimpleScope struct {
	parentScope Scope
	attrs       map[string]*Attribute
}

// NewSimpleScope creates a new SimpleScope
func NewSimpleScope(attrs []*Attribute, parentScope Scope) Scope {

	return newSimpleScope(attrs, parentScope)
}

// NewSimpleScope creates a new SimpleScope
func newSimpleScope(attrs []*Attribute, parentScope Scope) *SimpleScope {

	scope := &SimpleScope{
		parentScope: parentScope,
		attrs:       make(map[string]*Attribute),
	}

	for _, attr := range attrs {
		scope.attrs[attr.Name()] = attr
	}

	return scope
}

// NewSimpleScopeFromMap creates a new SimpleScope
func NewSimpleScopeFromMap(attrs map[string]*Attribute, parentScope Scope) *SimpleScope {

	scope := &SimpleScope{
		parentScope: parentScope,
		attrs:       attrs,
	}

	return scope
}

// GetAttr implements Scope.GetAttr
func (s *SimpleScope) GetAttr(name string) (attr *Attribute, exists bool) {

	attr, found := s.attrs[name]

	if found {
		return attr, true
	}

	if s.parentScope != nil {
		return s.parentScope.GetAttr(name)
	}

	return nil, false
}

// SetAttrValue implements Scope.SetAttrValue
func (s *SimpleScope) SetAttrValue(name string, value interface{}) error {

	attr, found := s.attrs[name]

	if found {
		attr.SetValue(value)
		return nil
	}

	return errors.New("attribute not in scope")
}

// AddAttr implements MutableScope.AddAttr
func (s *SimpleScope) AddAttr(name string, valueType Type, value interface{}) *Attribute {

	attr, found := s.attrs[name]

	if found {
		attr.SetValue(value)
	} else {
		//todo handle error, add error to AddAttr signature
		attr, _ = NewAttribute(name, valueType, value)
		s.attrs[name] = attr
	}

	return attr
}

// SimpleSyncScope is a basic implementation of a synchronized scope
type SimpleSyncScope struct {
	scope *SimpleScope
	mutex sync.Mutex
}

// NewSimpleSyncScope creates a new SimpleSyncScope
func NewSimpleSyncScope(attrs []*Attribute, parentScope Scope) MutableScope {

	var syncScope SimpleSyncScope
	syncScope.scope = newSimpleScope(attrs, parentScope)

	return &syncScope
}

// GetAttr implements Scope.GetAttr
func (s *SimpleSyncScope) GetAttr(name string) (value *Attribute, exists bool) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.scope.GetAttr(name)
}

// SetAttrValue implements Scope.SetAttrValue
func (s *SimpleSyncScope) SetAttrValue(name string, value interface{}) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.scope.SetAttrValue(name, value)
}

// AddAttr implements MutableScope.AddAttr
func (s *SimpleSyncScope) AddAttr(name string, valueType Type, value interface{}) *Attribute {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.scope.AddAttr(name, valueType, value)
}

var (
	globalScope = NewSimpleSyncScope(nil, nil)
)

// GetGlobalScope gets the global scope the application
func GetGlobalScope() MutableScope {
	return globalScope
}

// FixedScope is an implementation of a empty scope fixed to a particular set of metadata
type FixedScope struct {
	attrs    map[string]*Attribute
	metadata map[string]*Attribute
}

// NewFixedScope creates a new SimpleScope
func NewFixedScope(metadata map[string]*Attribute) *FixedScope {

	scope := &FixedScope{
		metadata: make(map[string]*Attribute),
		attrs:    make(map[string]*Attribute),
	}

	scope.metadata = metadata

	return scope
}

func NewFixedScopeFromMap(metadata map[string]*Attribute) *FixedScope {

	scope := &FixedScope{
		metadata: metadata,
		attrs:    make(map[string]*Attribute),
	}

	return scope
}

// GetAttr implements Scope.GetAttr
func (s *FixedScope) GetAttr(name string) (attr *Attribute, exists bool) {

	attr, found := s.attrs[name]

	if found {
		return attr, true
	} else {
		metaAttr, found := s.metadata[name]
		if found {
			attr, _ := NewAttribute(name, metaAttr.Type(), metaAttr.value)
			s.attrs[name] = attr
			return attr, true
		}
	}
	return nil, false
}

// GetAttrs gets the attributes set in the scope
func (s *FixedScope) GetAttrs() map[string]*Attribute {
	return s.attrs
}

// SetAttrValue implements Scope.SetAttrValue
func (s *FixedScope) SetAttrValue(name string, value interface{}) error {

	attr, found := s.attrs[name]

	if found {
		attr.SetValue(value)
		return nil
	} else {
		metaAttr, found := s.metadata[name]
		if found {
			attr, err := NewAttribute(name, metaAttr.Type(), value)
			s.attrs[name] = attr
			return err
		}
	}

	return fmt.Errorf("attribute '%s' not in scope", name)
}

//todo fix up all the scopes!

// FixedScope is an implementation of a empty scope fixed to a particular set of metadata
type FlexableScope struct {
	attrs    map[string]*Attribute
	metadata map[string]*Attribute
}

// NewFlexableScope creates a new SimpleScope
func NewFlexableScope(metadata map[string]*Attribute) *FlexableScope {

	scope := &FlexableScope{
		metadata: make(map[string]*Attribute),
		attrs:    make(map[string]*Attribute),
	}

	scope.metadata = metadata

	return scope
}

func NewFlexableScopeFromMap(metadata map[string]*Attribute) *FlexableScope {

	scope := &FlexableScope{
		metadata: metadata,
		attrs:    make(map[string]*Attribute),
	}

	return scope
}

// GetAttr implements Scope.GetAttr
func (s *FlexableScope) GetAttr(name string) (attr *Attribute, exists bool) {

	attr, found := s.attrs[name]

	if found {
		return attr, true
	} else {
		metaAttr, found := s.metadata[name]
		if found {
			attr, _ := NewAttribute(name, metaAttr.Type(), metaAttr.value)
			s.attrs[name] = attr
			return attr, true
		}
	}
	return nil, false
}

// GetAttrs gets the attributes set in the scope
func (s *FlexableScope) GetAttrs() map[string]*Attribute {
	return s.attrs
}

// SetAttrValue implements Scope.SetAttrValue
func (s *FlexableScope) SetAttrValue(name string, value interface{}) error {

	attr, found := s.attrs[name]

	if found {
		attr.SetValue(value)
		return nil
	} else {
		metaAttr, found := s.metadata[name]
		if found {
			attr, err := NewAttribute(name, metaAttr.Type(), value)
			s.attrs[name] = attr
			return err
		}
	}

	t, err := GetType(value)
	if err != nil {
		t = TypeAny
	}

	s.attrs[name], err = NewAttribute(name, t, value)

	return err
}
