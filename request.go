package goesty

import (
	"context"
	"net/http"
	"reflect"
)

// Request is a wrapper for a http Request that provides convenience methods
type Request struct {
	r *Runtime
	*http.Request
	vars              map[string]string
	attributes        map[string]interface{} // for storing request-scoped values
	selectedRoutePath string                 // root path + route path that matched the request, e.g. /meetings/{id}/attendees
}

// withVars set vars for Request
func withVars(r *http.Request, vars map[string]string) *http.Request {
	ctx := context.WithValue(r.Context(), varsKey, vars)
	return r.WithContext(ctx)
}

func contextFromRequest(v *Value, r *Request) {
	// to slice
	v.val = reflect.ValueOf(r.Context())
}

func defaultValueFromRequest(v *Value, r *Request) {
	// from query
	v.setString(r.FormValue(v.name))
}
