package extension

func MessageName(n string) string {
	return "gql_Type_" + n
}

func EnumName(n string) string {
	return "gql_Enum_" + n
}

func QueryName() string {
	return "gql_Query"
}
