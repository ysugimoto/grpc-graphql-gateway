package runtime

import (
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

func derefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func MarshalResponse(resp interface{}) interface{} {
	// If response is nil, nothing to do.
	if resp == nil {
		return resp
	}
	v := derefValue(reflect.ValueOf(resp))
	switch v.Kind() {
	case reflect.Struct:
		return marshalStruct(v)
	case reflect.Slice:
		return marshalSlice(v)
	default:
		return primitive(v)
	}
}

// Marshal reflect value to []interface{} with lower camel case field
func marshalSlice(v reflect.Value) []interface{} {
	size := v.Len()
	ret := make([]interface{}, size)

	for i := 0; i < size; i++ {
		vv := derefValue(v.Index(i))
		switch vv.Kind() {
		case reflect.Struct:
			ret[i] = marshalStruct(vv)
		case reflect.Slice:
			ret[i] = marshalSlice(vv)
		default:
			ret[i] = primitive(vv)
		}
	}
	return ret
}

// Marshal reflect value to map[string]interface{} with lower camel case field
func marshalStruct(v reflect.Value) map[string]interface{} {
	ret := make(map[string]interface{})
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		// If "json" tag is not set in struct field, it's not related to response field
		// So we can skip marshaling
		tag := t.Field(i).Tag.Get("json")
		if tag == "" {
			continue
		}

		name := strcase.ToLowerCamel(strings.TrimSuffix(tag, ",omitempty"))
		vv := derefValue(v.Field(i))

		switch vv.Kind() {
		case reflect.Struct:
			ret[name] = marshalStruct(vv)
		case reflect.Slice:
			ret[name] = marshalSlice(vv)
		default:
			ret[name] = primitive(vv)
		}
	}
	return ret
}

// Get primitive type value
// Protobuf in Go only use a few scalar types.
// See: https://developers.google.com/protocol-buffers/docs/proto3#scalar
func primitive(v reflect.Value) interface{} {
	// Guard by cheking IsValid due to prevent panic in runtime
	if !v.IsValid() {
		return nil
	}

	switch v.Type().Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return v.Bool()
	case reflect.Int:
		return int(v.Int())
	case reflect.Int32:
		return int32(v.Int())
	case reflect.Int64:
		return int64(v.Int())
	case reflect.Uint:
		return uint(v.Uint())
	case reflect.Uint32:
		return uint32(v.Uint())
	case reflect.Uint64:
		return uint64(v.Uint())
	case reflect.Float32:
		return float32(v.Float())
	case reflect.Float64:
		return float64(v.Float())
	default:
		return nil
	}
}
