package spec

type DependType int

const (
	DependTypeMessage DependType = iota
	DependTypeInput
	DependTypeEnum
	DependTypeInterface
)

// shorthand alias
type ms map[string]struct{}

type Dependencies struct {
	message    ms
	enum       ms
	input      ms
	interfaces ms
}

func NewDependencies() *Dependencies {
	return &Dependencies{
		message:    ms{},
		enum:       ms{},
		input:      ms{},
		interfaces: ms{},
	}
}

func (d *Dependencies) Depend(t DependType, pkg string) {
	switch t {
	case DependTypeMessage:
		d.message[pkg] = struct{}{}
	case DependTypeEnum:
		d.enum[pkg] = struct{}{}
	case DependTypeInput:
		d.input[pkg] = struct{}{}
	case DependTypeInterface:
		d.interfaces[pkg] = struct{}{}
	}
}

func (d *Dependencies) IsDepended(t DependType, pkg string) bool {
	var ok bool
	switch t {
	case DependTypeMessage:
		_, ok = d.message[pkg]
	case DependTypeEnum:
		_, ok = d.enum[pkg]
	case DependTypeInput:
		_, ok = d.input[pkg]
	case DependTypeInterface:
		_, ok = d.interfaces[pkg]
	}
	return ok
}

func (d *Dependencies) GetDependendencies() map[string][]string {
	ret := map[string][]string{
		"message":   {},
		"enum":      {},
		"input":     {},
		"interface": {},
	}
	for p := range d.message {
		ret["message"] = append(ret["message"], p)
	}
	for p := range d.enum {
		ret["enum"] = append(ret["enum"], p)
	}
	for p := range d.input {
		ret["input"] = append(ret["input"], p)
	}
	for p := range d.interfaces {
		ret["interface"] = append(ret["interface"], p)
	}
	return ret
}
