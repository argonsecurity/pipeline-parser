package models

type Pipeline struct {
	Id           *string
	Name         *string
	Triggers     *Triggers
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
	FileReference        *FileReference
}

type Repository struct {
	Entity
	URL          *string
	Organization *Entity
}
