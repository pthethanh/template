package template

import (
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

// GeneralFuncMap return general func map.
func GeneralFuncMap() map[string]interface{} {
	return map[string]interface{}{
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
		rs.WriteString(fmt.Sprintf("%v", printableValue(reflect.ValueOf(v))))
	}
	return rs.String()
}

// Join join the string representation of the values together.
// String will be joined as whole.
// Map, slice, array will be joined using its value, one by one.
func Join(sep string, values ...interface{}) string {
	rs := make([]string, 0)
	for _, val := range values {
		v, isNil := indirect(reflect.ValueOf(val))
		if isNil {
			return ""
		}
		switch v.Kind() {
		case reflect.String:
			rs = append(rs, v.String())
		case reflect.Array, reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				rs = append(rs, fmt.Sprintf("%v", printableValue(v.Index(i))))
			}
		case reflect.Map:
			r := v.MapRange()
			for r.Next() {
				rs = append(rs, fmt.Sprintf("%v", printableValue(r.Value())))
			}
		default:
			rs = append(rs, fmt.Sprintf("%v", printableValue(v)))
		}
	}
	return strings.Join(rs, sep)
}

// Contains check whether all the values exist in the collection.
// The collection must be a slice, array, string or a map.
func Contains(collection reflect.Value, values ...reflect.Value) bool {
	for _, val := range values {
		if ok := contains(collection, val); !ok {
			return false
		}
	}
	return true
}

// ContainsAny check whether one of the value exist in the collection.
// The collection must be a slice, array, string or a map.
func ContainsAny(collection reflect.Value, values ...reflect.Value) bool {
	for _, val := range values {
		if ok := contains(collection, val); ok {
			return true
		}
	}
	return false
}

func contains(collection reflect.Value, val reflect.Value) bool {
	v, isNil := indirect(collection)
	if isNil {
		return false
	}
	val, isNil = indirect(val)
	switch v.Kind() {
	case reflect.String:
		// accept all kinds of val.
		return strings.Contains(v.String(), fmt.Sprintf("%v", val))
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			iv, vIsNil := indirect(v.Index(i))
			// accept compare nil, nil
			if isNil && vIsNil || (!val.IsValid() && !iv.IsValid()) {
				return true
			}
			if ok, _ := eq(val, iv); ok {
				return true
			}
		}
	case reflect.Map:
		r := v.MapRange()
		for r.Next() {
			iv, vIsNil := indirect(r.Value())
			// accept compare nil, nil
			if isNil && vIsNil || (!val.IsValid() && !iv.IsValid()) {
				return true
			}
			if ok, _ := eq(iv, val); ok {
				return true
			}
		}
	default:
		return false
	}
	return false
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
