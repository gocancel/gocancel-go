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
	ID                *string          `json:"id,omitempty"`
	CategoryID        *string          `json:"category_id,omitempty"`
	Name              *string          `json:"name,omitempty"`
	Slug              *string          `json:"slug,omitempty"`
	Email             *string          `json:"email,omitempty"`
	URL               *string          `json:"url,omitempty"`
	Phone             *string          `json:"phone,omitempty"`
	Fax               *string          `json:"fax,omitempty"`
	Address           *Address         `json:"address,omitempty"`
	RequiresConsent   *bool            `json:"requires_consent,omitempty"`
	RequiresProofOfID *bool            `json:"requires_proof_of_id,omitempty"`
	Locales           []*string        `json:"locales,omitempty"`
	Metadata          *AccountMetadata `json:"metadata,omitempty"`
	CreatedAt         *Timestamp       `json:"created_at,omitempty"`
	UpdatedAt         *Timestamp       `json:"updated_at,omitempty"`
}

func (o Organization) String() string {
	return Stringify(o)
}

// OrganizationLocale represents the localized variant of the organization.
type OrganizationLocale struct {
	ID                *string                 `json:"id,omitempty"`
	Name              *string                 `json:"name,omitempty"`
	Slug              *string                 `json:"slug,omitempty"`
	Locale            *string                 `json:"locale,omitempty"`
	Email             *string                 `json:"email,omitempty"`
	URL               *string                 `json:"url,omitempty"`
	Phone             *string                 `json:"phone,omitempty"`
	Fax               *string                 `json:"fax,omitempty"`
	Address           *Address                `json:"address,omitempty"`
	RequiresConsent   *bool                   `json:"requires_consent,omitempty"`
	RequiresProofOfID *bool                   `json:"requires_proof_of_id,omitempty"`
	Providers         []*OrganizationProvider `json:"providers,omitempty"`
	LetterTemplate    *LetterTemplate         `json:"letter_template,omitempty"`
	Metadata          *AccountMetadata        `json:"metadata,omitempty"`
	CreatedAt         *Timestamp              `json:"created_at,omitempty"`
	UpdatedAt         *Timestamp              `json:"updated_at,omitempty"`
}

func (o OrganizationLocale) String() string {
	return Stringify(o)
}

// OrganizationProvider represents the provider of the organization.
type OrganizationProvider struct {
	ID        *string          `json:"id,omitempty"`
	Name      *string          `json:"name,omitempty"`
	Type      *string          `json:"type,omitempty"`
	Method    *string          `json:"method,omitempty"`
	Metadata  *AccountMetadata `json:"metadata,omitempty"`
	CreatedAt *Timestamp       `json:"created_at,omitempty"`
	UpdatedAt *Timestamp       `json:"updated_at,omitempty"`
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

type organizationLocaleRoot struct {
	OrganizationLocale *OrganizationLocale `json:"organization_locale"`
}

type organizationsRoot struct {
	Organizations []*Organization `json:"organizations"`
	Metadata      *Metadata       `json:"metadata"`
}

// List lists all organizations
func (s *OrganizationsService) List(ctx context.Context, opts *OrganizationsListOptions) ([]*Organization, *Response, error) {
	u, err := addOptions("api/v2/organizations", opts)
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

// Get fetches an organization.
func (s *OrganizationsService) Get(ctx context.Context, organization string) (*Organization, *Response, error) {
	u := fmt.Sprintf("api/v2/organizations/%s", organization)
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

// GetLocale fetches a localized organization.
func (s *OrganizationsService) GetLocale(ctx context.Context, organization string, locale string) (*OrganizationLocale, *Response, error) {
	u := fmt.Sprintf("api/v2/organizations/%s/locales/%s", organization, locale)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(organizationLocaleRoot)
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root.OrganizationLocale, resp, nil
}

// GetLetterTemplate fetches a localized letter template for an organization.
func (s *OrganizationsService) GetLetterTemplate(ctx context.Context, organization string, locale string) (*LetterTemplate, *Response, error) {
	u := fmt.Sprintf("api/v2/organizations/%s/locales/%s/letter_template", organization, locale)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(letterTemplateRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.LetterTemplate, resp, nil
}
