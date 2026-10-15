package openapiconvert

import "testing"

func TestSchemaKeyName(t *testing.T) {
	cases := []struct{ in, want string }{
		// Existing behavior, which must not change.
		{"user", "User"},
		{"user-profile", "UserProfile"},
		{"user_profile", "UserProfile"},
		{"some thing.else", "SomeThingElse"},
		{"myObject", "Myobject"},
		{"USER", "User"},
		{"2fa", "_2fa"},
		{"", "Object"},
		// Names generated for generic instantiations keep their casing.
		{"ListOfUser", "ListOfUser"},
		{"PageOfPageOfUser", "PageOfPageOfUser"},
		{"ListOfUserOrNull", "ListOfUserOrNull"},
		{"ListOfArrayOfString", "ListOfArrayOfString"},
		{"ListOfUser2", "ListOfUser2"},
	}
	for _, c := range cases {
		if got := schemaKeyName(c.in); got != c.want {
			t.Errorf("schemaKeyName(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}