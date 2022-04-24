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
	EnvironmentVariables *EnvironmentVariables
	Runner               *Runner
	Conditions           *[]Condition
	TokenPermissions     *map[string]Permission
	Settings             *map[string]any
}

type Repository struct {
	Entity
	URL          *string
	Organization *Entity
}
