package builder

type Package struct {
	pkgName string
}

func NewPackage(p string) *Package {
	return &Package{
		pkgName: p,
	}
}

func (b *Package) BuildQuery() string {
	return ""
}

func (b *Package) BuildProgram() string {
	return "package " + b.pkgName + "\n"
}
