package gocancel

import (
	"fmt"
	"net/http"
)

// Error is the response returned when a call is unsuccessful.
type Error struct {
	Response *http.Response // HTTP response that caused this error

	Code    string `json:"code"`    // error code
	Message string `json:"message"` // error message
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v %v: %v (%d) %v",
		e.Response.Request.Method, e.Response.Request.URL,
		e.Code, e.Response.StatusCode, e.Message)
}

// Is returns whether the provided error equals this error.
func (r *Error) Is(target error) bool {
	v, ok := target.(*Error)
	if !ok {
		return false
	}

	if r.Message != v.Message || r.Code != v.Code || !compareHttpResponse(r.Response, v.Response) {
		return false
	}

	return true
}

// errorRoot deserializes the outer JSON object returned in an error response
// from the API.
type errorRoot struct {
	Error *Error `json:"error,omitempty"`
}
