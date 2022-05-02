package models

type Pipeline struct {
	Id           *string
	Name         *string
	Triggers     *[]Trigger
	Jobs         *[]Job
	Imports      *[]string
	Parameters   *[]Parameter
	Repository   *Repository
	Organization *Entity
	Defaults     *Defaults
}

type Defaults struct {
	EnvironmentVariables *EnvironmentVariablesRef
	Runner               *Runner
	Conditions           *[]Condition
	TokenPermissions     *TokenPermissions
	Settings             *map[string]any
}

type Repository struct {
	Entity
	URL          *string
	Organization *Entity
}
