package spec

import (
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
)

const (
	mainPackage = "main"
)

type PackageGetter interface {
	Package() string
	GoPackage() string
	Filename() string
}

type Package struct {
	Name                    string
	CamelName               string
	Path                    string
	GeneratedFilenamePrefix string
}

func NewPackage(g PackageGetter) *Package {
	p := &Package{}
	p.Name = strings.TrimSuffix(filepath.Base(g.Filename()), filepath.Ext(g.Filename()))
	p.GeneratedFilenamePrefix = strings.TrimSuffix(g.Filename(), filepath.Ext(g.Filename()))

	if pkg := g.GoPackage(); pkg != "" {
		// Support custom package definitions like example.com/path/to/package:packageName
		if index := strings.Index(pkg, ";"); index > -1 {
			p.Path = pkg[0:index]
		} else {
			p.Path = pkg
		}
	}

	p.CamelName = strcase.ToCamel(p.Name)
	return p
}

func NewGooglePackage(m PackageGetter) *Package {
	name := filepath.Base(m.GoPackage())

	return &Package{
		Name:      "gql_ptypes_" + strings.ToLower(name),
		CamelName: strcase.ToCamel(name),
		Path:      "github.com/ysugimoto/grpc-graphql-gateway/ptypes/" + strings.ToLower(name),
	}
}

func NewGoPackageFromString(pkg string) *Package {
	p := &Package{}
	// Support custom package definitions like example.com/path/to/package:packageName
	if index := strings.Index(pkg, ";"); index > -1 {
		p.Name = pkg[index+1:]
		p.Path = pkg[0:index]
	} else {
		p.Name = filepath.Base(pkg)
		p.Path = pkg
	}
	return p
}
