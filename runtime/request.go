package runtime

import (
	"errors"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/iancoleman/strcase"
)

type GraphqlRequest struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                 `json:"operationName"`
}

// ParseRequest parses graphql query and variables from each request methods
func parseRequest(r *http.Request) (*GraphqlRequest, error) {
	var body []byte

	// Get request body
	switch r.Method {
	case http.MethodPost:
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, errors.New("malformed request body, " + err.Error())
		}
		body = buf
	case http.MethodGet:
		body = []byte(r.URL.Query().Get("query"))
	default:
		return nil, errors.New("invalid request method: '" + r.Method + "'")
	}

	// And try to parse
	var req GraphqlRequest
	if err := json.Unmarshal(body, &req); err != nil {
		// If error, the request body may come with single query line
		req.Query = string(body)
	}
	return &req, nil
}

func MarshalRequest(args interface{}, v interface{}, isCamel bool) error {
	if args == nil {
		return errors.New("Resolved params should be non-nil")
	}
	m, ok := args.(map[string]interface{}) // graphql.ResolveParams or nested object
	if !ok {
		return errors.New("Failed to type conversion of map[string]interface{}")
	}
	if isCamel {
		m = toLowerCaseKeys(m)
	}
	buf, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, &v)
}

func toLowerCaseKeys(args map[string]interface{}) map[string]interface{} {
	lc := make(map[string]interface{})
	for k, v := range args {
		lc[strcase.ToSnake(k)] = marshal(v)
	}
	return lc
}

func marshal(v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		return toLowerCaseKeys(t)
	case []interface{}:
		ret := make([]interface{}, len(t))
		for i, si := range t {
			ret[i] = marshal(si)
		}
		return ret
	default:
		return t
	}
}
