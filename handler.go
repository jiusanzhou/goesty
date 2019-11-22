package goesty

import (
	"fmt"
	"net/http"
	"reflect"
)

// newHandler init a new item t present all thing
// for implement http.Handler
//
// func(ctx context.Context, name string) (User, error) {
//
// }
func newHandler(v interface{}, r *Runtime) (*Handler, error) {
	// 1. parse type of v
	//	1. must should be function (TODO: or struct)
	// 	2. take params and returns out
	// 2. init request init function
	// 3. init response init function
	vkr := r.invoker

	ins, ous, ok := vkr.Invoke(v)
	if !ok {
		return nil, ErrUnknownSignature
	}

	// TODO:

	return &Handler{
		r:   r,
		m:   reflect.ValueOf(v),
		ins: ins,
		ous: ous,
	}, nil
}

// Handler present wraped request and response
type Handler struct {
	r *Runtime

	// handler function from v
	m reflect.Value
	// function inputs
	ins Values
	// function outs
	ous Values

	// this handler supported http methods
	methods []string
	// params we need to handler
	params []Parameter

	// init data to handler
	_paramIndex       int
	_paramKey         string
	_paramIndexOpened bool
	_paramKeyOpened   bool
	_paramCalled      int

	// TODO: copy options from runtime
}

// Param to set in position
func (h *Handler) Param(index int) *Handler {
	// check if the runtime has just called but not set
	if h._paramIndexOpened {
		panic("Param method need be closed by InBody/InPath... method")
	}
	// temp store
	h._paramIndex = index
	h._paramIndexOpened = true
	return h
}

// Key to set in struct
func (h *Handler) Key(key string) *Handler {
	// check if the runtime has just called but not set
	if h._paramKeyOpened {
		panic("Key method need be closed by InBody/InPath... method")
	}
	// temp store
	h._paramKey = key
	h._paramKeyOpened = true

	return h
}

// InPath take value from path
// TODO: maybe we need to use mux or something else
func (h *Handler) InPath(name string) *Handler {
	h.SetParam(PathExtrator(name), name)
	return h
}

// InQuery take value from query
func (h *Handler) InQuery(name string) *Handler {
	// set value from request's query
	h.SetParam(QueryExtrator(name), name)
	return h
}

// InHeader take value from header
func (h *Handler) InHeader(name string) *Handler {
	h.SetParam(HeaderExtrator(name), name)
	return h
}

// InCookie take value from cookie
func (h *Handler) InCookie(name string) *Handler {
	h.SetParam(CookieExtrator(name), name)
	return h
}

// InForm take value from form
func (h *Handler) InForm(name string) *Handler {
	h.SetParam(FormExtrator(name), name)
	return h
}

// SetParam set param
func (h *Handler) SetParam(extrator ParameterExtrator, key string) {
	// if param has opened
	var idx int
	if h._paramIndexOpened {
		idx = h._paramIndex
	} else {
		// try to guess the position
		idx = h._paramCalled
		// increase
		h._paramCalled++
	}

	// we need to check if the index overflow
	if idx >= len(h.ins) {
		panic(fmt.Sprintf("Param position(%d) overflow", idx))
	}

	// set inputs' value from function
	// take the value out check the value's type
	var valwrap = h.ins[idx]
	if valwrap == nil {
		// never bee here
		panic("can't get input value")
	}

	// some time we need to replace the  value from function
	// but at something we can't just replace simple, like if the value if a struct
	// and fields' value from request, not the whole body
	// we need to take the correct to value to effect.

	// must take it out

	// just replace the value from function
	// TODO: split to different functions
	valwrap.addValueFromRequest(func(v *Value, r *Request) {
		// take string out from extrate parameter
		var valstr, _ = extrator(r, key)
		v.setString(valstr) // TODO: handle with error error
	})

	// TODO: input is complex

	// set value to

	// reset all flags
	h._paramIndexOpened = false
	h._paramKeyOpened = false
}

// ServeHTTP implement from net/http Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. init params from request
	// 	1. check if some params is required
	// 	2. turn type from string to correct param type
	// 2. init default value from returns's type
	// 3. call func with those params and returns
	// 4. dumps returns to response writer

	// build my request
	mr := &Request{
		r:       h.r,
		Request: r,
		vars:    h.r.varsFunc(r),
		// TODO:
	}

	ins := h.ins.FromRequest(mr) // []reflect.Value
	rets := h.m.Call(ins)

	mw := &Response{
		r:              h.r,
		ResponseWriter: w,
	}

	// set rets to outs
	// TODO: finish those code
	h.ous.IntoResponse(mw, rets)
	mw.Flush()
}

// HandlerOption present function to change value in handler
type HandlerOption func(h *Handler)
