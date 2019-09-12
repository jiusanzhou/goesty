package goesty

import (
	"net/http"
)

type InKind string

const (
	InPath   InKind = "path"
	InQuery  InKind = "query"
	InHeader InKind = "header"

	// it's not easy to implement
	InBody InKind = "body"
)

type Parameter struct {
	In       InKind
	Required bool
	DataType string
	DataFormat string
}

// ParameterExtrator take paramter
type ParameterExtrator func(r *Request, key string) (string, bool)
