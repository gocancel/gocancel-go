package gocancel

// LetterTemplate represents an embeddable letter template.
type LetterTemplate struct {
	ID        *string                `json:"id,omitempty"`
	Template  *string                `json:"template,omitempty"`
	Fields    []*LetterTemplateField `json:"fields,omitempty"`
	CreatedAt *Timestamp             `json:"created_at,omitempty"`
	UpdatedAt *Timestamp             `json:"updated_at,omitempty"`
}

func (l LetterTemplate) String() string {
	return Stringify(l)
}

type LetterTemplateField struct {
	Key      *string                      `json:"key,omitempty"`
	Type     *string                      `json:"type,omitempty"`
	Default  *string                      `json:"default,omitempty"`
	Label    *string                      `json:"label,omitempty"`
	Required *bool                        `json:"required,omitempty"`
	Position *int                         `json:"position,omitempty"`
	Options  []*LetterTemplateFieldOption `json:"options,omitempty"`
}

type LetterTemplateFieldOption struct {
	Value *string `json:"value,omitempty"`
	Label *string `json:"label,omitempty"`
}

type letterTemplateRoot struct {
	LetterTemplate *LetterTemplate `json:"letter_template"`
}
