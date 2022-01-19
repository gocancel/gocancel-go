package gocancel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nsf/jsondiff"
)

// setup sets up a test HTTP server along with a gocancel.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)

	// client is the GoCancel client being tested and is
	// configured to use test server.
	client = NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; !cmp.Equal(got, want) {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testFormValues(t *testing.T, r *http.Request, want url.Values) {
	t.Helper()

	_ = r.ParseForm()
	if got := r.Form; !cmp.Equal(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func testURLParseError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

// Test whether the marshaling of v produces JSON that corresponds
// to the want string.
func testJSONMarshal(t *testing.T, v interface{}, want string) {
	t.Helper()
	// Unmarshal the wanted JSON, to verify its correctness, and marshal it back
	// to sort the keys.
	u := reflect.New(reflect.TypeOf(v)).Interface()
	if err := json.Unmarshal([]byte(want), &u); err != nil {
		t.Errorf("Unable to unmarshal JSON for %v: %v", want, err)
	}
	w, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", u)
	}

	// Marshal the target value.
	j, err := json.Marshal(v)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", v)
	}

	diffOpts := jsondiff.DefaultConsoleOptions()
	res, diff := jsondiff.Compare([]byte(w), []byte(j), &diffOpts)

	if res != jsondiff.FullMatch {
		t.Errorf("json.Marshal(%q) is invalid, want %s", v, diff)
	}
}

// Test how bad options are handled. Method f under test should
// return an error.
func testBadOptions(t *testing.T, methodName string, f func() error) {
	t.Helper()
	if methodName == "" {
		t.Error("testBadOptions: must supply method methodName")
	}
	if err := f(); err == nil {
		t.Errorf("bad options %v err = nil, want error", methodName)
	}
}

// Test function under NewRequest failure and then s.client.Do failure.
// Method f should be a regular call that would normally succeed, but
// should return an error when NewRequest or s.client.Do fails.
func testNewRequestAndDoFailure(t *testing.T, methodName string, client *Client, f func() (*Response, error)) {
	t.Helper()
	if methodName == "" {
		t.Error("testNewRequestAndDoFailure: must supply method methodName")
	}

	client.BaseURL.Path = ""
	resp, err := f()
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' %v resp = %#v, want nil", methodName, resp)
	}
	if err == nil {
		t.Errorf("client.BaseURL.Path='' %v err = nil, want error", methodName)
	}
}

func TestClient(t *testing.T) {
	c := NewClient(nil)
	c2 := c.Client()
	if c.client == c2 {
		t.Error("Client returned same http.Client, but should be different")
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}

	if got, want := c.UserAgent, userAgent; got != want {
		t.Errorf("NewClient UserAgent is %v, want %v", got, want)
	}

	c2 := NewClient(nil)
	if c.client == c2.client {
		t.Error("NewClient returned same http.Clients, but they should differ")
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL, outURL := "/foo", defaultBaseURL+"foo"
	inBody, outBody := &Category{ID: String("foo")}, `{"id":"foo"}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%q) Body is %v, want %v", inBody, got, want)
	}

	// test that default user-agent is attached to the request
	if got, want := req.Header.Get("User-Agent"), c.UserAgent; got != want {
		t.Errorf("NewRequest() User-Agent is %v, want %v", got, want)
	}
}

func TestNewRequest_invalidJSON(t *testing.T) {
	c := NewClient(nil)

	type T struct {
		A map[interface{}]interface{}
	}
	_, err := c.NewRequest("GET", ".", &T{})

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("Expected a JSON error; got %#v.", err)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient(nil)
	_, err := c.NewRequest("GET", ":", nil)
	testURLParseError(t, err)
}

func TestNewRequest_badMethod(t *testing.T) {
	c := NewClient(nil)
	if _, err := c.NewRequest("BOGUS\nMETHOD", ".", nil); err == nil {
		t.Fatal("NewRequest returned nil; expected error")
	}
}

func TestNewRequest_emptyUserAgent(t *testing.T) {
	c := NewClient(nil)
	c.UserAgent = ""
	req, err := c.NewRequest("GET", ".", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if _, ok := req.Header["User-Agent"]; ok {
		t.Fatal("constructed request contains unexpected User-Agent header")
	}
}

// If a nil body is passed to greneonline.NewRequest, make sure that nil is also
// passed to http.NewRequest. In most cases, passing an io.Reader that returns
// no content is fine, since there is no difference between an HTTP request
// body that is an empty string versus one that is not set at all. However in
// certain cases, intermediate systems may treat these differently resulting in
// subtle errors.
func TestNewRequest_emptyBody(t *testing.T) {
	c := NewClient(nil)
	req, err := c.NewRequest("GET", ".", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if req.Body != nil {
		t.Fatalf("constructed request contains a non-nil Body")
	}
}

func TestNewRequest_errorForNoTrailingSlash(t *testing.T) {
	tests := []struct {
		rawurl    string
		wantError bool
	}{
		{rawurl: "https://example.com/api/v3", wantError: true},
		{rawurl: "https://example.com/api/v3/", wantError: false},
	}
	c := NewClient(nil)
	for _, test := range tests {
		u, err := url.Parse(test.rawurl)
		if err != nil {
			t.Fatalf("url.Parse returned unexpected error: %v.", err)
		}
		c.BaseURL = u
		if _, err := c.NewRequest(http.MethodGet, "test", nil); test.wantError && err == nil {
			t.Fatalf("Expected error to be returned.")
		} else if !test.wantError && err != nil {
			t.Fatalf("NewRequest returned unexpected error: %v.", err)
		}
	}
}

func TestBareDo_returnsOpenBody(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	expectedBody := "Hello from the other side !"

	mux.HandleFunc("/test-url", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, expectedBody)
	})

	ctx := context.Background()
	req, err := client.NewRequest("GET", "test-url", nil)
	if err != nil {
		t.Fatalf("client.NewRequest returned error: %v", err)
	}

	resp, err := client.BareDo(ctx, req)
	if err != nil {
		t.Fatalf("client.BareDo returned error: %v", err)
	}

	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll returned error: %v", err)
	}
	if string(got) != expectedBody {
		t.Fatalf("Expected %q, got %q", expectedBody, string(got))
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatalf("resp.Body.Close() returned error: %v", err)
	}
}

func TestDo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	body := new(foo)
	ctx := context.Background()
	_, _ = client.Do(ctx, req, body)

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_nilContext(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	req, _ := client.NewRequest("GET", ".", nil)
	//nolint:staticcheck //lint:ignore SA1012 we explicitly pass nil to test error
	_, err := client.Do(nil, req, nil)

	if !errors.Is(err, errNonNilContext) {
		t.Errorf("Expected context must be non-nil error")
	}
}

func TestDo_httpError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := context.Background()
	resp, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, got no error.")
	}
	if resp.StatusCode != 400 {
		t.Errorf("Expected HTTP 400 error, got %d status code.", resp.StatusCode)
	}
}

// Test handling of an error caused by the internal http client's Do()
// function. A redirect loop is pretty unlikely to occur within the GoCancel
// API, but does allow us to exercise the right code path.
func TestDo_redirectLoop(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "api/v1", http.StatusFound)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := context.Background()
	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected a URL error; got %#v.", err)
	}
}

func TestDo_noContent(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	var body json.RawMessage

	req, _ := client.NewRequest("GET", ".", nil)
	ctx := context.Background()
	_, err := client.Do(ctx, req, &body)
	if err != nil {
		t.Fatalf("Do returned unexpected error: %v", err)
	}
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader(`{"error": {"code": "not_found", "message": "Cannot find resource"}}`)),
	}
	err := CheckResponse(res).(*Error)

	if err == nil {
		t.Errorf("Expected error response.")
	}

	want := &Error{
		Response: res,
		Code:     "not_found",
		Message:  "Cannot find resource",
	}
	if !errors.Is(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

func TestCompareHttpResponse(t *testing.T) {
	testcases := map[string]struct {
		h1       *http.Response
		h2       *http.Response
		expected bool
	}{
		"both are nil": {
			expected: true,
		},
		"both are non nil - same StatusCode": {
			expected: true,
			h1:       &http.Response{StatusCode: 200},
			h2:       &http.Response{StatusCode: 200},
		},
		"both are non nil - different StatusCode": {
			expected: false,
			h1:       &http.Response{StatusCode: 200},
			h2:       &http.Response{StatusCode: 404},
		},
		"one is nil, other is not": {
			expected: false,
			h2:       &http.Response{},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			v := compareHttpResponse(tc.h1, tc.h2)
			if tc.expected != v {
				t.Errorf("Expected %t, got %t for (%#v, %#v)", tc.expected, v, tc.h1, tc.h2)
			}
		})
	}
}
