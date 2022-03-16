package gocancel

import (
	"context"
	"fmt"
)

// AccountsService provides access to the account related functions
// in the GoCancel API.
type AccountsService service

// Account represents a GoCancel account.
type Account struct {
	ID           *string    `json:"id,omitempty"`
	Name         *string    `json:"name,omitempty"`
	SandboxMode  *bool      `json:"sandbox_mode,omitempty"`
	SandboxEmail *string    `json:"sandbox_email,omitempty"`
	CreatedAt    *Timestamp `json:"created_at,omitempty"`
	UpdatedAt    *Timestamp `json:"updated_at,omitempty"`
}

func (w Account) String() string {
	return Stringify(w)
}

// AccountMetadata represents key-value attributes for a specific account.
type AccountMetadata map[string]interface{}

type accountRoot struct {
	Account *Account `json:"account"`
}

// Get fetches an account.
func (s *AccountsService) Get(ctx context.Context, account string) (*Account, *Response, error) {
	u := fmt.Sprintf("api/v1/accounts/%s", account)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(accountRoot)
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root.Account, resp, nil
}
