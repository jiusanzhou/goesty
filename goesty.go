package goesty

var (
	defualtRuntime *Runtime
)

func init() {
	defualtRuntime = newRuntime(
	// add all default options
	)
}

// NewRuntime create a handler factory
func NewRuntime(opts ...Option) *Runtime {
	return newRuntime(opts...)
}

// New export from default runtime
func New(v interface{}, opts ...Option) (*Handler, error) {
	return defualtRuntime.New(v, opts...)
}

// MustNew export from default runtime
func MustNew(v interface{}, opts ...Option) *Handler {
	return defualtRuntime.MustNew(v, opts...)
}
