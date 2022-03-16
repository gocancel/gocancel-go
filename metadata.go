package gocancel

// Metadata represents the cursors for list responses.
type Metadata struct {
	NextCursor     string `json:"next_cursor,omitempty"`
	PreviousCursor string `json:"previous_cursor,omitempty"`
}
