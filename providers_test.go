package gocancel

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestProvider_marshal(t *testing.T) {
	testJSONMarshal(t, &Provider{}, "{}")

	o := &Provider{
		ID:             String("2637d3dd-b556-409f-8f36-cd2f6d08ab77"),
		ProviderType:   String("email"),
		ProviderMethod: String("single"),
		Configuration:  &ProviderConfiguration{"foo": "bar"},
		Locales: []*ProviderLocale{
			{
				ID:            String("2637d3dd-b556-409f-8f36-cd2f6d08ab77"),
				Locale:        String("nl-NL"),
				Configuration: &ProviderConfiguration{"foo": "bar"},
				CreatedAt:     &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
				UpdatedAt:     &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
			},
		},
		CreatedAt: &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
		UpdatedAt: &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
	}
	want := `
		{
			"id":"2637d3dd-b556-409f-8f36-cd2f6d08ab77",
			"provider_type": "email",
			"provider_method": "single",
			"configuration": {
				"foo": "bar"
			},
			"locales": [
				{
					"id":"2637d3dd-b556-409f-8f36-cd2f6d08ab77",
					"locale": "nl-NL",
					"configuration": {
						"foo": "bar"
					},
					"created_at":"2021-05-27T11:49:05Z",
					"updated_at":"2021-05-27T11:49:05Z"
				}
			],
			"created_at":"2021-05-27T11:49:05Z",
			"updated_at":"2021-05-27T11:49:05Z"
		}
	`
	testJSONMarshal(t, o, want)
}

func TestProvidersService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/providers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, url.Values{"sort[created_at]": {"desc"}})

		fmt.Fprint(w, `{"providers": [{"id":"b"}]}`)
	})

	ctx := context.Background()
	opts := &ProvidersListOptions{Sort: ProvidersSortOptions{CreatedAt: "desc"}}
	providers, _, err := client.Providers.List(ctx, opts)
	if err != nil {
		t.Errorf("Providers.List returned error: %v", err)
	}

	want := []*Provider{{ID: String("b")}}
	if !cmp.Equal(providers, want) {
		t.Errorf("Providers.List returned %+v, want %+v", providers, want)
	}

	const methodName = "List"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Providers.List(ctx, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProvidersService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/providers/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"provider": {"id":"b"}}`)
	})

	ctx := context.Background()
	provider, _, err := client.Providers.Get(ctx, "b")
	if err != nil {
		t.Fatalf("Providers.Get returned error: %v", err)
	}

	want := &Provider{ID: String("b")}
	if !cmp.Equal(provider, want) {
		t.Errorf("Providers.Get returned %+v, want %+v", provider, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Providers.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Providers.Get(ctx, "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
