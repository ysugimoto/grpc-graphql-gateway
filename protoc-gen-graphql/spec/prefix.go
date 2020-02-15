package spec

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
