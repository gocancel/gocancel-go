package gocancel

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestAccount_marshal(t *testing.T) {
	testJSONMarshal(t, &Account{}, "{}")

	o := &Account{
		ID:           String("7ef66283-3f9d-4cab-8f31-72f3f734652b"),
		Name:         String("ACME"),
		SandboxMode:  Bool(true),
		SandboxEmail: String("sandbox@acme.com"),
		CreatedAt:    &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
		UpdatedAt:    &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
	}
	want := `
		{
			"id":"7ef66283-3f9d-4cab-8f31-72f3f734652b",
			"name":"ACME",
			"sandbox_mode":true,
			"sandbox_email":"sandbox@acme.com",
			"created_at":"2021-05-27T11:49:05Z",
			"updated_at":"2021-05-27T11:49:05Z"
		}
	`
	testJSONMarshal(t, o, want)
}

func TestAccountsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/accounts/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"account": {"id":"b"}}`)
	})

	ctx := context.Background()
	account, _, err := client.Accounts.Get(ctx, "b")
	if err != nil {
		t.Fatalf("Accounts.Get returned error: %v", err)
	}

	want := &Account{ID: String("b")}
	if !cmp.Equal(account, want) {
		t.Errorf("Accounts.Get returned %+v, want %+v", account, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Accounts.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Accounts.Get(ctx, "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
