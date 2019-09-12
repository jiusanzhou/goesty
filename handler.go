package goesty

import (
	"net/http"
	"reflect"
)

// newHandler init a new item t present all thing
// for implement http.Handler
//
// func(ctx context.Context, name string) (User, error) {
// 		
// }
func newHandler(v interface{}, vkr HandlerInvoker) (*Handler, error) {
	// 1. parse type of v
	//	1. must should be function (TODO: or struct)
	// 	2. take params and returns out
	// 2. init request init function
	// 3. init response init function

	ins, ous, err := vkr.Invoker(v)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Handler present wraped request and response
type Handler struct {
	Method string

	m reflect.Method // handler function from v
}

// ServeHTTP implement from net/http Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. init params from request
	// 	1. check if some params is required
	// 	2. turn type from string to correct param type
	// 2. init default value from returns's type
	// 3. call func with those params and returns
	// 4. dumps returns to response writer
}
