package gocancel

import (
	"testing"
	"time"
)

func TestLetterTemplate_marshal(t *testing.T) {
	testJSONMarshal(t, &LetterTemplate{}, "{}")

	o := &LetterTemplate{
		ID:        String("26468553-08bb-47c4-a28c-d80dec6ef3b2"),
		Template:  String("Dear {{ name }}"),
		CreatedAt: &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},
		UpdatedAt: &Timestamp{time.Date(2021, time.May, 27, 11, 49, 05, 0, time.UTC)},

		Fields: []*LetterTemplateField{
			{
				Key:      String("name"),
				Type:     String("string"),
				Label:    String("Name"),
				Required: Bool(true),
				Position: Int(0),

				Options: []*LetterTemplateFieldOption{
					{
						Value: String("foo"),
						Label: String("bar"),
					},
				},
			},
		},
	}
	want := `
		{
			"id":"26468553-08bb-47c4-a28c-d80dec6ef3b2",
			"template": "Dear {{ name }}",
			"fields": [
				{
					"key": "name",
					"type": "string",
					"label": "Name",
					"required": true,
					"position": 0,
					"options": [{"value": "foo", "label": "bar"}]
				}
			],
			"created_at":"2021-05-27T11:49:05Z",
			"updated_at":"2021-05-27T11:49:05Z"
		}
	`
	testJSONMarshal(t, o, want)
}
