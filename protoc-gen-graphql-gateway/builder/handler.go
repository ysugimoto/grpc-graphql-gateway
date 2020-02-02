package builder

var handlerTemplate = `

type ErrorHandler func(errs []gqlerrors.FormattedError)

const (
	optNameErrorHandler = "errorhandler"
	optNameAllowCORS = "allowcors"
)

type Option struct {
	name string
	value interface{}
}

func WithErrorHandler(eh ErrorHandler) Option {
	return Option {
		name: optNameErrorHandler,
		value: eh,
	}
}

func WithCORS() Option {
	return Option {
		name: optNameAllowCORS,
		value: true,
	}
}

type GraphqlResolver struct {
	schema graphql.Schema
	errorHandler ErrorHandler
	allowCORS bool
}

func New(opts ...Option) *GraphqlResolver {
	var eh ErrorHandler
	var cors bool

	for _, o := range opts {
		switch o.name {
			case optNameErrorHandler:
				eh = o.value.(ErrorHandler)
			case optNameAllowCORS:
				cors = true
		}
	}

	return &GraphqlResolver {
		errorHandler: eh,
		allowCORS: cors,
	}
}

func corsHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.URL.Host)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Max-Age", "1728000")
}

func respondError(w http.ResponseWriter, status int, message string) {
	m := []byte(message)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(m)))
	w.WriteHeader(status)
	if len(m) > 0 {
		w.Write(m)
	}
}

func (g *GraphqlResolver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if g.allowCORS {
		corsHeader(w, r)
	}
	var query string
	switch r.Method {
		case http.MethodOptions:
			respondError(w, http.StatusNoContent, "")
			return
		case http.MethodPost:
			buf, err := ioutil.ReadAll(r.Body)
			if err != nil {
				respondError(w, http.StatusBadRequest, "malformed request body")
				return
			}
			query = string(buf)
		case http.MethodGet:
			query = r.URL.Query().Get("query")
		default:
			respondError(w, http.StatusBadRequest, "invalid request method: '" + r.Method + "'")
			return
	}

	result := graphql.Do(graphql.Params{
		Schema: g.schema,
		RequestString: query,
		Context: r.Context(),
	})
	if len(result.Errors) > 0 {
		if g.errorHandler != nil {
			g.errorHandler(result.Errors)
		}
	}
	out, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprint(len(out)))
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}`

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (b *Handler) BuildQuery() string {
	return ""
}

func (b *Handler) BuildProgram() string {
	return handlerTemplate
}
