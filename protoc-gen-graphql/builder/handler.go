package builder

var handlerTemplate = `

func graphqlHandler(endpoint string, v interface{}) (runtime.GraphqlHandler, error) {
	var c *runtime.Connection
	if v == nil {
		c = runtime.NewConnection(nil)
	} else {
		switch t := v.(type) {
		case *grpc.ClientConn:
			c = runtime.NewConnection(t)
		case *runtime.Connection:
			c = t
		default:
			return nil, errors.New("invalid type conversion")
		}
	}

	schema := createSchema(c)

	return func(w http.ResponseWriter, r *http.Request) *graphql.Result {
		if r.URL.Path != endpoint {
			runtime.Respond(w, http.StatusNotFound, "endpoint not found")
			return nil
		}
		query, variables, err := runtime.ParseRequest(r)
		if err != nil {
			runtime.Respond(w, http.StatusBadRequest, err.Error())
			return nil
		}

		return graphql.Do(graphql.Params{
			Schema: schema,
			RequestString: query,
			VariableValues: variables,
			Context: r.Context(),
		})
	}, nil
}

func RegisterGraphqlHandler(mux *runtime.ServeMux, v interface{}, endpoint string) (err error) {
	mux.Handler, err = graphqlHandler(endpoint, v)
	return
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
