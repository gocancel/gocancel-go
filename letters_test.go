package gocancel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestLetter_marshal(t *testing.T) {
	testJSONMarshal(t, &Letter{}, "{}")

	o := &Letter{
		ID:                    String("26468553-08bb-47c4-a28c-d80dec6ef3b2"),
		AccountID:             String("f172758f-7718-41f4-95d6-d3fd931e0326"),
		OrganizationID:        String("f172758f-7718-41f4-95d6-d3fd931e0326"),
		OrganizationName:      String("Foo"),
		ProductID:             String("f172758f-7718-41f4-95d6-d3fd931e0326"),
		ProductName:           String("Bar"),
		ProviderID:            String("f172758f-7718-41f4-95d6-d3fd931e0326"),
		ProviderConfiguration: &ProviderConfiguration{"foo": "bar"},
		Locale:                String("nl-NL"),
		State:                 String("generating"),
		ProofOfIDs:            []*string{String("1d7e5cf6-a871-48cd-b98a-1ecc6acbda96"), String("0971e527-ea0d-4ba2-a87b-0e8e8d4f83a2")},
		Email:                 String("cancellations@foo.com"),
		Fax:                   String("71-336-4530"),
		Parameters:            &LetterParameters{"foo": "bar"},
		SignatureType:         String("text"),
		SignatureData:         String("John Doe"),
		SandboxMode:           Bool(false),
		SandboxEmail:          String("sandbox@example.com"),
		Metadata:              &Metadata{"foo": "bar"},
		CreatedAt:             &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
		UpdatedAt:             &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},

		Address: &Address{
			AddressLine1:       String("1098 Alta Ave"),
			PostalCode:         String("94043"),
			Locality:           String("Mountain View"),
			AdministrativeArea: String("CA"),
			Country:            String("US"),
		},

		LetterTemplate: &LetterTemplate{
			Template: String("Dear {{ name }}"),
			Fields: []*LetterTemplateField{
				{
					Key:      String("name"),
					Type:     String("string"),
					Label:    String("Name"),
					Required: Bool(true),
					Position: Int(0),

					Options: []*LetterTemplateFieldOption{
						{
							Value: String("foo"),
							Label: String("bar"),
						},
					},
				},
			},
		},
	}
	want := `
		{
			"id":"26468553-08bb-47c4-a28c-d80dec6ef3b2",
			"account_id": "f172758f-7718-41f4-95d6-d3fd931e0326",
			"organization_id": "f172758f-7718-41f4-95d6-d3fd931e0326",
			"organization_name": "Foo",
			"product_id": "f172758f-7718-41f4-95d6-d3fd931e0326",
			"product_name": "Bar",
            "provider_id": "f172758f-7718-41f4-95d6-d3fd931e0326",
			"provider_configuration": {
				"foo": "bar"
			},
			"locale": "nl-NL",
			"state": "generating",
			"proof_of_ids": [
                "1d7e5cf6-a871-48cd-b98a-1ecc6acbda96",
                "0971e527-ea0d-4ba2-a87b-0e8e8d4f83a2"
            ],
			"email": "cancellations@foo.com",
            "fax": "71-336-4530",
			"address": {
				"address_line1": "1098 Alta Ave",
				"address_line2": null,
				"postal_code": "94043",
				"dependent_locality": null,
				"locality": "Mountain View",
				"administrative_area": "CA",
				"country": "US"
			},
			"parameters": {
				"foo": "bar"
			},
			"letter_template": {
                "template": "Dear {{ name }}",
                "fields": [
                    {
                        "key": "name",
                        "type": "string",
                        "label": "Name",
                        "required": true,
                        "position": 0,
                        "options": [{"value": "foo", "label": "bar"}]
                    }
                ]
            },
			"signature_type": "text",
			"signature_data": "John Doe",
			"sandbox_mode": false,
			"sandbox_email": "sandbox@example.com",
			"metadata": {
				"foo": "bar"
			},
			"created_at":"2021-05-27T11:49:05Z",
			"updated_at":"2021-05-27T11:49:05Z"
		}
	`
	testJSONMarshal(t, o, want)
}

func TestLettersService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/letters", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, url.Values{"sort[created_at]": {"asc"}})

		fmt.Fprint(w, `{"letters": [{"id":"b"}]}`)
	})

	ctx := context.Background()
	opts := &LettersListOptions{Sort: LettersSortOptions{CreatedAt: "asc"}}
	letters, _, err := client.Letters.List(ctx, opts)
	if err != nil {
		t.Errorf("Letters.List returned error: %v", err)
	}

	want := []*Letter{{ID: String("b")}}
	if !cmp.Equal(letters, want) {
		t.Errorf("Letters.List returned %+v, want %+v", letters, want)
	}

	const methodName = "List"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Letters.List(ctx, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestLettersService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &CreateLetterRequest{OrganizationID: "foo"}

	mux.HandleFunc("/api/v1/letters", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		v := new(CreateLetterRequest)
		_ = json.NewDecoder(r.Body).Decode(v)

		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"letter": {"id":"b"}}`)
	})

	ctx := context.Background()
	letter, _, err := client.Letters.Create(ctx, input)
	if err != nil {
		t.Fatalf("Letters.Create returned error: %v", err)
	}

	want := &Letter{ID: String("b")}
	if !cmp.Equal(letter, want) {
		t.Errorf("Letters.Create returned %+v, want %+v", letter, want)
	}

	const methodName = "Create"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Letters.Create(ctx, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestLettersService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/letters/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"letter": {"id":"b"}}`)
	})

	ctx := context.Background()
	letter, _, err := client.Letters.Get(ctx, "b")
	if err != nil {
		t.Fatalf("Letters.Get returned error: %v", err)
	}

	want := &Letter{ID: String("b")}
	if !cmp.Equal(letter, want) {
		t.Errorf("Letters.Get returned %+v, want %+v", letter, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Letters.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Letters.Get(ctx, "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestLettersService_DownloadDocument(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/letters/b/document", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeLetterDocument)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=hello-world.pdf")
		fmt.Fprint(w, "Hello World")
	})

	ctx := context.Background()
	reader, _, err := client.Letters.DownloadDocument(ctx, "b")
	if err != nil {
		t.Fatalf("Letters.DownloadDocument returned error: %v", err)
	}

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Errorf("Letters.DownloadDocument returned bad reader: %v", err)
	}

	want := []byte("Hello World")
	if !bytes.Equal(content, want) {
		t.Errorf("Letters.DownloadDocument returned %+v, want %+v", content, want)
	}

	const methodName = "DownloadDocument"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Letters.DownloadDocument(ctx, "\n")
		return err
	})
}

func TestLettersService_MarkAsDrafted(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &MarkLetterAsDraftedRequest{}

	mux.HandleFunc("/api/v1/letters/b/mark_as_drafted", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		v := new(MarkLetterAsDraftedRequest)
		_ = json.NewDecoder(r.Body).Decode(v)

		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"letter": {"id":"b"}}`)
	})

	ctx := context.Background()
	letter, _, err := client.Letters.MarkAsDrafted(ctx, "b", input)
	if err != nil {
		t.Fatalf("Letters.MarkAsDrafted returned error: %v", err)
	}

	want := &Letter{ID: String("b")}
	if !cmp.Equal(letter, want) {
		t.Errorf("Letters.MarkAsDrafted returned %+v, want %+v", letter, want)
	}

	const methodName = "MarkAsDrafted"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Letters.MarkAsDrafted(ctx, "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Letters.MarkAsDrafted(ctx, "b", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
