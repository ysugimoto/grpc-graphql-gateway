package spec

// Query spec wraps MethodDescriptorProto.
type Query struct {
	*Method
	Input  *Message
	Output *Message
}

func NewQuery(m *Method, input, output *Message) *Query {
	return &Query{
		Method: m,
		Input:  input,
		Output: output,
	}
}

func (q *Query) QueryType() string {
	if q.Method.ExposeQuery() != "" {
		field := q.Method.ExposeQueryFields(q.Output)[0]
		return field.FieldType()
	}

	typeName := PrefixType(q.Output.Name())
	if resp := q.Method.QueryResponse(); resp != nil {
		if resp.GetRepeated() {
			typeName = "graphql.NewList(" + typeName + ")"
		}
		if !resp.GetOptional() {
			typeName = "graphql.NewNonNull(" + typeName + ")"
		}
	}
	return typeName
}

func (q *Query) Args() []*Field {
	return q.Input.Fields()
}
