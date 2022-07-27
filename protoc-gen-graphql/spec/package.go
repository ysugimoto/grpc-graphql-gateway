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
	FileName                string
}

func NewPackage(g PackageGetter) *Package {
	p := &Package{}
	p.GeneratedFilenamePrefix = strings.TrimSuffix(g.Filename(), filepath.Ext(g.Filename()))
	p.FileName = filepath.Base(p.GeneratedFilenamePrefix)

	if pkg := g.GoPackage(); pkg != "" {
		// Support custom package definitions like example.com/path/to/package:packageName
		if index := strings.Index(pkg, ";"); index > -1 {
			p.Name = pkg[index+1:]
			p.Path = pkg[0:index]
		} else {
			p.Name = filepath.Base(pkg)
			p.Path = pkg
		}
	} else if pkg := g.Package(); pkg != "" {
		p.Name = pkg
	} else {
		p.Name = strings.ReplaceAll(
			g.Filename(),
			filepath.Ext(g.Filename()),
			"",
		)
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
