# gocancel-go

[![test](https://github.com/gocancel/gocancel-go/actions/workflows/test.yml/badge.svg)](https://github.com/gocancel/gocancel-go/actions/workflows/test.yml)

gocancel-go is a Go client library for accessing the GoCancel API.

You can view the GoCancel API docs here: [https://app.gocxl.com/docs](https://app.gocxl.com/docs).

## Installation

gocancel-go is compatible with modern Go releases in module mode, with Go installed:

```bash
go get github.com/gocancel/gocancel-go
```

will resolve and add the package to the current development module, along with its dependencies.

Alternatively the same can be achieved if you use import in a package:

```go
import "github.com/gocancel/gocancel-go"
```

and run `go get` without parameters.

## Usage

```go
import "github.com/gocancel/gocancel-go"
```

Construct a new client, then use the various services on the client to access different parts of the GoCancel API. For example:

```go
client := gocancel.NewClient(nil)

// list all organizations.
organizations, _, err := client.Organizations.List(context.Background(), nil)
```

### Authentication

The `gocancel-go` library does not directly handle authentication. Instead, when creating a new client, pass an `http.Client` that can handle authentication for you. The easiest and recommended way to do this is using the [oauth2](https://github.com/golang/oauth2) library, but you can always use any other library that provides an `http.Client`. If you have a OAuth2 client ID and client secret, you can use it with the oauth2 library using:

```go
import (
  "github.com/gocancel/gocancel-go"
  "golang.org/x/oauth2/clientcredentials"
)

func main() {
	ctx := context.Background()

	conf := &clientcredentials.Config{
		ClientID:     "... your client id ...",
		ClientSecret: "... your client secret ...",
		Scopes:       []string{"read:categories"},
		TokenURL:     gocancel.Endpoint.TokenURL,
	}
	tc := conf.Client(ctx)

	client := gocancel.NewClient(tc)

	// list all categories
	categories, _, err := client.Categories.List(ctx)
}
```

Note that when using an authenticated Client, all calls made by the client will include the specified OAuth client credentials. Therefore, authenticated clients should almost never be shared between different accounts.

See the [oauth2 docs](https://godoc.org/golang.org/x/oauth2) for complete instructions on using that library.

### Testing

The API client found in `gocancel-go` is HTTP based. Interactions with the HTTP API can be faked by serving up your own in-memory server within your test. One benefit of using this approach is that you don’t need to define an interface in your runtime code; you can keep using the concrete struct types returned by the client library.

To fake HTTP interactions we can make use of the `httptest` package found in the standard library. The server URL obtained from creating the test server can be passed to the client struct by overriding the `BaseURL` field. This instructs the client to route all traffic to your test server rather than the live API. Here is an example of what doing this looks like:

```go
import (
  "fmt"
  "net/http"
  "net/http/httptest"
  "net/url"

  "github.com/gocancel/gocancel-go"
)

// mux is the HTTP request multiplexer used with the test server.
mux := http.NewServeMux()

mux.HandleFunc("/api/v1/categories/foo", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"id":"foo"}`)
})

mux.HandleFunc("/api/v1/organizations/bar", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"id":"bar"}`)
})

// ts is a test HTTP server used to provide mock API responses.
ts := httptest.NewServer(mux)
defer ts.Close()

// client is the GoCancel client being tested.
client := gocancel.NewClient(nil)
url, _ := url.Parse(server.URL + "/")
// Pass the test server's URL to the API client.
client.BaseURL = url
```
