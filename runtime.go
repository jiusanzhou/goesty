package goesty

import (
	"net/http"
)

// MiddlewareFunc is a function which receives an http.Handler and returns another http.Handler.
// Typically, the returned handler is a closure which does something with the http.ResponseWriter and http.Request passed
// to it, and then calls the handler passed as parameter to the MiddlewareFunc.
type MiddlewareFunc func(http.Handler) http.Handler

// Option is a function to re configurate runtime
type Option func(r *Runtime)

// Runtime present all function to generate handler
type Runtime struct {
	invoker HandlerInvokers // invoker create handler from v(interface{})

	// responseForward

	// middlewares
	middlewares []MiddlewareFunc

	// special middleware vars getter
	varsFunc func(r *http.Request) map[string]string

	// response hook
	responseHook func(r *Response)

	// entity readerWriter
	entitor EntityReaderWriter

	// only marshal data to body
	onlyData bool

	codeFromError func(err error) StatusCode

	statusCodeFromCode func(code StatusCode) int
}

/**
- generate a object from input(func/struct)
- create handler from object
  - request to some input fields
  - output field to response
  - error take code out and set for status code
*/

func newRuntime(opts ...Option) *Runtime {
	r := &Runtime{
		// add default invoker
		invoker: HandlerInvokers{
			HandlerInvokerFunc(BaseHandlerInvoker),
		},
		// add defualt entity
		entitor:            NewEntityAccessorJSON(MIME_JSON),
		statusCodeFromCode: defaultHTTPStatusFromCode,
		codeFromError:      defaultCodeFromError,
		responseHook:       defaultResponseHook,
	}

	for _, o := range opts {
		o(r)
	}
	return r
}

// newHandler create a new handler
func (r *Runtime) newHandler(v interface{}, opts ...Option) (*Handler, error) {
	// TODO: copy a new runtime deep
	// opts just for this handler

	// use the invoker

	// TODO: generate a handler and call middleware

	return newHandler(v, r)
}

// New create a new handler or handler func
// return error if input wrong type
func (r *Runtime) New(v interface{}, opts ...Option) (*Handler, error) {
	return r.newHandler(v, opts...)
}

// MustNew create a new handler or handler func
// panic if any error
func (r *Runtime) MustNew(v interface{}, opts ...Option) *Handler {
	h, err := r.newHandler(v, opts...)
	if err != nil {
		panic(err)
	}
	return h
}

// OptionMiddleware add middlewares
func OptionMiddleware(wares ...MiddlewareFunc) Option {
	return func(r *Runtime) {
		r.middlewares = append(r.middlewares, wares...)
	}
}

// OptionVarsFunc set vars getter
func OptionVarsFunc(fn func(r *http.Request) map[string]string) Option {
	return func(r *Runtime) {
		r.varsFunc = fn
	}
}

// OptionResponseHook set response func
func OptionResponseHook(fn func(r *Response)) Option {
	return func(rt *Runtime) {
		var oldfn = rt.responseHook
		if oldfn != nil {
			rt.responseHook = func(r *Response) {
				oldfn(r)
				fn(r)
			}
		} else {
			rt.responseHook = fn
		}
	}
}

// OptionEntity set entity read writer
func OptionEntity(entitor EntityReaderWriter) Option {
	return func(r *Runtime) {
		r.entitor = entitor
	}
}

// OptionCodeFromError turn error code
func OptionCodeFromError(fn func(err error) StatusCode) Option {
	return func(r *Runtime) {
		var oldfn = r.codeFromError
		r.codeFromError = func(err error) StatusCode {
			var v = oldfn(err)
			if v != CodeInternal {
				return v
			}
			return fn(err)
		}
	}
}

// OptionStatusCodeFromCode turn code to http status code
func OptionStatusCodeFromCode(fn func(code StatusCode) int) Option {
	return func(r *Runtime) {
		var oldfn = r.statusCodeFromCode
		r.statusCodeFromCode = func(code StatusCode) int {
			var v = oldfn(code)
			if v != http.StatusInternalServerError {
				return v
			}
			return fn(code)
		}
	}
}
