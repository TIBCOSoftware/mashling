package bson

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var defaultStructCodec = &StructCodec{
	cache:  make(map[reflect.Type]*structDescription),
	parser: DefaultStructTagParser,
}

// StructCodec is the Codec used for struct values.
type StructCodec struct {
	cache  map[reflect.Type]*structDescription
	l      sync.RWMutex
	parser StructTagParser
}

var _ Codec = &StructCodec{}

// NewStructCodec returns a StructCodec that uses p for struct tag parsing.
func NewStructCodec(p StructTagParser) (*StructCodec, error) {
	if p == nil {
		return nil, errors.New("a StructTagParser must be provided to NewStructCodec")
	}

	return &StructCodec{
		cache:  make(map[reflect.Type]*structDescription),
		parser: p,
	}, nil
}

// EncodeValue handles encoding generic struct types.
func (sc *StructCodec) EncodeValue(r EncodeContext, vw ValueWriter, i interface{}) error {
	val := reflect.ValueOf(i)
	for {
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
			continue
		}

		break
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("%T can only process structs, but got a %T", sc, val)
	}

	sd, err := sc.describeStruct(r.Registry, val.Type())
	if err != nil {
		return err
	}

	dw, err := vw.WriteDocument()
	if err != nil {
		return err
	}
	var rv reflect.Value
	for _, desc := range sd.fl {
		if desc.inline == nil {
			rv = val.Field(desc.idx)
		} else {
			rv = val.FieldByIndex(desc.inline)
		}

		codec := desc.codec

		iszero := sc.isZero
		if iz, ok := codec.(CodecZeroer); ok {
			iszero = iz.IsZero
		}

		if desc.omitEmpty && iszero(rv.Interface()) {
			continue
		}

		vw2, err := dw.WriteDocumentElement(desc.name)
		if err != nil {
			return err
		}

		ectx := EncodeContext{Registry: r.Registry, MinSize: desc.minSize}
		err = codec.EncodeValue(ectx, vw2, rv.Interface())
		if err != nil {
			return err
		}
	}

	if sd.inlineMap >= 0 {
		rv := val.Field(sd.inlineMap)
		collisionFn := func(key string) bool {
			_, exists := sd.fm[key]
			return exists
		}

		return defaultMapCodec.encodeValue(r, dw, rv, collisionFn)
	}

	return dw.WriteDocumentEnd()
}

// DecodeValue implements the Codec interface.
func (sc *StructCodec) DecodeValue(r DecodeContext, vr ValueReader, i interface{}) error {
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct || !val.CanAddr() {
		return fmt.Errorf("%T can only processes addressable structs, but got %T (addressable: %t)", sc, i, val.CanAddr())
	}

	sd, err := sc.describeStruct(r.Registry, val.Type())
	if err != nil {
		return err
	}

	var dFn decodeFn
	var inlineMap reflect.Value
	if sd.inlineMap >= 0 {
		inlineMap = val.Field(sd.inlineMap)
		if inlineMap.IsNil() {
			inlineMap.Set(reflect.MakeMap(inlineMap.Type()))
		}
		dFn, err = defaultMapCodec.decodeFn(r, inlineMap)
		if err != nil {
			return err
		}
	}

	dr, err := vr.ReadDocument()
	if err != nil {
		return err
	}

	for {
		name, vr, err := dr.ReadElement()
		if err == ErrEOD {
			break
		}
		if err != nil {
			return err
		}

		fd, exists := sd.fm[name]
		if !exists {
			if sd.inlineMap < 0 {
				// The encoding/json package requires a flag to return on error for non-existent fields.
				// This functionality seems appropriate for the struct codec.
				err = vr.Skip()
				if err != nil {
					return err
				}
				continue
			}

			key, elem, err := dFn(r, vr, name)
			if err != nil {
				return err
			}
			inlineMap.SetMapIndex(reflect.ValueOf(key), elem)
			continue
		}

		var field reflect.Value
		if fd.inline == nil {
			field = val.Field(fd.idx)
		} else {
			field = val.FieldByIndex(fd.inline)
		}

		if !field.CanSet() { // Being settable is a super set of being addressable.
			return fmt.Errorf("cannot decode element '%s' into field %v; it is not settable", name, field)
		}
		if field.Kind() == reflect.Ptr && field.IsNil() {
			field.Set(reflect.New(field.Type()).Elem())
		}
		field = field.Addr()

		dctx := DecodeContext{Registry: r.Registry, Truncate: fd.truncate}
		if ec, ok := fd.codec.(*elementCodec); ok {
			err = ec.decodeValue(dctx, vr, name, field.Interface().(**Element))
			if err != nil {
				return err
			}
			continue
		}

		err = fd.codec.DecodeValue(dctx, vr, field.Interface())
		if err != nil {
			fmt.Println(name)
			return err
		}
	}

	return nil
}

func (sc *StructCodec) isZero(i interface{}) bool {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		if z, ok := v.Interface().(Zeroer); ok {
			return z.IsZero()
		}
		return false
	}

	return false
}

type structDescription struct {
	fm        map[string]fieldDescription
	fl        []fieldDescription
	inlineMap int
}

type fieldDescription struct {
	name      string
	idx       int
	omitEmpty bool
	minSize   bool
	truncate  bool
	inline    []int
	codec     Codec
}

func (sc *StructCodec) describeStruct(r *Registry, t reflect.Type) (*structDescription, error) {
	// We need to analyze the struct, including getting the tags, collecting
	// information about inlining, and create a map of the field name to the field.
	sc.l.RLock()
	ds, exists := sc.cache[t]
	sc.l.RUnlock()
	if exists {
		return ds, nil
	}

	numFields := t.NumField()
	sd := &structDescription{
		fm:        make(map[string]fieldDescription, numFields),
		fl:        make([]fieldDescription, 0, numFields),
		inlineMap: -1,
	}

	for i := 0; i < numFields; i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" {
			// unexported, ignore
			continue
		}

		var codec Codec
		var err error

		switch sf.Type {
		case tElement: // We handle this as a special case within the struct codec.
			codec = defaultElementCodec
		default:
			codec, err = r.Lookup(sf.Type)
			if err != nil {
				return nil, err
			}
		}

		description := fieldDescription{idx: i, codec: codec}

		stags, err := sc.parser.ParseStructTags(sf)
		if err != nil {
			return nil, err
		}
		if stags.Skip {
			continue
		}
		description.name = stags.Name
		description.omitEmpty = stags.OmitEmpty
		description.minSize = stags.MinSize
		description.truncate = stags.Truncate

		if stags.Inline {
			switch sf.Type.Kind() {
			case reflect.Map:
				if sd.inlineMap >= 0 {
					return nil, errors.New("(struct " + t.String() + ") multiple inline maps")
				}
				if sf.Type.Key() != tString {
					return nil, errors.New("(struct " + t.String() + ") inline map must have a string keys")
				}
				sd.inlineMap = description.idx
			case reflect.Struct:
				inlinesf, err := sc.describeStruct(r, sf.Type)
				if err != nil {
					return nil, err
				}
				for _, fd := range inlinesf.fl {
					if _, exists := sd.fm[fd.name]; exists {
						return nil, fmt.Errorf("(struct %s) duplicated key %s", t.String(), fd.name)
					}
					if fd.inline == nil {
						fd.inline = []int{i, fd.idx}
					} else {
						fd.inline = append([]int{i}, fd.inline...)
					}
					sd.fm[fd.name] = fd
					sd.fl = append(sd.fl, fd)
				}
			default:
				return nil, fmt.Errorf("(struct %s) inline fields must be either a struct or a map", t.String())
			}
			continue
		}

		if _, exists := sd.fm[description.name]; exists {
			return nil, fmt.Errorf("struct %s) duplicated key %s", t.String(), description.name)
		}

		sd.fm[description.name] = description
		sd.fl = append(sd.fl, description)
	}

	sc.l.Lock()
	sc.cache[t] = sd
	sc.l.Unlock()

	return sd, nil
}
