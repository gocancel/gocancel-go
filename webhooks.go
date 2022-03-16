package gocancel

import (
	"context"
	"fmt"
)

// WebhooksService provides access to the webhook related functions
// in the GoCancel API.
type WebhooksService service

// Webhook represents a GoCancel webhook.
type Webhook struct {
	ID        *string          `json:"id,omitempty"`
	AccountID *string          `json:"account_id,omitempty"`
	Url       *string          `json:"url,omitempty"`
	Events    []*string        `json:"events,omitempty"`
	Locales   []*string        `json:"locales,omitempty"`
	Metadata  *AccountMetadata `json:"metadata,omitempty"`
	Active    *bool            `json:"active,omitempty"`
	CreatedAt *Timestamp       `json:"created_at,omitempty"`
	UpdatedAt *Timestamp       `json:"updated_at,omitempty"`
}

func (w Webhook) String() string {
	return Stringify(w)
}

type ProductsSortOptions struct {
	CreatedAt string `url:"created_at,omitempty"`
	UpdatedAt string `url:"updated_at,omitempty"`
}

type WebhooksSortOptions struct {
	CreatedAt string `url:"created_at,omitempty"`
	UpdatedAt string `url:"updated_at,omitempty"`
}

// WebhooksListOptions specifies the optional parameters to the
// WebhooksService.List method.
type WebhooksListOptions struct {
	// Active filters webhooks by their status.
	Active  bool                `url:"active,omitempty"`
	Cursor  string              `url:"cursor,omitempty"`
	Events  []string            `url:"events,omitempty"`
	Limit   int                 `url:"limit,omitempty"`
	Locales []string            `url:"locales,omitempty"`
	Sort    WebhooksSortOptions `url:"sort,omitempty"`
}

type webhookRoot struct {
	Webhook *Webhook `json:"webhook"`
}

type webhooksRoot struct {
	Webhooks []*Webhook `json:"webhooks"`
	Metadata *Metadata  `json:"metadata"`
}

// List lists all webhooks
func (s *WebhooksService) List(ctx context.Context, opts *WebhooksListOptions) ([]*Webhook, *Response, error) {
	u, err := addOptions("api/v1/webhooks", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(webhooksRoot)
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	resp.Metadata = root.Metadata

	return root.Webhooks, resp, nil
}

// Get fetches a webhook.
func (s *WebhooksService) Get(ctx context.Context, webhook string) (*Webhook, *Response, error) {
	u := fmt.Sprintf("api/v1/webhooks/%s", webhook)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(webhookRoot)
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root.Webhook, resp, nil
}
