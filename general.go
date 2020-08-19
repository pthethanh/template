package template

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

// GeneralFuncMap return general func map.
func GeneralFuncMap() template.FuncMap {
	return template.FuncMap{
		"is_true":      IsTrue,
		"is_empty":     IsEmpty,
		"default":      Default,
		"yesno":        YesNo,
		"coalesce":     Coalesce,
		"env":          os.Getenv,
		"contains":     Contains,
		"contains_any": ContainsAny,
		"file_size":    FileSizeFormat,
		"uuid":         UUID,
		"repeat":       Repeat,
		"join":         Join,
	}
}

// Repeat repeats the string representation of value n times.
func Repeat(n int, v interface{}) string {
	rs := &strings.Builder{}
	for i := 0; i < n; i++ {
		rs.WriteString(fmt.Sprintf("%v", v))
	}
	return rs.String()
}

// Join join the string representation of the values together.
func Join(sep string, values ...interface{}) string {
	rs := &strings.Builder{}
	for _, v := range values {
		rs.WriteString(fmt.Sprintf("%v%s", v, sep))
	}
	return rs.String()
}

// Contains check whether all the values exist in the collection.
// The collection must be a slice, array, string or a map.
func Contains(collection reflect.Value, values ...reflect.Value) bool {
	for _, val := range values {
		if ok, err := contains(collection, val); !ok || err != nil {
			return false
		}
	}
	return true
}

// ContainsAny check whether one of the value exist in the collection.
// The collection must be a slice, array, string or a map.
func ContainsAny(collection reflect.Value, values ...reflect.Value) bool {
	for _, val := range values {
		if ok, err := contains(collection, val); ok && err == nil {
			return true
		}
	}
	return false
}

func contains(collection reflect.Value, val reflect.Value) (bool, error) {
	v := indirectInterface(collection)
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
	default:
		return false, nil
	}
	return false, nil
}

// YesNo returns the first value if the last value has meaningful value/IsTrue, otherwise returns the second value.
func YesNo(y interface{}, n interface{}, v interface{}) interface{} {
	if IsTrue(v) {
		return y
	}
	return n
}

// Coalesce return first meaningful value (IsTrue).
func Coalesce(v ...interface{}) interface{} {
	for _, val := range v {
		if IsTrue(val) {
			return val
		}
	}
	return nil
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
	if truth, ok := template.IsTrue(v); truth && ok {
		return ok
	}
	return false
}

// IsEmpty report whether the value not holding meaningful value.
// Opposite with IsTrue.
func IsEmpty(v interface{}) bool {
	return !IsTrue(v)
}

// UUID return a UUID.
func UUID() string {
	return uuid.New().String()
}

// FileSizeFormat return human readable string of file size.
func FileSizeFormat(value interface{}) string {
	var size float64

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		size = float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		size = float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		size = v.Float()
	default:
		return ""
	}

	var KB float64 = 1 << 10
	var MB float64 = 1 << 20
	var GB float64 = 1 << 30
	var TB float64 = 1 << 40
	var PB float64 = 1 << 50

	filesizeFormat := func(filesize float64, suffix string) string {
		return strings.Replace(fmt.Sprintf("%.1f %s", filesize, suffix), ".0", "", -1)
	}

	var result string
	if size < KB {
		result = filesizeFormat(size, "bytes")
	} else if size < MB {
		result = filesizeFormat(size/KB, "KB")
	} else if size < GB {
		result = filesizeFormat(size/MB, "MB")
	} else if size < TB {
		result = filesizeFormat(size/GB, "GB")
	} else if size < PB {
		result = filesizeFormat(size/TB, "TB")
	} else {
		result = filesizeFormat(size/PB, "PB")
	}

	return result
}
