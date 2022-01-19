package gocancel

import (
	"context"
	"fmt"
	"io"
)

// DownloadProofOfID downloads a proof of ID for a letter
func (s *LettersService) DownloadProofOfID(ctx context.Context, letter string, proof_of_id string) (io.ReadCloser, *Response, error) {
	u := fmt.Sprintf("api/v1/letters/%v/proof_of_ids/%v", letter, proof_of_id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", mediaTypeLetterProofOfID)

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
