package spec

import (
	"errors"
	"strings"
)

// Params spec have plugin parameters
type Params struct {
	QueryOut string
	Verbose  bool
}

func NewParams(p string) (*Params, error) {
	params := &Params{}
	if p == "" {
		return params, nil
	}

	for _, v := range strings.Split(p, ",") {
		kv := strings.SplitN(v, "=", 2)
		switch kv[0] {
		case "verbose":
			params.Verbose = true
		case "query":
			if len(kv) == 1 {
				return nil, errors.New("argument " + kv[0] + " must have value")
			}
			params.QueryOut = kv[1]
		default:
			return nil, errors.New("Unacceptable argument " + kv[0] + " provided")
		}
	}
	return params, nil
}
