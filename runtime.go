package goesty

import (
	"go.zoe.im/x/httputil"
)

// Option is a function to re configurate runtime
type Option func(r *Runtime)

// Runtime present all function to generate handler
type Runtime struct {
	invoker HandlerInvoker // invoker create handler from v(interface{})

	// responseForward

	// convert error to status code
	statusCodeFromError func(err error) httputil.StatusCode

	// extend for x.httputil convert status code to http code
	httpStatusFromCode func(code httputil.StatusCode) int
}

/**
- generate a object from input(func/struct)
- create handler from object
  - request to some input fields
  - output field to response
  - error take code out and set for status code
*/

func (r *Runtime) newRuntime(opts ...Option) *Runtime {
	r = &Runtime{}
	for _, o := range opts {
		o(r)
	}
	return r
}

// newHandler create a new handler
func (r *Runtime) newHandler(v interface{}, opts ...Option) (*Handler, error) {
	// TODO: copy a new runtime deep
	// use the invoker
	return newHandler(v, r.invoker)
}

// New create a new handler or handler func
// return error if input wrong type
func (r *Runtime) New(v interface{}, opts ...Option) (http.Handler, error) {
	return r.newHandler(v, opts...)
}

// MustNew create a new handler or handler func
// panic if any error
func (r *Runtime) MustNew(v interface{}, opts ...Option) http.Handler {
	h, err := r.newHandler(v, opts...)
	if err != nil {
		panic(err)
	}
	return h
}
