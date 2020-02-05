package spec

// PrefixType adds prefix to avoid conflicting name
func PrefixType(name string) string {
	return "gql__type_" + name
}

// PrefixEnum adds prefix to avoid conflicting name
func PrefixEnum(name string) string {
	return "gql__enum_" + name
}

// PrefixInput adds prefix to avoid conflicting name
func PrefixInput(name string) string {
	return "gql__input_" + name
}
