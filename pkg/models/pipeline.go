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
	Defaults     *struct {
		EnvironmentVariables *EnvironmentVariables
		Runner               *Runner
		Conditions           *[]Condition
		TokenPermissions     *map[string]string
	}
}

type Repository struct {
	Entity
	URL          *string
	Organization *Entity
}
