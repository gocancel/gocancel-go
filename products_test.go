package gocancel

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestProduct_marshal(t *testing.T) {
	testJSONMarshal(t, &Product{}, "{}")

	o := &Product{
		ID:             String("26468553-08bb-47c4-a28c-d80dec6ef3b2"),
		OrganizationID: String("f172758f-7718-41f4-95d6-d3fd931e0326"),
		Name:           String("ACME"),
		Slug:           String("acme"),
		Email:          String("contact@acme.com"),
		URL:            String("https://acme.com"),
		Phone:          String("517-234-9141"),
		Fax:            String("745-756-0818"),
		Metadata:       &AccountMetadata{"foo": "bar"},
		CreatedAt:      &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
		UpdatedAt:      &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},

		Address: &Address{
			Name:               String("Google"),
			ForAttentionOf:     String("Mr. John Doe"),
			AddressLine1:       String("1098 Alta Ave"),
			PostalCode:         String("94043"),
			Locality:           String("Mountain View"),
			AdministrativeArea: String("CA"),
			Country:            String("US"),
		},

		Locales: []*ProductLocale{
			{
				ID:        String("f38c8fab-0fa6-40b6-bb0c-6b3dfa2fec05"),
				Name:      String("ACME"),
				Slug:      String("acme"),
				Email:     String("contact@acme.com"),
				URL:       String("https://acme.com"),
				Phone:     String("517-234-9141"),
				Fax:       String("745-756-0818"),
				Locale:    String("nl-NL"),
				Metadata:  &AccountMetadata{"foo": "bar"},
				CreatedAt: &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
				UpdatedAt: &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},

				Address: &Address{
					Name:               String("Google"),
					ForAttentionOf:     String("Mr. John Doe"),
					AddressLine1:       String("1098 Alta Ave"),
					PostalCode:         String("94043"),
					Locality:           String("Mountain View"),
					AdministrativeArea: String("CA"),
					Country:            String("US"),
				},

				Providers: []*ProductProvider{
					{
						ID:                    String("f8acd284-bb6a-4933-a244-dedb9797b1d5"),
						ProviderID:            String("e88524b8-1380-41fe-b8b4-a08daabf03c8"),
						ProviderType:          String("email"),
						ProviderMethod:        String("single"),
						ProviderConfiguration: &ProviderConfiguration{"foo": "bar"},
						Metadata:              &AccountMetadata{"foo": "bar"},
						CreatedAt:             &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
						UpdatedAt:             &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
					},
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
			},
		},
	}
	want := `
		{
			"id":"26468553-08bb-47c4-a28c-d80dec6ef3b2",
			"organization_id": "f172758f-7718-41f4-95d6-d3fd931e0326",
			"name": "ACME",
			"slug": "acme",
			"email": "contact@acme.com",
			"url": "https://acme.com",
			"phone": "517-234-9141",
			"fax": "745-756-0818",
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
			"locales": [
				{
					"id": "f38c8fab-0fa6-40b6-bb0c-6b3dfa2fec05",
					"name": "ACME",
					"slug": "acme",
					"email": "contact@acme.com",
					"url": "https://acme.com",
					"phone": "517-234-9141",
					"fax": "745-756-0818",
					"locale": "nl-NL",
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
							"provider_id": "e88524b8-1380-41fe-b8b4-a08daabf03c8",
							"provider_type": "email",
							"provider_method": "single",
							"provider_configuration": {
                                "foo": "bar"
                            },
							"metadata": {
								"foo": "bar"
							},
							"created_at":"2021-05-27T11:49:05Z",
							"updated_at":"2021-05-27T11:49:05Z"
						}
					],
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
					"metadata": {
						"foo": "bar"
					},
					"created_at":"2021-05-27T11:49:05Z",
					"updated_at":"2021-05-27T11:49:05Z"
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

func TestProductsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/products/a", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"product": {"id":"a"}}`)
	})

	ctx := context.Background()
	product, _, err := client.Products.Get(ctx, "a")
	if err != nil {
		t.Fatalf("Products.Get returned error: %v", err)
	}

	want := &Product{ID: String("a")}
	if !cmp.Equal(product, want) {
		t.Errorf("Products.Get returned %+v, want %+v", product, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Products.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Products.Get(ctx, "a")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
