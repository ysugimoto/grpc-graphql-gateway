package spec

import (
	"path/filepath"
)

type Package struct {
	Name string
	Path string
}

func NewPackage(p string) *Package {
	return &Package{
		Name: filepath.Base(p),
		Path: p,
	}
}
