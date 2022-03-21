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

func TestCategory_marshal(t *testing.T) {
	testJSONMarshal(t, &Category{}, "{}")

	o := &Category{
		ID:              String("f172758f-7718-41f4-95d6-d3fd931e0326"),
		Name:            String("Finance"),
		Slug:            String("finance"),
		RequiresConsent: Bool(true),
		Metadata:        &AccountMetadata{"foo": "bar"},
		CreatedAt:       &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
		UpdatedAt:       &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},

		Locales: []*CategoryLocale{
			{
				ID:              String("f38c8fab-0fa6-40b6-bb0c-6b3dfa2fec05"),
				Name:            String("Financieel"),
				Slug:            String("financieel"),
				Locale:          String("nl-NL"),
				RequiresConsent: Bool(true),
				Metadata:        &AccountMetadata{"foo": "bar"},
				CreatedAt:       &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
				UpdatedAt:       &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},

				Providers: []*CategoryProvider{
					{
						ID:                    String("c61320df-9d9c-4738-b4c1-12db3f41af6c"),
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
			"id":"f172758f-7718-41f4-95d6-d3fd931e0326",
			"name": "Finance",
			"slug": "finance",
			"requires_consent": true,
			"locales": [
				{
					"id": "f38c8fab-0fa6-40b6-bb0c-6b3dfa2fec05",
					"name": "Financieel",
					"slug": "financieel",
					"requires_consent": true,
					"locale": "nl-NL",
					"providers": [
						{
							"id": "c61320df-9d9c-4738-b4c1-12db3f41af6c",
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

func TestCategoriesService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/categories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, url.Values{"sort[created_at]": {"desc"}, "locales": {"nl-NL"}})

		fmt.Fprint(w, `{"categories": [{"id":"b"}], "metadata": {"next_cursor": "def", "previous_cursor": "abc"}}`)
	})

	ctx := context.Background()
	opts := &CategoriesListOptions{Sort: CategoriesSortOptions{CreatedAt: "desc"}, Locales: []string{"nl-NL"}}
	categories, resp, err := client.Categories.List(ctx, opts)
	if err != nil {
		t.Errorf("Categories.List returned error: %v", err)
	}

	want := []*Category{{ID: String("b")}}
	if !cmp.Equal(categories, want) {
		t.Errorf("Categories.List returned %+v, want %+v", categories, want)
	}

	metadata := &Metadata{NextCursor: "def", PreviousCursor: "abc"}
	if !cmp.Equal(resp.Metadata, metadata) {
		t.Errorf("Categories.List returned %+v, want %+v", resp.Metadata, metadata)
	}

	const methodName = "List"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Categories.List(ctx, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCategoriesService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/categories/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"category": {"id":"b"}}`)
	})

	ctx := context.Background()
	category, _, err := client.Categories.Get(ctx, "b")
	if err != nil {
		t.Fatalf("Categories.Get returned error: %v", err)
	}

	want := &Category{ID: String("b")}
	if !cmp.Equal(category, want) {
		t.Errorf("Categories.Get returned %+v, want %+v", category, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Categories.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Categories.Get(ctx, "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
