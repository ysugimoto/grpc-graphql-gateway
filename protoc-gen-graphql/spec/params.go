package spec

import (
	"errors"
	"strings"
)

// Params spec have plugin parameters
type Params struct {
	QueryOut string
}

func NewParams(p string) (*Params, error) {
	params := &Params{}

	for _, v := range strings.Split(p, ",") {
		kv := strings.SplitN(v, "=", 2)
		if len(kv) == 1 {
			return nil, errors.New("argument " + kv[0] + " must have value")
		}
		switch kv[0] {
		case "query":
			params.QueryOut = kv[1]
		default:
			return nil, errors.New("Unacceptable argument " + kv[0] + " provided")
		}
	}
	return params, nil
}
