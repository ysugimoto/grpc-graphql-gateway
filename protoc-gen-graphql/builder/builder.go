package builder

// Builder is interface for collect sections to generate code/query.
// BuildQuery() is used for GraphQL schema definition,
// BuildProgram is used for Go pgoram generation.
// Every builder should implement these methods to genearte appropriately.
type Builder interface {
	BuildQuery() (string, error)
	BuildProgram() (string, error)
}
