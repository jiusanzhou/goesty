package goesty

import "reflect"

type contextKey int

// InKind is a type about where paramter in
type InKind string

const (
	// InPath /user/{name}
	InPath InKind = "path"
	// InQuery /user?name={name}
	InQuery InKind = "query"
	// InHeader Authorization: {token}
	InHeader InKind = "header"
	// InBody {}
	// it's not easy to implement
	InBody InKind = "body" // json or protobuf?

	// InCookie in cookie, Cookie: session=xxxx
	InCookie InKind = "cookie"
	// InForm in post form
	InForm InKind = "form"

	varsKey contextKey = iota
	routeKey
)

// Parameter present state of paramter
type Parameter struct {
	h *Handler

	name     string
	in       InKind
	required bool
	def      interface{}
	defval   reflect.Value

	dataType   string
	dataFormat string

	ext ParameterExtrator
}

// Name ...
func (p *Parameter) Name(v string) *Parameter {
	p.name = v
	return p
}

// Default ...
func (p *Parameter) Default(v interface{}) *Parameter {
	p.def = v
	p.defval = reflect.ValueOf(v)
	return p
}

// Required ...
func (p *Parameter) Required(v bool) *Parameter {
	p.required = v
	return p
}

// In ...
func (p *Parameter) In(v InKind) *Parameter {
	p.in = v
	return p
}

// At ...
func (p *Parameter) At(index int) *Handler {
	// set param to handler
	p.initExt()

	if p.ext == nil {
		// !!! ingore
		return p.h
	}

	p.h.SetParam(index, p)

	return p.h
}

func (p *Parameter) initExt() {

	if p.ext != nil {
		return
	}

	switch p.in {
	case InPath:
		p.ext = PathExtrator(p.name)
	case InQuery:
		p.ext = QueryExtrator(p.name)
	case InHeader:
		p.ext = HeaderExtrator(p.name)
	case InForm:
		p.ext = FormExtrator(p.name)
	case InCookie:
		p.ext = CookieExtrator(p.name)
	case InBody:

	}
}

// NewParam ...
func NewParam(h *Handler) *Parameter {
	return &Parameter{
		h:  h,
		in: InQuery,
	}
}

// ParameterExtrator take paramter
type ParameterExtrator func(r *Request, key string) (string, bool)

// QueryExtrator extrate parameter from query
func QueryExtrator(name string) ParameterExtrator {
	return func(r *Request, key string) (string, bool) {
		return r.FormValue(name), true
	}
}

// PathExtrator take paramter
func PathExtrator(name string) ParameterExtrator {
	return func(r *Request, key string) (string, bool) {
		if r.vars == nil {
			return "", false
		}
		v, ok := r.vars[name]
		return v, ok
	}
}

// HeaderExtrator take paramter
func HeaderExtrator(name string) ParameterExtrator {
	return func(r *Request, key string) (string, bool) {
		return r.Header.Get(name), true
	}
}

// CookieExtrator take paramter from cookie
func CookieExtrator(name string) ParameterExtrator {
	return func(r *Request, key string) (string, bool) {
		if cookie, err := r.Cookie(name); err == nil {
			return cookie.Value, true
		}
		return "", false
	}
}

// FormExtrator take paramter from form
func FormExtrator(name string) ParameterExtrator {
	return func(r *Request, key string) (string, bool) {
		return r.FormValue(name), true
	}
}

// Vars returns the route variables for the current request, if any.
func Vars(r *Request) map[string]string {
	if rv := r.Context().Value(varsKey); rv != nil {
		return rv.(map[string]string)
	}
	return nil
}
