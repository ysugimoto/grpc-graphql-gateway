package spec

import (
	"bytes"
	"fmt"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// File spec wraps FileDescriptorProto
// and this spec will be passed in all other specs in order to get
// filename, package name, etc...
type File struct {
	descriptor *descriptor.FileDescriptorProto
	comments   Comments
}

func NewFile(d *descriptor.FileDescriptorProto) *File {
	return &File{
		descriptor: d,
		comments:   makeComments(d),
	}
}

func (f *File) Services() []*Service {
	var services []*Service
	for i, s := range f.descriptor.GetService() {
		services = append(services, NewService(s, f, 6, i))
	}
	return services
}

func (f *File) Messages() []*Message {
	var messages []*Message
	for i, m := range f.descriptor.GetMessageType() {
		messages = append(messages, NewMessage(m, f, 4, i))
	}

	return messages
}

func (f *File) Enums() []*Enum {
	var enums []*Enum
	for i, e := range f.descriptor.GetEnumType() {
		enums = append(enums, NewEnum(e, f, 5, i))
	}
	return enums
}

func (f *File) Package() string {
	return f.descriptor.GetPackage()
}

func (f *File) GoPackage() string {
	var pkgName string
	if opt := f.descriptor.GetOptions(); opt == nil {
		pkgName = f.Package()
	} else if p := opt.GetGoPackage(); p == "" {
		pkgName = f.Package()
	} else {
		pkgName = p
	}
	return pkgName
}

func (f *File) Filename() string {
	return f.descriptor.GetName()
}

func (f *File) getComment(paths []int) string {
	b := new(bytes.Buffer)
	for _, p := range paths {
		b.WriteString(fmt.Sprint(p))
	}

	if v, ok := f.comments[b.String()]; ok {
		return strings.TrimSpace(v)
	}
	return ""
}
