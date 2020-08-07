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

	messages []*Message
	services []*Service
	enums    []*Enum

	isCamel bool
}

func NewFile(d *descriptor.FileDescriptorProto, isCamel bool) *File {
	f := &File{
		descriptor: d,
		comments:   makeComments(d),

		services: make([]*Service, 0),
		messages: make([]*Message, 0),
		enums:    make([]*Enum, 0),
	}
	for i, s := range d.GetService() {
		f.services = append(f.services, NewService(s, f, 6, i))
	}
	for i, m := range d.GetMessageType() {
		f.messages = append(f.messages, f.messagesRecursive(m, []string{}, 4, i)...)
	}
	for i, e := range d.GetEnumType() {
		f.enums = append(f.enums, NewEnum(e, f, 5, i))
	}
	return f
}

func (f *File) Services() []*Service {
	return f.services
}

func (f *File) Messages() []*Message {
	return f.messages
}

func (f *File) Enums() []*Enum {
	return f.enums
}

func (f *File) messagesRecursive(d *descriptor.DescriptorProto, prefix []string, paths ...int) []*Message {
	messages := []*Message{
		NewMessage(d, f, prefix, f.isCamel, paths...),
	}

	prefix = append(prefix, d.GetName())
	for i, m := range d.GetNestedType() {
		p := make([]int, len(paths))
		copy(p, paths)
		messages = append(
			messages,
			f.messagesRecursive(m, prefix, append(p, 3, i)...)...,
		)
	}
	return messages
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
		return strings.ReplaceAll(strings.TrimSpace(v), "`", "`+\"`\"+`")
	}
	return ""
}
