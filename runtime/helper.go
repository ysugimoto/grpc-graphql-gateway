package runtime

import (
	"errors"
	"fmt"
	"net/http"

	"io/ioutil"
)

// Response responds HTTP response easily
func Respond(w http.ResponseWriter, status int, message string) {
	m := []byte(message)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(m)))
	w.WriteHeader(status)
	if len(m) > 0 {
		w.Write(m)
	}
}

// parseRequest parses graphql query and variables from each request methods
func parseRequest(r *http.Request) (
	query string,
	variables map[string]interface{},
	parseError error,
) {
	switch r.Method {
	case http.MethodPost:
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			parseError = errors.New("malformed request body, " + err.Error())
			return
		}
		query = string(buf)
	case http.MethodGet:
		query = r.URL.Query().Get("query")
	default:
		parseError = errors.New("invalid request method: '" + r.Method + "'")
	}
	return
}
