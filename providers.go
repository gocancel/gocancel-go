package gocancel

import (
	"context"
	"fmt"
)

// ProvidersService provides access to the provider related functions
// in the GoCancel API.
type ProvidersService service

// Provider represents a GoCancel provider.
type Provider struct {
	ID             *string                `json:"id,omitempty"`
	ProviderType   *string                `json:"provider_type,omitempty"`
	ProviderMethod *string                `json:"provider_method,omitempty"`
	Configuration  *ProviderConfiguration `json:"configuration,omitempty"`
	Locales        []*ProviderLocale      `json:"locales,omitempty"`
	CreatedAt      *Timestamp             `json:"created_at,omitempty"`
	UpdatedAt      *Timestamp             `json:"updated_at,omitempty"`
}

func (p Provider) String() string {
	return Stringify(p)
}

type ProviderLocale struct {
	ID            *string                `json:"id,omitempty"`
	Locale        *string                `json:"locale,omitempty"`
	Configuration *ProviderConfiguration `json:"configuration,omitempty"`
	CreatedAt     *Timestamp             `json:"created_at,omitempty"`
	UpdatedAt     *Timestamp             `json:"updated_at,omitempty"`
}

type ProviderConfiguration map[string]interface{}

type ProvidersSortOptions struct {
	CreatedAt string `url:"created_at,omitempty"`
	UpdatedAt string `url:"updated_at,omitempty"`
}

type ProvidersListOptions struct {
	Cursor string               `url:"cursor,omitempty"`
	Limit  int                  `url:"limit,omitempty"`
	Sort   ProvidersSortOptions `url:"sort,omitempty"`
}

type providerRoot struct {
	Provider *Provider `json:"provider"`
}

type providersRoot struct {
	Providers []*Provider `json:"providers"`
	Metadata  *Metadata   `json:"metadata"`
}

// List lists all providers
func (s *ProvidersService) List(ctx context.Context, opts *ProvidersListOptions) ([]*Provider, *Response, error) {
	u, err := addOptions("api/v1/providers", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(providersRoot)
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	resp.Metadata = root.Metadata

	return root.Providers, resp, nil
}

// Get fetches a provider.
func (s *ProvidersService) Get(ctx context.Context, provider string) (*Provider, *Response, error) {
	u := fmt.Sprintf("api/v1/providers/%s", provider)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(providerRoot)
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root.Provider, resp, nil
}
