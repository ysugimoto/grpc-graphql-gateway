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
	Name      string
	CamelName string
	Path      string
}

func NewPackage(g PackageGetter) *Package {
	p := &Package{}
	if pkg := g.GoPackage(); pkg != "" {
		p.Name = filepath.Base(pkg)
		p.Path = pkg
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
