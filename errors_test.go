package gocancel

import (
	"errors"
	"net/http"
	"testing"
)

func TestError_Is(t *testing.T) {
	err := &Error{
		Response: &http.Response{},
		Code:     "c",
		Message:  "m",
	}
	testcases := map[string]struct {
		wantSame   bool
		otherError error
	}{
		"errors are same": {
			wantSame: true,
			otherError: &Error{
				Response: &http.Response{},
				Code:     "c",
				Message:  "m",
			},
		},
		"errors have different values - Message": {
			wantSame: false,
			otherError: &Error{
				Response: &http.Response{},
				Code:     "c",
				Message:  "m1",
			},
		},
		"errors have different values - Response is nil": {
			wantSame: false,
			otherError: &Error{
				Code:    "c",
				Message: "m",
			},
		},
		"errors have different values - Code": {
			wantSame: false,
			otherError: &Error{
				Response: &http.Response{},
				Code:     "c1",
				Message:  "m",
			},
		},
		"errors have different types": {
			wantSame:   false,
			otherError: errors.New("Github"),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			if tc.wantSame != err.Is(tc.otherError) {
				t.Errorf("Error = %#v, want %#v", err, tc.otherError)
			}
		})
	}
}
