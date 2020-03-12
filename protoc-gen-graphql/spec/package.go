package spec

import (
	"path/filepath"

	"github.com/iancoleman/strcase"
)

const (
	mainPackage = "main"
)

type Package struct {
	Name      string
	CamelName string
	Path      string
}

func NewPackage(p string) *Package {
	return &Package{
		Name:      filepath.Base(p),
		CamelName: strcase.ToCamel(filepath.Base(p)),
		Path:      p,
	}
}
