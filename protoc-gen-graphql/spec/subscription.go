package spec

import (
	"log"

	"github.com/ysugimoto/grpc-graphql-gateway/graphql"
)

// Subscription wraps a streaming gRPC method for GraphQL subscription.
type Subscription struct {
	*Method
	Input   *Message
	Output  *Message
	isCamel bool
}

func NewSubscription(m *Method, input, output *Message, isCamel bool) *Subscription {
	return &Subscription{Method: m, Input: input, Output: output, isCamel: isCamel}
}

func (s *Subscription) SubscriptionName() string {
	return s.Schema.GetName()
}

// Args returns the GraphQL args for starting the stream.
// Typically, this should mirror PluckRequest().
func (s *Subscription) Args() []*Field {
	return s.PluckRequest()
}

func (s *Subscription) SubscriptionType() string {
	return PrefixType(s.Output.Name())
}

func (s *Subscription) IsCamel() bool {
	return s.isCamel
}

// Request returns the GraphQL subscription request metadata.
func (s *Subscription) Request() *graphql.GraphqlRequest {
	return s.Schema.GetRequest()
}

// Response returns the GraphQL subscription response metadata.
func (s *Subscription) Response() *graphql.GraphqlResponse {
	return s.Schema.GetResponse()
}

// PluckRequest returns the subset of input fields specified by the `plucks` option,
// or all input fields if none are specified. Pattern follows Mutation.PluckRequest().
func (s *Subscription) PluckRequest() []*Field {
	var plucks []string
	if req := s.Request(); req != nil {
		plucks = req.GetPlucks()
	}
	if len(plucks) == 0 {
		return s.Input.Fields()
	}
	var fields []*Field
	for _, f := range s.Input.Fields() {
		for _, p := range plucks {
			if p != f.Name() {
				continue
			}
			fields = append(fields, f)
		}
	}
	return fields
}

// PluckResponse returns the subset of output fields specified by the `pluck` option,
// or all output fields if none is specified. Pattern follows Mutation.PluckResponse().
func (s *Subscription) PluckResponse() []*Field {
	var pluck string
	if resp := s.Response(); resp != nil {
		pluck = resp.GetPluck()
	}
	if pluck == "" {
		return s.Output.Fields()
	}
	var fields []*Field
	for _, f := range s.Output.Fields() {
		if pluck != f.Name() {
			continue
		}
		fields = append(fields, f)
	}
	return fields
}

func (s *Subscription) IsPluckResponse() bool {
	resp := s.Response()
	if resp == nil {
		return false
	}
	return resp.GetPluck() != ""
}

func (s *Subscription) InputType() string {
	if s.Method.GoPackage() != s.Input.GoPackage() { // If input type is from a different package
		if IsGooglePackage(s.Input) { // You might need to ensure IsGooglePackage is usable here
			ptypeName, err := getImplementedPtypes(s.Input) // and getImplementedPtypes
			if err != nil {
				log.Fatalln("[PROTOC-GEN-GRAPHQL] Error (Subscription.InputType for Google Ptype):", err)
			}
			return "gql_ptypes_" + ptypeName + "." + s.Input.Name() // Assumes s.Input.Name() is simple for ptypes
		}
		return s.Input.StructName(false) // Uses the (now fixed) StructName
	}
	// If in the same package, use TypeName for correct Go type name
	return s.Input.TypeName()
}
