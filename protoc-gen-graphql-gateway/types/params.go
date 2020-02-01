package types

import (
	"errors"
	"strings"
)

type Params struct {
	QueryOut       string
	ProgramOut     string
	ProgramPackage string
}

func NewParams(p string) (*Params, error) {
	params := &Params{
		ProgramPackage: "main",
	}

	for _, v := range strings.Split(p, ",") {
		kv := strings.SplitN(v, "=", 2)
		if len(kv) == 1 {
			return nil, errors.New("argument " + kv[0] + " must have value")
		}
		switch kv[0] {
		case "query":
			params.QueryOut = kv[1]
		case "go":
			params.ProgramOut = kv[1]
		case "gopkg":
			params.ProgramPackage = kv[1]
		default:
			return nil, errors.New("Unacceptable argument " + kv[0] + " provided")
		}
	}
	return params, nil
}
