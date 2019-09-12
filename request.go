package goesty

import (
	"net/http"
)

// Request is a wrapper for a http Request that provides convenience methods
type Request struct {
	Request           *http.Request
	pathParameters    map[string]string
	attributes        map[string]interface{} // for storing request-scoped values
	selectedRoutePath string                 // root path + route path that matched the request, e.g. /meetings/{id}/attendees
}