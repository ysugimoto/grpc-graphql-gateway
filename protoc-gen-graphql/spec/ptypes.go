package spec

import (
	"fmt"
	"strings"

	"path/filepath"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// Before protoc v3.14.0, go_package option name does not have "pb" suffix.
var supportedPtypes = []string{
	"timestamp",
	"wrappers",
	"empty",
}

// After protoc v3.14.0, go_package option name have been changed.
// @see https://github.com/protocolbuffers/protobuf/releases/tag/v3.14.0
var supportedPtypesLaterV3_14_0 = []string{
	"timestamppb",
	"wrapperspb",
	"emptypb",
}

func getSupportedPtypeNames(cv *plugin.Version) []string {
	if cv.GetMajor() >= 3 && cv.GetMinor() >= 14 {
		return supportedPtypesLaterV3_14_0
	}
	return supportedPtypes
}

func getImplementedPtypes(m *Message) (string, error) {
	ptype := strings.ToLower(filepath.Base(m.GoPackage()))

	var found bool
	for _, v := range getSupportedPtypeNames(m.CompilerVersion) {
		if ptype == v {
			found = true
		}
	}
	if !found {
		return "", fmt.Errorf("google's ptype \"%s\" does not implement for now.", ptype)
	}

	return ptype, nil
}
