package goesty

import (
	"net/http"
)

// Response is a wrapper on the actual http ResponseWriter
// It provides several convenience methods to prepare and write response content.
type Response struct {
	http.ResponseWriter
	requestAccept string        // mime-type what the Http Request says it wants to receive
	routeProduces []string      // mime-types what the Route says it can produce
	statusCode    int           // HTTP status code that has been written explicitly (if zero then net/http has written 200)
	contentLength int           // number of bytes written for the response body
	prettyPrint   bool          // controls the indentation feature of XML and JSON serialization. It is initialized using var PrettyPrintResponses.
	err           error         // err property is kept when WriteError is called
	hijacker      http.Hijacker // if underlying ResponseWriter supports it
}
