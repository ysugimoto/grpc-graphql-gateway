package spec

type DependType int

const (
	DependTypeMessage DependType = iota
	DependTypeInput
	DependTypeEnum
)

// shorthand alias
type ms map[string]struct{}

type Dependencies struct {
	message ms
	enum    ms
	input   ms
}

func NewDependencies() *Dependencies {
	return &Dependencies{
		message: ms{},
		enum:    ms{},
		input:   ms{},
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
	}
	return ok
}

func (d *Dependencies) GetDependendencies() map[string][]string {
	ret := map[string][]string{
		"message": []string{},
		"enum":    []string{},
		"input":   []string{},
	}
	for p, _ := range d.message {
		ret["message"] = append(ret["message"], p)
	}
	for p, _ := range d.enum {
		ret["enum"] = append(ret["enum"], p)
	}
	for p, _ := range d.input {
		ret["input"] = append(ret["input"], p)
	}
	return ret
}
