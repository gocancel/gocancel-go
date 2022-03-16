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

func TestWebhook_marshal(t *testing.T) {
	testJSONMarshal(t, &Webhook{}, "{}")

	o := &Webhook{
		ID:        String("2637d3dd-b556-409f-8f36-cd2f6d08ab77"),
		AccountID: String("f7906b74-d9b9-4da9-950a-38e86252e328"),
		Url:       String("https://example.com"),
		Events:    []*string{String("organization.created")},
		Locales:   []*string{String("nl-NL")},
		Active:    Bool(true),
		Metadata:  &AccountMetadata{"foo": "bar"},
		CreatedAt: &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
		UpdatedAt: &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
	}
	want := `
		{
			"id":"2637d3dd-b556-409f-8f36-cd2f6d08ab77",
			"account_id": "f7906b74-d9b9-4da9-950a-38e86252e328",
			"url": "https://example.com",
			"events": ["organization.created"],
			"locales": ["nl-NL"],
			"active": true,
			"metadata": {
				"foo": "bar"
			},
			"created_at":"2021-05-27T11:49:05Z",
			"updated_at":"2021-05-27T11:49:05Z"
		}
	`
	testJSONMarshal(t, o, want)
}

func TestWebhooksService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/webhooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, url.Values{"active": {"true"}, "sort[created_at]": {"desc"}})
		fmt.Fprint(w, `{"webhooks": [{"id":"b"}], "metadata": {"next_cursor": "def", "previous_cursor": "abc"}}`)
	})

	ctx := context.Background()
	opts := &WebhooksListOptions{Active: true, Sort: WebhooksSortOptions{CreatedAt: "desc"}}
	webhooks, resp, err := client.Webhooks.List(ctx, opts)
	if err != nil {
		t.Errorf("Webhooks.List returned error: %v", err)
	}

	want := []*Webhook{{ID: String("b")}}
	if !cmp.Equal(webhooks, want) {
		t.Errorf("Webhooks.List returned %+v, want %+v", webhooks, want)
	}

	metadata := &Metadata{NextCursor: "def", PreviousCursor: "abc"}
	if !cmp.Equal(resp.Metadata, metadata) {
		t.Errorf("Webhooks.List returned %+v, want %+v", resp.Metadata, metadata)
	}

	const methodName = "List"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Webhooks.List(ctx, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestWebhooksService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/webhooks/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"webhook": {"id":"b"}}`)
	})

	ctx := context.Background()
	webhook, _, err := client.Webhooks.Get(ctx, "b")
	if err != nil {
		t.Fatalf("Webhooks.Get returned error: %v", err)
	}

	want := &Webhook{ID: String("b")}
	if !cmp.Equal(webhook, want) {
		t.Errorf("Webhooks.Get returned %+v, want %+v", webhook, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Webhooks.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Webhooks.Get(ctx, "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
