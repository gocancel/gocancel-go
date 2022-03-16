package gocancel

import (
	"context"
	"fmt"
)

type OrganizationProductsSortOptions struct {
	CreatedAt string `url:"created_at,omitempty"`
	UpdatedAt string `url:"updated_at,omitempty"`
}

type OrganizationProductsListOptions struct {
	Cursor   string                          `url:"cursor,omitempty"`
	Limit    int                             `url:"limit,omitempty"`
	Locales  []string                        `url:"locales,omitempty"`
	Metadata map[string]string               `url:"metadata,omitempty"`
	Slug     string                          `url:"slug,omitempty"`
	Sort     OrganizationProductsSortOptions `url:"sort,omitempty"`
	URL      string                          `url:"url,omitempty"`
}

// List lists all products of an organization
func (s *OrganizationsService) ListProducts(ctx context.Context, organization string, opts *OrganizationProductsListOptions) ([]*Product, *Response, error) {
	u, err := addOptions(fmt.Sprintf("api/v1/organizations/%s/products", organization), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(productsRoot)
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	resp.Metadata = root.Metadata

	return root.Products, resp, nil
}

// Get fetches a product of an organization.
func (s *OrganizationsService) GetProduct(ctx context.Context, organization string, product string) (*Product, *Response, error) {
	u := fmt.Sprintf("api/v1/organizations/%s/products/%s", organization, product)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(productRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Product, resp, nil
}
