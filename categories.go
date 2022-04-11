package gocancel

import (
	"context"
	"fmt"
)

// CategoriesService provides access to the category related functions
// in the GoCancel API.
type CategoriesService service

// Category represents a GoCancel category.
type Category struct {
	ID              *string          `json:"id,omitempty"`
	Name            *string          `json:"name,omitempty"`
	Slug            *string          `json:"slug,omitempty"`
	RequiresConsent *bool            `json:"requires_consent,omitempty"`
	Locales         []*string        `json:"locales,omitempty"`
	Metadata        *AccountMetadata `json:"metadata,omitempty"`
	CreatedAt       *Timestamp       `json:"created_at,omitempty"`
	UpdatedAt       *Timestamp       `json:"updated_at,omitempty"`
}

func (c Category) String() string {
	return Stringify(c)
}

// CategoryLocale represents the localized variant of the category.
type CategoryLocale struct {
	ID              *string             `json:"id,omitempty"`
	Name            *string             `json:"name,omitempty"`
	Slug            *string             `json:"slug,omitempty"`
	RequiresConsent *bool               `json:"requires_consent,omitempty"`
	Locale          *string             `json:"locale,omitempty"`
	Providers       []*CategoryProvider `json:"providers,omitempty"`
	LetterTemplate  *LetterTemplate     `json:"letter_template,omitempty"`
	Metadata        *AccountMetadata    `json:"metadata,omitempty"`
	CreatedAt       *Timestamp          `json:"created_at,omitempty"`
	UpdatedAt       *Timestamp          `json:"updated_at,omitempty"`
}

func (o CategoryLocale) String() string {
	return Stringify(o)
}

// CategoryProvider represents the provider of the category.
type CategoryProvider struct {
	ID     *string `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	Type   *string `json:"type,omitempty"`
	Method *string `json:"method,omitempty"`
}

func (c CategoryProvider) String() string {
	return Stringify(c)
}

type CategoriesSortOptions struct {
	CreatedAt string `url:"created_at,omitempty"`
	UpdatedAt string `url:"updated_at,omitempty"`
}

type CategoriesListOptions struct {
	Cursor   string                `url:"cursor,omitempty"`
	Limit    int                   `url:"limit,omitempty"`
	Locales  []string              `url:"locales,omitempty"`
	Metadata map[string]string     `url:"metadata,omitempty"`
	Slug     string                `url:"slug,omitempty"`
	Sort     CategoriesSortOptions `url:"sort,omitempty"`
}

type categoryRoot struct {
	Category *Category `json:"category"`
}

type categoryLocaleRoot struct {
	CategoryLocale *CategoryLocale `json:"category_locale"`
}

type categoriesRoot struct {
	Categories []*Category `json:"categories"`
	Metadata   *Metadata   `json:"metadata"`
}

// List lists all categories
func (s *CategoriesService) List(ctx context.Context, opts *CategoriesListOptions) ([]*Category, *Response, error) {
	u, err := addOptions("api/v2/categories", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(categoriesRoot)
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	resp.Metadata = root.Metadata

	return root.Categories, resp, nil
}

// Get fetches a category.
func (s *CategoriesService) Get(ctx context.Context, category string) (*Category, *Response, error) {
	u := fmt.Sprintf("api/v2/categories/%s", category)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(categoryRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Category, resp, nil
}

// GetLocale fetches a localized category.
func (s *CategoriesService) GetLocale(ctx context.Context, category string, locale string) (*CategoryLocale, *Response, error) {
	u := fmt.Sprintf("api/v2/categories/%s/locales/%s", category, locale)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(categoryLocaleRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.CategoryLocale, resp, nil
}

// GetLetterTemplate fetches a localized letter template for a category.
func (s *CategoriesService) GetLetterTemplate(ctx context.Context, category string, locale string) (*LetterTemplate, *Response, error) {
	u := fmt.Sprintf("api/v2/categories/%s/locales/%s/letter_template", category, locale)
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
