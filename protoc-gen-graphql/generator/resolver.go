package generator

import (
	"errors"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
)

// Resolver is struct for resolve messages, enumurations by package name.
// In protobuf, each message, fields, etc.. has nested message and
// can be defined with other package's message, therefore it's hard to resolve recursively.
// To resolve them, this struct keeps "all" message, enums by each names which can get from descriptors
// and resolve recursively with keeping uniqueness.
type Resolver struct {
	// factory stacks
	messages map[string]*spec.Message
	enums    map[string]*spec.Enum
}

// Create all stacks on constructing,
// so we should instantiate only once.
func NewResolver(files []*spec.File) *Resolver {
	messages := make(map[string]*spec.Message)
	enums := make(map[string]*spec.Enum)

	for _, f := range files {
		for _, m := range f.Messages() {
			messages[m.FullPath()] = m
		}
		for _, e := range f.Enums() {
			enums[e.FullPath()] = e
		}
	}

	return &Resolver{
		messages: messages,
		enums:    enums,
	}
}

// Find message from pakcage name.
// Trick: this function will be passed to some builders.
func (r *Resolver) Find(t string) *spec.Message {
	if v, ok := r.messages[t]; !ok {
		return nil
	} else {
		return v
	}
}

// Find enum from pakcage name.
// Trick: this function will be passed to some builders.
func (r *Resolver) FindEnum(t string) *spec.Enum {
	if v, ok := r.enums[t]; !ok {
		return nil
	} else {
		return v
	}
}

// ResolveTypes resolves all types which is used only in whole queries and mutations.
func (r *Resolver) ResolveTypes(
	queries []*spec.Method,
	mutations []*spec.Method,
) (
	types []*spec.Message,
	enums []*spec.Enum,
	inputs []*spec.Message,
	packages []*spec.Package,
	resolveErr error,
) {
	var methods []*spec.Method
	methods = append(methods, queries...)
	methods = append(methods, mutations...)

	cache := NewCache()

	for _, m := range methods {
		msg := r.Find(m.Input())
		if msg == nil {
			resolveErr = errors.New("input " + m.Input() + " is not defined in " + m.Package())
			return
		}
		if !cache.Exists("m_" + msg.Name()) {
			if m.Mutation != nil {
				inputs = append(inputs, msg)
			} else {
				types = append(types, msg)
			}
			cache.Add("m_" + msg.Name())
		}
		if !cache.Exists("p_"+msg.GoPackage()) && !spec.IsGooglePackage(msg) {
			packages = append(packages, spec.NewPackage(msg.GoPackage()))
			cache.Add("p_" + msg.GoPackage())
		}
		ts, es, ps, err := r.resolveRecursive(msg.Fields(), cache)
		if err != nil {
			resolveErr = err
		}
		types = append(types, ts...)
		enums = append(enums, es...)
		packages = append(packages, ps...)

		msg = r.Find(m.Output())
		if msg == nil {
			resolveErr = errors.New("output " + m.Output() + " is not defined in " + m.Package())
			return
		}
		if m.ExposeQuery() == "" {
			if !cache.Exists("m_" + msg.Name()) {
				types = append(types, msg)
				cache.Add("m_" + msg.Name())
			}
			if !cache.Exists("p_"+msg.GoPackage()) && !spec.IsGooglePackage(msg) {
				packages = append(packages, spec.NewPackage(msg.GoPackage()))
				cache.Add("p_" + msg.GoPackage())
			}
		}
		ts, es, ps, err = r.resolveRecursive(m.ExposeQueryFields(msg), cache)
		if err != nil {
			resolveErr = err
			return
		}
		types = append(types, ts...)
		enums = append(enums, es...)
		packages = append(packages, ps...)

		if m.ExposeMutation() == "" {
			if !cache.Exists("m_" + msg.Name()) {
				types = append(types, msg)
				cache.Add("m_" + msg.Name())
			}
			if !cache.Exists("p_" + msg.GoPackage()) {
				packages = append(packages, spec.NewPackage(msg.GoPackage()))
				cache.Add("p_" + msg.GoPackage())
			}
		}
		ts, es, ps, err = r.resolveRecursive(m.ExposeQueryFields(msg), cache)
		if err != nil {
			resolveErr = err
			return
		}
		types = append(types, ts...)
		enums = append(enums, es...)
		packages = append(packages, ps...)
	}

	return
}

// resolveRecursive resolves all types in fields recursively.
func (r *Resolver) resolveRecursive(
	fields []*spec.Field,
	c *Cache,
) (
	types []*spec.Message,
	enums []*spec.Enum,
	packages []*spec.Package,
	resolveErr error,
) {
	for _, f := range fields {
		switch f.Type() {
		case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
			mm, ok := r.messages[f.TypeName()]
			if !ok {
				resolveErr = errors.New("failed to resolve message: " + f.TypeName())
			}
			f.TypeMessage = mm
			if !c.Exists("m_" + mm.Name()) {
				types = append(types, mm)
				c.Add("m_" + mm.Name())
			}
			if !c.Exists("p_"+mm.GoPackage()) && !spec.IsGooglePackage(mm) {
				packages = append(packages, spec.NewPackage(mm.GoPackage()))
				c.Add("p_" + mm.GoPackage())
			}
			ts, es, ps, err := r.resolveRecursive(mm.Fields(), c)
			if err != nil {
				resolveErr = err
			}
			types = append(types, ts...)
			enums = append(enums, es...)
			packages = append(packages, ps...)
		case descriptor.FieldDescriptorProto_TYPE_ENUM:
			en, ok := r.enums[f.TypeName()]
			if !ok {
				resolveErr = errors.New("failed to resolve enum: " + f.TypeName())
				return
			}
			f.TypeEnum = en
			if !c.Exists("e_" + en.Name()) {
				enums = append(enums, en)
				c.Add("e_" + en.Name())
			}
			if !c.Exists("p_"+en.GoPackage()) && !spec.IsGooglePackage(en) {
				packages = append(packages, spec.NewPackage(en.GoPackage()))
				c.Add("p_" + en.GoPackage())
			}
		}
	}

	return
}
