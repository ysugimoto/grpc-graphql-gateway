package spec

import (
	"strings"
)

// PrefixType adds prefix to avoid conflicting name
func PrefixType(name string) string {
	return "Gql__type_" + name
}

// PrefixEnum adds prefix to avoid conflicting name
func PrefixEnum(name string) string {
	return "Gql__enum_" + name
}

// PrefixInput adds prefix to avoid conflicting name
func PrefixInput(name string) string {
	return "Gql__input_" + name
}

// FormatTypePath formats expected type path name
// e.g .book.SomeMessage.Nested -> book_SomeMessage_Nested
func FormatTypePath(typeName string) string {
	return strings.TrimPrefix(
		strings.ReplaceAll(typeName, ".", "_"),
		"_",
	)
}
