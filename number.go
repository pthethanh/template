package template

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

type (
	operator int
	calFunc  = func(values ...interface{}) (float64, error)
)

const (
	mul operator = iota
	add
	div
	sub
	pow
)

// NumberFuncMap return number func map.
func NumberFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"mul": cal(mul),
		"add": cal(add),
		"sum": cal(add),
		"div": cal(div),
		"sub": cal(sub),
		"pow": cal(pow),
	}
}

func cal(op operator) calFunc {
	return func(values ...interface{}) (float64, error) {
		r := float64(1.0)
		for i, v := range values {
			rv := reflect.ValueOf(v)
			ik, err := basicKind(rv)
			if err != nil {
				return 0, err
			}
			iv := 0.0
			switch ik {
			case uintKind, intKind:
				iv = float64(rv.Int())
			case floatKind:
				iv = rv.Float()
			case stringKind:
				// string will be converted to float
				iv, err = strconv.ParseFloat(v.(string), 64)
				if err != nil {
					return 0, err
				}
			default:
				return 0, fmt.Errorf("value must be number kind, kind: %v", ik)
			}
			if i == 0 {
				r = iv
				continue
			}
			switch op {
			case mul:
				r *= iv
			case add:
				r += iv
			case div:
				r /= iv
			case sub:
				r -= iv
			case pow:
				r = math.Pow(r, iv)
			}
		}
		return r, nil
	}
}
