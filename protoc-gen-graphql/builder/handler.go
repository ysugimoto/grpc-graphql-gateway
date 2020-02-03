package builder

var handlerTemplate = `

func graphqlHandler(endpoint string, conn *grpc.ClientConn) runtime.GraphqlHandler {
	schema := createSchema(conn)

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
	}
}

func RegisterGraphqlHandler(mux *runtime.ServeMux, conn *grpc.ClientConn, endpoint string) {
	mux.Handler = graphqlHandler(endpoint, conn)
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
