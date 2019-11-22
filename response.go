package goesty

import (
	"net/http"
)

var (
	// StatusSuccess present this request success
	StatusSuccess = "success"
	// StatusFailed repsent this request failed
	StatusFailed = "failed"
)

// Response is a wrapper on the actual http ResponseWriter
// It provides several convenience methods to prepare and write response content.
type Response struct {
	r *Runtime

	http.ResponseWriter `json:"-"`

	Error error `json:"-"`

	Code     StatusCode  `json:"code"`
	Data     interface{} `json:"data,omitempty"`
	Status   string      `json:"status"`
	ErrorStr string      `json:"error,omitempty"`

	kvs map[string]interface{}

	prettyPrint bool
	statusCode  int
}

// Flush data to response writer
func (r *Response) Flush() {
	// TODO:
	// callback for data

	// callback for error

	// called hooks
	r.r.responseHook(r)

	// callback for code
	r.Code = r.r.codeFromError(r.Error)

	// callback for statusCode
	r.statusCode = r.r.statusCodeFromCode(r.Code)

	// flush response to responseWriter
	if r.Error != nil {
		r.ErrorStr = r.Error.Error()
	}

	if r.Code != CodeOK || r.Error != nil {
		r.Status = StatusFailed
		if r.Code == CodeOK {
			// set default code
			r.Code = CodeInternal
		}
	} else {
		r.Status = StatusSuccess
	}

	// TODO: fix data
	data := r.Data.([]interface{})
	if len(data) == 1 {
		r.Data = data[0]
	}

	// fxcking things
	if r.r.onlyData {
		r.r.entitor.Write(r, r.statusCode, r.Data)
	} else {
		r.r.entitor.Write(r, r.statusCode, r)
	}
}

func defaultResponseDump(r *Response) {

}

func defaultResponseHook(r *Response) {
	// do nothing
}

func defaultValueIntoResponse(v *Value, r *Response) {
	if r.Data == nil {
		r.Data = []interface{}{}
	}
	data := r.Data.([]interface{})
	data = append(data, v.val.Interface())
	r.Data = data
}
