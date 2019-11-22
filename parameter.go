package goesty

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
	In         InKind
	Required   bool
	DataType   string
	DataFormat string
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
