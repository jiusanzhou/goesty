package goesty

import (
	"errors"
)

// This is the comment tag that carries parameters for open API generation.
const (
	tagName = "api"
)

var (
	ErrUnknownSignature     = errors.New("unknown fucntion signature")
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

// HandlerInvokers slice of vkers
type HandlerInvokers []HandlerInvoker

func (ivks HandlerInvokers) Invoke(v interface{}) (Values, Values, error) {
	for _, iv := ivks {
		ins, ous, ok := iv.Invoke(v)
		if ok {
			return ins, ous, nil
		}
	}
	return nil, nil, ErrUnknownSignature
}

// --------------------- Invoker ----------------- //

/**
ptr 取完地址再处理,并返回回调函数
// TODO: 暂时不处理. 限定基础类型的参数可以完全来自path。string,int 等基础类型,直接处理返回
map[string]string 接收query/path 中未被捕捉的参数, 其余map接收body
struct 遍历field,并处理tag获取相应的key值,处理
interface{} 接收body

context.Context 为特殊类型值(r.Context()),不应进入本处理逻辑
*/

// buildValueCreator build value creator from request
func buildValueCreator(typ reflect.Type, depth int) (Value, error) {

	// we can take correct value from typ
	if typ == reflect.TypeOf((*context.Context)(nil)).Elem() {
		// context.Context we just return. should we check the depth?
		return newValue(typ, func(w http.ResponseWriter, r *http.Request) interface{} {
			// TODO: we should set request and response writer
			return r.Context()
		})
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
			return newValue(typ, func(w http.ResponseWriter, r *http.Request) interface{} {
				// 
			})
		}

		// parse tag from field, TODO: use object from tag in x
		return newValue(typ, func(w http.ResponseWriter, r *http.Request) interface{} {

		})
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
	default:
		return nil, ErrUnsupportedParamType
	}
}