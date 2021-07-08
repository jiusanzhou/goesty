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
	//	1. must should be function (TODO: or struct for model CRUD?)
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
		r:      r,
		m:      reflect.ValueOf(v),
		ins:    ins,
		ous:    ous,
		params: make([]*Parameter, len(ins)),
	}, nil

	// auto add params
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
	params []*Parameter

	// TODO: copy options from runtime

	// parse struct be here
}

// Param ... reutrn Param to be set
func (h *Handler) Param(name string) *Parameter {
	return NewParam(h).Name(name)
}

// InPath take value from path
// TODO: maybe we need to use mux or something else
func (h *Handler) InPath(name string) *Parameter {
	return h.Param(name).In(InPath)
}

// InQuery take value from query
func (h *Handler) InQuery(name string) *Parameter {
	return h.Param(name).In(InQuery)
}

// InHeader take value from header
func (h *Handler) InHeader(name string) *Parameter {
	return h.Param(name).In(InHeader)
}

// InCookie take value from cookie
func (h *Handler) InCookie(name string) *Parameter {
	return h.Param(name).In(InCookie)
}

// InForm take value from form
func (h *Handler) InForm(name string) *Parameter {
	return h.Param(name).In(InForm)
}

// SetParam set param
func (h *Handler) SetParam(idx int, p *Parameter) {
	// if param has opened

	// we need to check if the index overflow
	if idx >= len(h.ins) {
		panic(fmt.Sprintf("Param position(%d) overflow", idx))
	}

	p.h.params[idx] = p

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
		var valstr, _ = p.ext(r, p.name)
		// default value
		if valstr == "" {
			v.setValue(p.defval)
		} else {
			v.setString(valstr) // TODO: handle with error error
		}
	})

	// TODO: input is complex

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
