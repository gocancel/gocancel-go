package gocancel

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_ListProducts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/organizations/a/products", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, url.Values{"sort[created_at]": {"asc"}})

		fmt.Fprint(w, `{"products": [{"id":"b"}]}`)
	})

	ctx := context.Background()
	opts := &OrganizationProductsListOptions{Sort: OrganizationProductsSortOptions{CreatedAt: "asc"}}
	products, _, err := client.Organizations.ListProducts(ctx, "a", opts)
	if err != nil {
		t.Errorf("Organizations.ListProducts returned error: %v", err)
	}

	want := []*Product{{ID: String("b")}}
	if !cmp.Equal(products, want) {
		t.Errorf("Organizations.ListProducts returned %+v, want %+v", products, want)
	}

	const methodName = "ListProducts"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListProducts(ctx, "a", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetProduct(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/organizations/a/products/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"product": {"id":"b"}}`)
	})

	ctx := context.Background()
	product, _, err := client.Organizations.GetProduct(ctx, "a", "b")
	if err != nil {
		t.Fatalf("Products.Get returned error: %v", err)
	}

	want := &Product{ID: String("b")}
	if !cmp.Equal(product, want) {
		t.Errorf("Products.Get returned %+v, want %+v", product, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.GetProduct(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetProduct(ctx, "a", "b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
