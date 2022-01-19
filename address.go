package gocancel

// Address represents an embeddable address.
type Address struct {
	Name               *string `json:"name,omitempty"`
	ForAttentionOf     *string `json:"for_attention_of,omitempty"`
	AddressLine1       *string `json:"address_line1,omitempty"`
	AddressLine2       *string `json:"address_line2,omitempty"`
	PostalCode         *string `json:"postal_code,omitempty"`
	DependentLocality  *string `json:"dependent_locality,omitempty"`
	Locality           *string `json:"locality,omitempty"`
	AdministrativeArea *string `json:"administrative_area,omitempty"`
	Country            *string `json:"country,omitempty"`
}

func (a Address) String() string {
	return Stringify(a)
}
