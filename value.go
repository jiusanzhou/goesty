package goesty

import (
	"reflect"
)

// Value contains reflect.*
type Value struct {
	typ reflect.Type
	val reflect.Value

	ValueFrom func(r *http.Request) interface{}
}

func newValue(typ reflect.Type, valueFrom) Value {
	return Value{
		typ: typ,
		ValueFrom: valueFrom,
	}
}

// Couple value ordered
type Values []Value

// Values returns all reflect.Value
func (vs Values) Values() []reflect.Value {
	var rvs = []reflect.Value
	for _, v := range vs {
		rvs = append(rvs, v.val)
	}
	return rvs
}