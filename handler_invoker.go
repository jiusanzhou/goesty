package goesty

import (
	"context"
	"errors"
	"reflect"
)

// This is the comment tag that carries parameters for open API generation.
const (
	tagName = "api"
)

var (
	// ErrUnknownSignature ...
	ErrUnknownSignature = errors.New("unknown function signature")
	// ErrUnsupportedParamType ...
	ErrUnsupportedParamType = errors.New("unsupported param type")
)

type errUnknownSignature struct {
	name string
}

func (err errUnknownSignature) Error() string {
	return "unknown fucntion signature: " + err.name
}

// HandlerInvoker create a new handler
type HandlerInvoker interface {
	Invoke(v interface{}) (Values, Values, bool)
}

// HandlerInvokerFunc is an adpater to allow original function as invoker
type HandlerInvokerFunc func(v interface{}) (Values, Values, bool)

// Invoke calls f(w, r).
func (f HandlerInvokerFunc) Invoke(v interface{}) (Values, Values, bool) {
	return f(v)
}

// HandlerInvokers slice of vkers
type HandlerInvokers []HandlerInvoker

// Invoke ...
func (ivks HandlerInvokers) Invoke(v interface{}) (Values, Values, bool) {
	for _, iv := range ivks {
		ins, ous, ok := iv.Invoke(v)
		if ok {
			return ins, ous, ok
		}
	}
	return nil, nil, false // ErrUnknownSignature
}

// --------------------- Invoker ----------------- //
// TODO: auto from struct take method out with name's

// BaseHandlerInvoker v: (args ...) (obj, error) | (obj) | (error)
func BaseHandlerInvoker(v interface{}) (Values, Values, bool) {

	handlerTyp := reflect.TypeOf(v)
	// handlerVal := reflect.ValueOf(v)

	// must be a func
	if handlerTyp.Kind() != reflect.Func {
		return nil, nil, false
	}

	// get all args
	largs := handlerTyp.NumIn()
	lrets := handlerTyp.NumOut()

	var ins = Values{}
	var ous = Values{}

	for i := 0; i < largs; i++ {
		var ou, _ = buildValueCreator(handlerTyp.In(i), 0)
		ins = append(ins, ou)
	}

	for i := 0; i < lrets; i++ {
		var ou, _ = buildValueCreator(handlerTyp.Out(i), 0)
		ous = append(ous, ou)
	}

	return ins, ous, true
}

// ContextHandlerInvoker v: (contex.Context, args ...) (obj, error) | (obj) | (error)
func ContextHandlerInvoker(v interface{}) (Values, Values, bool) {

	return nil, nil, false
}

/**
ptr 取完地址再处理,并返回回调函数
// TODO: 暂时不处理. 限定基础类型的参数可以完全来自path。string,int 等基础类型,直接处理返回
map[string]string 接收query/path 中未被捕捉的参数, 其余map接收body
struct 遍历field,并处理tag获取相应的key值,处理
interface{} 接收body

context.Context 为特殊类型值(r.Context()),不应进入本处理逻辑
*/

// buildValueCreator build value creator from request
// defualt set value from
func buildValueCreator(typ reflect.Type, depth int) (*Value, error) {

	// we can take correct value from typ
	if typ == reflect.TypeOf((*context.Context)(nil)).Elem() {
		// context.Context we just return. should we check the depth?
		return newValue(typ), nil
	}

	switch typ.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, // Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.String:
		// basic type

		// field comes from function params
		// so it's a little difficate to disting where we should take value from
		if depth == 0 {
			// TODO:
			return newValue(typ), nil
		}

		// parse tag from field, TODO: use object from tag in x
		return newValue(typ), nil
	case reflect.Ptr:
		// take elem type

	case reflect.Struct:
		// iter fields
	case reflect.Interface:
		// unmarshal from body if depth == 0 else just take from what it will be
	case reflect.Map:
		// if depth == 0 only supported map[string]interface{} and take data from path and query
		// else in `querys | paths | headers`
	case reflect.Array:
		// if depth == 0 unsupported
	}
	return newValue(typ), ErrUnsupportedParamType
}
