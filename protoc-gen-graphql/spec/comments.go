package spec

import (
	"bytes"
	"fmt"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// CommentType indicates comment format.
type CommentType int

const (
	GraphqlComment CommentType = iota
	GoComment
)

type Comments map[string]string

// makeComments makes comment map, key is paths (it depends on descriptor)
// see: https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto#L776-L798
func makeComments(d *descriptor.FileDescriptorProto) Comments {
	m := make(map[string]string)

	for _, l := range d.GetSourceCodeInfo().GetLocation() {
		paths := l.GetPath()
		// skip odd paths because it indicates field type
		if len(paths)%2 > 0 || l.GetLeadingComments() == "" {
			continue
		}
		b := new(bytes.Buffer)
		for _, p := range paths {
			b.WriteString(fmt.Sprint(p))
		}
		m[b.String()] = l.GetLeadingComments()
	}
	return Comments(m)
}
