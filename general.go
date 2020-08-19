package template

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strings"
)

// GeneralFuncMap return general func map.
func GeneralFuncMap() template.FuncMap {
	return template.FuncMap{
		"is_true":  IsTrue,
		"defined":  IsTrue,
		"empty":    IsEmpty,
		"default":  Default,
		"ternary":  Ternary,
		"coalesce": Coalesce,
		"env":      os.Getenv,
		"contains": Contains,
	}
}

// Contains check whether the idx is in item.
func Contains(val reflect.Value, item reflect.Value) (bool, error) {
	v := indirectInterface(item)
	if !v.IsValid() {
		return false, errors.New("invalid value")
	}
	rVal := indirectInterface(val)
	if !rVal.IsValid() {
		return false, errors.New("invalid value")
	}
	switch v.Kind() {
	case reflect.String:
		// accept all kinds of val.
		return strings.Contains(v.String(), fmt.Sprintf("%v", rVal)), nil
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			ok, err := eq(rVal, indirectInterface(v.Index(i)))
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
	case reflect.Map:
		r := v.MapRange()
		for r.Next() {
			ok, err := eq(r.Value(), val)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
	case reflect.Invalid:
		return false, nil
	default:
		return false, nil
	}
	return false, nil
}

// Ternary returns the first value if the last value has meaningful value/IsTrue, otherwise returns the second value.
func Ternary(vt interface{}, vf interface{}, v interface{}) interface{} {
	if IsTrue(v) {
		return vt
	}

	return vf
}

// Coalesce return first meaningful value (IsTrue).
func Coalesce(v ...reflect.Value) reflect.Value {
	for _, val := range v {
		if IsTrue(val) {
			return val
		}
	}
	return zero
}

// Default return default value if the given value is not meaningful (not IsTrue).
func Default(df interface{}, v interface{}) interface{} {
	if IsEmpty(v) {
		return df
	}
	return v
}

// IsTrue reports whether the value is 'true', in the sense of not the zero of its type,
// and whether the value has a meaningful truth value. This is the definition of
// truth used by if and other such actions.
func IsTrue(v interface{}) bool {
	if ok, _ := template.IsTrue(v); ok {
		return ok
	}
	return false
}

// IsEmpty report whether the value not holding meaningful value.
// Opposite with IsTrue.
func IsEmpty(v interface{}) bool {
	return !IsTrue(v)
}
