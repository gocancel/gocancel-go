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

func TestOrganization_marshal(t *testing.T) {
	testJSONMarshal(t, &Organization{}, "{}")

	o := &Organization{
		ID:                String("26468553-08bb-47c4-a28c-d80dec6ef3b2"),
		CategoryID:        String("f172758f-7718-41f4-95d6-d3fd931e0326"),
		Name:              String("ACME"),
		Slug:              String("acme"),
		Email:             String("contact@acme.com"),
		URL:               String("https://acme.com"),
		Phone:             String("517-234-9141"),
		Fax:               String("745-756-0818"),
		RequiresConsent:   Bool(true),
		RequiresProofOfID: Bool(true),
		Locales:           []*string{String("nl-NL")},
		Metadata:          &AccountMetadata{"foo": "bar"},
		CreatedAt:         &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
		UpdatedAt:         &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},

		Address: &Address{
			Name:               String("Google"),
			ForAttentionOf:     String("Mr. John Doe"),
			AddressLine1:       String("1098 Alta Ave"),
			PostalCode:         String("94043"),
			Locality:           String("Mountain View"),
			AdministrativeArea: String("CA"),
			Country:            String("US"),
		},
	}
	want := `
		{
			"id":"26468553-08bb-47c4-a28c-d80dec6ef3b2",
			"category_id": "f172758f-7718-41f4-95d6-d3fd931e0326",
			"name": "ACME",
			"slug": "acme",
			"email": "contact@acme.com",
			"url": "https://acme.com",
			"phone": "517-234-9141",
			"fax": "745-756-0818",
			"requires_consent": true,
			"requires_proof_of_id": true,
			"address": {
				"name": "Google",
				"for_attention_of": "Mr. John Doe",
				"address_line1": "1098 Alta Ave",
				"address_line2": null,
				"postal_code": "94043",
				"dependent_locality": null,
				"locality": "Mountain View",
				"administrative_area": "CA",
				"country": "US"
			},
			"locales": ["nl-NL"],
			"metadata": {
				"foo": "bar"
			},
			"created_at":"2021-05-27T11:49:05Z",
			"updated_at":"2021-05-27T11:49:05Z"
		}
	`
	testJSONMarshal(t, o, want)
}

func TestOrganizationLocale_marshal(t *testing.T) {
	testJSONMarshal(t, &OrganizationLocale{}, "{}")

	o := &OrganizationLocale{
		ID:                String("f38c8fab-0fa6-40b6-bb0c-6b3dfa2fec05"),
		Name:              String("ACME"),
		Slug:              String("acme"),
		Email:             String("contact@acme.com"),
		URL:               String("https://acme.com"),
		Phone:             String("517-234-9141"),
		Fax:               String("745-756-0818"),
		RequiresConsent:   Bool(true),
		RequiresProofOfID: Bool(true),
		Locale:            String("nl-NL"),
		Metadata:          &AccountMetadata{"foo": "bar"},
		CreatedAt:         &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
		UpdatedAt:         &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},

		Address: &Address{
			Name:               String("Google"),
			ForAttentionOf:     String("Mr. John Doe"),
			AddressLine1:       String("1098 Alta Ave"),
			PostalCode:         String("94043"),
			Locality:           String("Mountain View"),
			AdministrativeArea: String("CA"),
			Country:            String("US"),
		},

		Providers: []*OrganizationProvider{
			{
				ID:     String("f8acd284-bb6a-4933-a244-dedb9797b1d5"),
				Name:   String("Email"),
				Type:   String("email"),
				Method: String("single"),
			},
		},
	}
	want := `
		{
			"id": "f38c8fab-0fa6-40b6-bb0c-6b3dfa2fec05",
			"name": "ACME",
			"slug": "acme",
			"email": "contact@acme.com",
			"url": "https://acme.com",
			"phone": "517-234-9141",
			"fax": "745-756-0818",
			"locale": "nl-NL",
			"requires_consent": true,
			"requires_proof_of_id": true,
			"address": {
				"name": "Google",
				"for_attention_of": "Mr. John Doe",
				"address_line1": "1098 Alta Ave",
				"address_line2": null,
				"postal_code": "94043",
				"dependent_locality": null,
				"locality": "Mountain View",
				"administrative_area": "CA",
				"country": "US"
			},
			"providers": [
				{
					"id": "f8acd284-bb6a-4933-a244-dedb9797b1d5",
					"name": "Email",
					"type": "email",
					"method": "single"
				}
			],
			"metadata": {
				"foo": "bar"
			},
			"created_at":"2021-05-27T11:49:05Z",
			"updated_at":"2021-05-27T11:49:05Z"
		}
	`
	testJSONMarshal(t, o, want)
}

func TestOrganizationsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v2/organizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, url.Values{"sort[created_at]": {"desc"}, "locales": {"nl-NL"}})

		fmt.Fprint(w, `{"organizations": [{"id":"b"}], "metadata": {"next_cursor": "def", "previous_cursor": "abc"}}`)
	})

	ctx := context.Background()
	opts := &OrganizationsListOptions{Sort: OrganizationsSortOptions{CreatedAt: "desc"}, Locales: []string{"nl-NL"}}
	organizations, resp, err := client.Organizations.List(ctx, opts)
	if err != nil {
		t.Errorf("Organizations.List returned error: %v", err)
	}

	want := []*Organization{{ID: String("b")}}
	if !cmp.Equal(organizations, want) {
		t.Errorf("Organizations.List returned %+v, want %+v", organizations, want)
	}

	metadata := &Metadata{NextCursor: "def", PreviousCursor: "abc"}
	if !cmp.Equal(resp.Metadata, metadata) {
		t.Errorf("Organizations.List returned %+v, want %+v", resp.Metadata, metadata)
	}

	const methodName = "List"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.List(ctx, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v2/organizations/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"organization": {"id":"b"}}`)
	})

	ctx := context.Background()
	organization, _, err := client.Organizations.Get(ctx, "b")
	if err != nil {
		t.Fatalf("Organizations.Get returned error: %v", err)
	}

	want := &Organization{ID: String("b")}
	if !cmp.Equal(organization, want) {
		t.Errorf("Organizations.Get returned %+v, want %+v", organization, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.Get(ctx, "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetLocale(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v2/organizations/b/locales/c", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"organization_locale": {"id":"b"}}`)
	})

	ctx := context.Background()
	organizationLocale, _, err := client.Organizations.GetLocale(ctx, "b", "c")
	if err != nil {
		t.Fatalf("Organizations.GetLocale returned error: %v", err)
	}

	want := &OrganizationLocale{ID: String("b")}
	if !cmp.Equal(organizationLocale, want) {
		t.Errorf("Organizations.GetLocale returned %+v, want %+v", organizationLocale, want)
	}

	const methodName = "GetLocale"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetLocale(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetLocale(ctx, "b", "c")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetLetterTemplate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v2/organizations/b/locales/c/letter_template", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"letter_template": {"id":"c"}}`)
	})

	ctx := context.Background()
	letterTemplate, _, err := client.Organizations.GetLetterTemplate(ctx, "b", "c")
	if err != nil {
		t.Fatalf("Organizations.GetLetterTemplate returned error: %v", err)
	}

	want := &LetterTemplate{ID: String("c")}
	if !cmp.Equal(letterTemplate, want) {
		t.Errorf("Organizations.GetLetterTemplate returned %+v, want %+v", letterTemplate, want)
	}

	const methodName = "GetLetterTemplate"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetLetterTemplate(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetLetterTemplate(ctx, "b", "c")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
