package gocancel

import (
	"context"
	"fmt"
	"io"
)

// LettersService provides access to the letter related functions
// in the GoCancel API.
type LettersService service

// Letter represents a GoCancel letter.
type Letter struct {
	ID                    *string                `json:"id,omitempty"`
	AccountID             *string                `json:"account_id,omitempty"`
	OrganizationID        *string                `json:"organization_id,omitempty"`
	OrganizationName      *string                `json:"organization_name,omitempty"`
	ProductID             *string                `json:"product_id,omitempty"`
	ProductName           *string                `json:"product_name,omitempty"`
	ProviderID            *string                `json:"provider_id,omitempty"`
	ProviderConfiguration *ProviderConfiguration `json:"provider_configuration,omitempty"`
	Locale                *string                `json:"locale,omitempty"`
	State                 *string                `json:"state,omitempty"`
	ProofOfIDs            []*string              `json:"proof_of_ids,omitempty"`
	Parameters            *LetterParameters      `json:"parameters,omitempty"`
	Email                 *string                `json:"email,omitempty"`
	Fax                   *string                `json:"fax,omitempty"`
	Address               *Address               `json:"address,omitempty"`
	LetterTemplate        *LetterTemplate        `json:"letter_template,omitempty"`
	SignatureType         *string                `json:"signature_type,omitempty"`
	SignatureData         *string                `json:"signature_data,omitempty"`
	SandboxMode           *bool                  `json:"sandbox_mode,omitempty"`
	SandboxEmail          *string                `json:"sandbox_email,omitempty"`
	Metadata              *Metadata              `json:"metadata,omitempty"`
	CreatedAt             *Timestamp             `json:"created_at,omitempty"`
	UpdatedAt             *Timestamp             `json:"updated_at,omitempty"`
}

func (l Letter) String() string {
	return Stringify(l)
}

// LetterParameters represents key-value parameters for a letter.
type LetterParameters map[string]interface{}

type LettersSortOptions struct {
	CreatedAt string `url:"created_at,omitempty"`
	UpdatedAt string `url:"updated_at,omitempty"`
}

type LettersListOptions struct {
	Sort LettersSortOptions `url:"sort,omitempty"`
}

type letterRoot struct {
	Letter *Letter `json:"letter"`
}

type lettersRoot struct {
	Letters []*Letter `json:"letters"`
}

// List lists all letters.
func (s *LettersService) List(ctx context.Context, opts *LettersListOptions) ([]*Letter, *Response, error) {
	u, err := addOptions("api/v1/letters", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(lettersRoot)
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root.Letters, resp, nil
}

// CreateLetterRequest represents a `create letter` request.
type CreateLetterRequest struct {
	OrganizationID string           `json:"organization_id,omitempty"`
	ProductID      string           `json:"product_id,omitempty"`
	ProviderID     string           `json:"provider_id,omitempty"`
	Locale         string           `json:"locale,omitempty"`
	Parameters     LetterParameters `json:"parameters,omitempty"`
	ProofOfIDs     []string         `json:"proof_of_ids,omitempty"`
	SignatureType  string           `json:"signature_type,omitempty"`
	SignatureData  string           `json:"signature_data,omitempty"`
	Metadata       Metadata         `json:"metadata,omitempty"`
	Drafted        bool             `json:"drafted,omitempty"`
}

// Create creates a letter.
func (s *LettersService) Create(ctx context.Context, request *CreateLetterRequest) (*Letter, *Response, error) {
	req, err := s.client.NewRequest("POST", "api/v1/letters", request)
	if err != nil {
		return nil, nil, err
	}

	root := new(letterRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Letter, resp, nil
}

// Get fetches a letter.
func (s *LettersService) Get(ctx context.Context, letter string) (*Letter, *Response, error) {
	u := fmt.Sprintf("api/v1/letters/%s", letter)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", mediaTypeLetterDocument)

	root := new(letterRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Letter, resp, nil
}

// DownloadDocument fetches a letter document as binary stream.
func (s *LettersService) DownloadDocument(ctx context.Context, letter string) (io.ReadCloser, *Response, error) {
	u := fmt.Sprintf("api/v1/letters/%s/document", letter)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", mediaTypeLetterDocument)

	resp, err := s.client.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	if err := CheckResponse(resp); err != nil {
		resp.Body.Close()
		return nil, nil, err
	}

	return resp.Body, newResponse(resp), nil
}

// MarkLetterAsDraftedRequest represents a `mark letter as drafted` request.
type MarkLetterAsDraftedRequest struct {
	ProviderKey string `json:"provider_key,omitempty"`
}

// MarkAsDrafted marks a letter as drafted.
func (s *LettersService) MarkAsDrafted(ctx context.Context, letter string, request *MarkLetterAsDraftedRequest) (*Letter, *Response, error) {
	u := fmt.Sprintf("api/v1/letters/%s/mark_as_drafted", letter)
	req, err := s.client.NewRequest("POST", u, request)
	if err != nil {
		return nil, nil, err
	}

	root := new(letterRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Letter, resp, nil
}
