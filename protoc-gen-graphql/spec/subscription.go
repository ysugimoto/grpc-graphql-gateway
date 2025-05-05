package spec

import (
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

func (s *Subscription) InputType() string {
	return s.Input.FullPath() // or however you qualify your Go structs
}
