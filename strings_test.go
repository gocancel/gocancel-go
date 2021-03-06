package gocancel

import (
	"fmt"
	"testing"
)

func TestStringify(t *testing.T) {
	var nilPointer *string

	var tests = []struct {
		in  interface{}
		out string
	}{
		// basic types
		{"foo", `"foo"`},
		{123, `123`},
		{1.5, `1.5`},
		{false, `false`},
		{
			[]string{"a", "b"},
			`["a" "b"]`,
		},
		{
			struct {
				A []string
			}{nil},
			// nil slice is skipped
			`{}`,
		},
		{
			struct {
				A string
			}{"foo"},
			// structs not of a named type get no prefix
			`{A:"foo"}`,
		},

		// pointers
		{nilPointer, `<nil>`},
		{String("foo"), `"foo"`},
		{Int(123), `123`},
		{Bool(false), `false`},
		{
			[]*string{String("a"), String("b")},
			`["a" "b"]`,
		},

		// actual GoCancel API structs
		{
			Category{ID: String("26468553-08bb-47c4-a28c-d80dec6ef3b2"), Name: String("foo")},
			`gocancel.Category{ID:"26468553-08bb-47c4-a28c-d80dec6ef3b2", Name:"foo"}`,
		},
		{
			Organization{Locales: []*OrganizationLocale{{ID: String("f38c8fab-0fa6-40b6-bb0c-6b3dfa2fec05")}}},
			`gocancel.Organization{Locales:[gocancel.OrganizationLocale{ID:"f38c8fab-0fa6-40b6-bb0c-6b3dfa2fec05"}]}`,
		},
	}

	for i, tt := range tests {
		s := Stringify(tt.in)
		if s != tt.out {
			t.Errorf("%d. Stringify(%q) => %q, want %q", i, tt.in, s, tt.out)
		}
	}
}

// Directly test the String() methods on various GoCancel types. We don't do an
// exaustive test of all the various field types, since TestStringify() above
// takes care of that. Rather, we just make sure that Stringify() is being
// used to build the strings, which we do by verifying that pointers are
// stringified as their underlying value.
func TestString(t *testing.T) {
	var tests = []struct {
		in  interface{}
		out string
	}{
		{Category{ID: String("n")}, `gocancel.Category{ID:"n"}`},
		{Organization{ID: String("n")}, `gocancel.Organization{ID:"n"}`},
		{OrganizationLocale{ID: String("n")}, `gocancel.OrganizationLocale{ID:"n"}`},
	}

	for i, tt := range tests {
		s := tt.in.(fmt.Stringer).String()
		if s != tt.out {
			t.Errorf("%d. String() => %q, want %q", i, tt.in, tt.out)
		}
	}
}
