package gocancel

import (
	"context"
	"fmt"
)

// OrganizationsService provides access to the organization related functions
// in the GoCancel API.
type OrganizationsService service

// Organization represents a GoCancel organization.
type Organization struct {
	ID         *string               `json:"id,omitempty"`
	CategoryID *string               `json:"category_id,omitempty"`
	Name       *string               `json:"name,omitempty"`
	Slug       *string               `json:"slug,omitempty"`
	Email      *string               `json:"email,omitempty"`
	URL        *string               `json:"url,omitempty"`
	Phone      *string               `json:"phone,omitempty"`
	Fax        *string               `json:"fax,omitempty"`
	Address    *Address              `json:"address,omitempty"`
	Locales    []*OrganizationLocale `json:"locales,omitempty"`
	Metadata   *AccountMetadata      `json:"metadata,omitempty"`
	CreatedAt  *Timestamp            `json:"created_at,omitempty"`
	UpdatedAt  *Timestamp            `json:"updated_at,omitempty"`
}

func (o Organization) String() string {
	return Stringify(o)
}

// OrganizationLocale represents the localized variant of the organization.
type OrganizationLocale struct {
	ID             *string                 `json:"id,omitempty"`
	Name           *string                 `json:"name,omitempty"`
	Slug           *string                 `json:"slug,omitempty"`
	Locale         *string                 `json:"locale,omitempty"`
	Email          *string                 `json:"email,omitempty"`
	URL            *string                 `json:"url,omitempty"`
	Phone          *string                 `json:"phone,omitempty"`
	Fax            *string                 `json:"fax,omitempty"`
	Address        *Address                `json:"address,omitempty"`
	Providers      []*OrganizationProvider `json:"providers,omitempty"`
	LetterTemplate *LetterTemplate         `json:"letter_template,omitempty"`
	Metadata       *AccountMetadata        `json:"metadata,omitempty"`
	CreatedAt      *Timestamp              `json:"created_at,omitempty"`
	UpdatedAt      *Timestamp              `json:"updated_at,omitempty"`
}

func (o OrganizationLocale) String() string {
	return Stringify(o)
}

// OrganizationProvider represents the provider of the organization.
type OrganizationProvider struct {
	ID                    *string                `json:"id,omitempty"`
	ProviderID            *string                `json:"provider_id,omitempty"`
	ProviderType          *string                `json:"provider_type,omitempty"`
	ProviderMethod        *string                `json:"provider_method,omitempty"`
	ProviderConfiguration *ProviderConfiguration `json:"provider_configuration,omitempty"`
	Metadata              *AccountMetadata       `json:"metadata,omitempty"`
	CreatedAt             *Timestamp             `json:"created_at,omitempty"`
	UpdatedAt             *Timestamp             `json:"updated_at,omitempty"`
}

func (o OrganizationProvider) String() string {
	return Stringify(o)
}

type OrganizationsSortOptions struct {
	CreatedAt string `url:"created_at,omitempty"`
	UpdatedAt string `url:"updated_at,omitempty"`
}

type OrganizationsListOptions struct {
	Category string                   `url:"category,omitempty"`
	Cursor   string                   `url:"cursor,omitempty"`
	Limit    int                      `url:"limit,omitempty"`
	Locales  []string                 `url:"locales,omitempty"`
	Metadata map[string]string        `url:"metadata,omitempty"`
	Slug     string                   `url:"slug,omitempty"`
	Sort     OrganizationsSortOptions `url:"sort,omitempty"`
	URL      string                   `url:"url,omitempty"`
}

type organizationRoot struct {
	Organization *Organization `json:"organization"`
}

type organizationsRoot struct {
	Organizations []*Organization `json:"organizations"`
	Metadata      *Metadata       `json:"metadata"`
}

// List lists all organizations
func (s *OrganizationsService) List(ctx context.Context, opts *OrganizationsListOptions) ([]*Organization, *Response, error) {
	u, err := addOptions("api/v1/organizations", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(organizationsRoot)
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	resp.Metadata = root.Metadata

	return root.Organizations, resp, nil
}

// Get fetches a organization.
func (s *OrganizationsService) Get(ctx context.Context, organization string) (*Organization, *Response, error) {
	u := fmt.Sprintf("api/v1/organizations/%s", organization)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(organizationRoot)
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root.Organization, resp, nil
}
