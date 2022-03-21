package gocancel

import (
	"context"
	"fmt"
)

// ProductsService provides access to the product related functions
// in the GoCancel API.
type ProductsService service

// Product represents a product of an organization.
type Product struct {
	ID                *string          `json:"id,omitempty"`
	OrganizationID    *string          `json:"organization_id,omitempty"`
	Name              *string          `json:"name,omitempty"`
	Slug              *string          `json:"slug,omitempty"`
	Email             *string          `json:"email,omitempty"`
	URL               *string          `json:"url,omitempty"`
	Phone             *string          `json:"phone,omitempty"`
	Fax               *string          `json:"fax,omitempty"`
	Address           *Address         `json:"address,omitempty"`
	RequiresConsent   *bool            `json:"requires_consent,omitempty"`
	RequiresProofOfID *bool            `json:"requires_proof_of_id,omitempty"`
	Locales           []*ProductLocale `json:"locales,omitempty"`
	Metadata          *AccountMetadata `json:"metadata,omitempty"`
	CreatedAt         *Timestamp       `json:"created_at,omitempty"`
	UpdatedAt         *Timestamp       `json:"updated_at,omitempty"`
}

func (p Product) String() string {
	return Stringify(p)
}

// ProductLocale represents the localized variant of the product.
type ProductLocale struct {
	ID                *string            `json:"id,omitempty"`
	Name              *string            `json:"name,omitempty"`
	Slug              *string            `json:"slug,omitempty"`
	Locale            *string            `json:"locale,omitempty"`
	Email             *string            `json:"email,omitempty"`
	URL               *string            `json:"url,omitempty"`
	Phone             *string            `json:"phone,omitempty"`
	Fax               *string            `json:"fax,omitempty"`
	Address           *Address           `json:"address,omitempty"`
	RequiresConsent   *bool              `json:"requires_consent,omitempty"`
	RequiresProofOfID *bool              `json:"requires_proof_of_id,omitempty"`
	Providers         []*ProductProvider `json:"providers,omitempty"`
	LetterTemplate    *LetterTemplate    `json:"letter_template,omitempty"`
	Metadata          *AccountMetadata   `json:"metadata,omitempty"`
	CreatedAt         *Timestamp         `json:"created_at,omitempty"`
	UpdatedAt         *Timestamp         `json:"updated_at,omitempty"`
}

func (o ProductLocale) String() string {
	return Stringify(o)
}

// ProductProvider represents the provider of the product.
type ProductProvider struct {
	ID                    *string                `json:"id,omitempty"`
	ProviderID            *string                `json:"provider_id,omitempty"`
	ProviderType          *string                `json:"provider_type,omitempty"`
	ProviderMethod        *string                `json:"provider_method,omitempty"`
	ProviderConfiguration *ProviderConfiguration `json:"provider_configuration,omitempty"`
	Metadata              *AccountMetadata       `json:"metadata,omitempty"`
	CreatedAt             *Timestamp             `json:"created_at,omitempty"`
	UpdatedAt             *Timestamp             `json:"updated_at,omitempty"`
}

func (o ProductProvider) String() string {
	return Stringify(o)
}

type productRoot struct {
	Product *Product `json:"product"`
}

type productsRoot struct {
	Products []*Product `json:"products"`
	Metadata *Metadata  `json:"metadata"`
}

// Get fetches a product.
func (s *ProductsService) Get(ctx context.Context, product string) (*Product, *Response, error) {
	u := fmt.Sprintf("api/v1/products/%s", product)
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
