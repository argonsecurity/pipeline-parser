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

type Job struct {
	ID                   *string
	Name                 *string
	Steps                *[]Step
	ContinueOnError      *bool
	PreSteps             *[]Step
	PostSteps            *[]Step
	EnvironmentVariables *EnvironmentVariables
	Runner               *Runner
	Conditions           *[]Condition
	ConcurrencyGroup     *string
	Inputs               *[]Parameter
	TimeoutMS            *int
	Tags                 *[]string
}

type Repository struct {
	Entity
	URL          *string
	Organization *Entity
}

type DockerMetadata struct {
	Image                 *string
	Label                 *string
	RegistryURL           *string
	RegistryCredentialsID *string
}
