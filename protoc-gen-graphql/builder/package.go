package builder

import (
	"fmt"
	"strings"
)

// Package builder generates package section for Go program.
type Package struct {
	pkgName string
}

func NewPackage(p string) *Package {
	return &Package{
		pkgName: p,
	}
}

func (b *Package) BuildQuery() (string, error) {
	return "", nil
}

func (b *Package) BuildProgram() (string, error) {
	return strings.TrimSpace(fmt.Sprintf(`
// This file is generated from proroc-gen-graphql, DO NOT EDIT!
package %s`,
		b.pkgName,
	)), nil
}
