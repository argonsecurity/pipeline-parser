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

type Trigger struct {
	Branches    *Filter
	Paths       *Filter
	PullRequest *bool
	Manual      *bool
	Disabled    *bool
	Push        *bool
	Scheduled   *string
	Events      *[]string
}

type Filter struct {
	AllowList *[]string
	DenyList  *[]string
}

type EnvironmentVariables map[string]string
type Variable struct {
	Context   *string
	Name      *string
	Value     *string
	Parent    *Entity
	Masked    *bool
	Protected *bool
}

type Condition string
type Parameter struct {
	Name        *string
	Value       *string
	Description *string
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
	Timeout              *int
	Tags                 *[]string
}

type Step struct {
	ID                   *string
	Name                 *string
	Type                 *string
	FailsPipeline        *bool
	Disabled             *bool
	EnvironmentVariables *EnvironmentVariables
	WorkingDirectory     *string
	Timeout              *int
	Conditions           *[]Condition
}

type Repository struct {
	Entity
	URL          *string
	Organization *Entity
}

type Entity struct {
	ID   *string
	Name *string
}

type DockerMetadata struct {
	Image                 *string
	Label                 *string
	RegistryURL           *string
	RegistryCredentialsID *string
}
type Runner struct {
	Type           *string
	Labels         *[]string
	OS             *string
	Arch           *string
	DockerMetadata *DockerMetadata
}
