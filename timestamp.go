package gocancel

import (
	"time"
)

// Timestamp represents a time that can be unmarshalled from a JSON string
// formatted as an RFC3339 timestamp. All
// exported methods of time.Time can be called on Timestamp.
type Timestamp struct {
	time.Time
}

func (t Timestamp) String() string {
	return t.Time.String()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// Time is expected in RFC3339 format.
func (t *Timestamp) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	t.Time, err = time.Parse(`"`+time.RFC3339+`"`, str)
	return
}

// Equal reports whether t and u are equal based on time.Equal
func (t Timestamp) Equal(u Timestamp) bool {
	return t.Time.Equal(u.Time)
}
