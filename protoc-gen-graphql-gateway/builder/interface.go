package builder

type Builder interface {
	BuildQuery() string
	BuildProgram() string
}
