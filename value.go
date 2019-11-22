package goesty

import (
	"context"
	"fmt"
	"reflect"
)

var (
	errType = reflect.TypeOf((*error)(nil)).Elem()
)

// UpdateValueFromRequest is a function to get value
type UpdateValueFromRequest func(v *Value, r *Request)

// FlushValueIntoResponse is a function to set value
type FlushValueIntoResponse func(v *Value, r *Response)

// Value contains reflect.*
type Value struct {
	name string

	typ reflect.Type
	val reflect.Value

	values Values

	valueUpdaters []UpdateValueFromRequest
	valueFlushers []FlushValueIntoResponse

	isContext bool
	isError bool
}

// return a new value from request, like copy it
func (v Value) newFromRequest(r *Request) *Value {
	mv := &Value{
		name: v.name,
		typ: v.typ,
		values: v.values, // TODO:
		valueUpdaters: v.valueUpdaters,
	}

	// if is context returns directlly
	if v.isContext {
		mv.val = reflect.ValueOf(r.Context())
		return mv
	}

	//take interface value V
	val := reflect.New(v.typ)
	vv := val.Interface()
	pv := interface{}(nil)
	if val.CanAddr() {
		pv = val.Addr().Interface()
	}
	//convert V or &V into a setter:
	for _, t := range []interface{}{vv, pv} {
		if s, ok := t.(Setter); ok {
			vv = s
		}
		// TODO: implement 3 types
		// if tm, ok := t.(encoding.TextUnmarshaler); ok {
		// 	v = &textValue{tm}
		// } else if bm, ok := t.(encoding.BinaryUnmarshaler); ok {
		// 	v = &binaryValue{bm}
		// } else if d, ok := t.(*time.Duration); ok {
		// 	v = newDurationValue(d)
		// } else if s, ok := t.(Setter); ok {
		// 	v = s
		// }
	}
	//implements setter (flag.Value)?
	if s, ok := vv.(Setter); ok {
		//NOTE: replacing val removes our ability to set
		//the value, resolved by flag.Value handling all Set calls.
		val = reflect.ValueOf(s)
	}
	//val must be concrete at this point
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// store the val
	mv.val = val

	// call data picker from request
	// TODO: struct value 
	for _, up := range v.valueUpdaters {
		up(mv, r)
	}

	return mv
}

func (v *Value) addValueFromRequest(up UpdateValueFromRequest) {
	v.valueUpdaters = append(v.valueUpdaters, up)
}

// intoResponse take the real value and set it to response at correct position
func (v Value) intoResponse(r *Response, val reflect.Value) {
	// TODO: first of all we need to check the type of value is correct
	mv := &Value{
		name: v.name,
		typ: v.typ,
		val: val,
		values: v.values, // TODO:
	}

	// special for error
	if v.isError {
		// special
		inf := val.Interface()
		if inf != nil {
			r.Error = val.Interface().(error)
		}
		return
	}

	// put me to response

	// call data flusher to response
	for _, fl := range v.valueFlushers {
		fl(mv, r)
	}
}

func (v *Value) addValueIntoResponse(fl FlushValueIntoResponse) {
	v.valueFlushers = append(v.valueFlushers, fl)
}

func newValue(typ reflect.Type) *Value {
	mv := &Value{
		typ: typ,
		// TODO:
		isContext: reflect.TypeOf((*context.Context)(nil)).Elem() == typ,
		isError: reflect.TypeOf((*error)(nil)).Elem() == typ,
	}

	// add defautl request extrator
	mv.addValueFromRequest(defaultValueFromRequest)

	// add default reponse dumps
	mv.addValueIntoResponse(defaultValueIntoResponse)

	return mv
}

// Values value ordered
type Values []*Value

// FromRequest returns all reflect.Value from request
func (vs Values) FromRequest(r *Request) []reflect.Value {
	var rvs = []reflect.Value{}
	for _, v := range vs {
		rvs = append(rvs, v.newFromRequest(r).val)
	}
	return rvs
}

// IntoResponse dumps value to reponse
func (vs Values) IntoResponse(r *Response, vals []reflect.Value) {
	// always, len(vs) == len(vals)
	// TODO: no need to check every time
	for i, v := range vs {
		v.intoResponse(r, vals[i])
	}
}

// Setter is any type which can be set from a string.
// This includes flag.Value.
// Hack use this to set value from string
type Setter interface {
	Set(string) error
}

// setString
func (v *Value) setString(valstr string) error {
	var elem reflect.Value
	if v.val.CanAddr() {
		elem = v.val.Addr() //pointer to concrete type
	} else {
		elem = v.val //possibly interface type
	}

	mv := elem.Interface()

	//convert string into value
	if fv, ok := mv.(Setter); ok {
		fv.Set(valstr)
	} else if elem.Kind() == reflect.Ptr {
		//magic set with scanf
		fmt.Sscanf(valstr, "%v", mv)
	} else {
		return fmt.Errorf("could not be set")
	}

	return nil
}

// setValue
func (v *Value) setValue(val reflect.Value) error {
	var elem reflect.Value
	if v.val.CanAddr() {
		elem = v.val.Addr() //pointer to concrete type
	} else {
		elem = v.val //possibly interface type
	}

	elem.Set(val)
	return nil
}