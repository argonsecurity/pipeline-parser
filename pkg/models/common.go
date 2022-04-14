package models

type EnvironmentVariables map[string]any
type Condition string

type Filter struct {
	AllowList []string
	DenyList  []string
}

type Variable struct {
	Context   *string
	Name      *string
	Value     *string
	Parent    *Entity
	Masked    *bool
	Protected *bool
}

type Parameter struct {
	Name        *string
	Value       *string
	Description *string
	Default     any
}

type Entity struct {
	ID   *string
	Name *string
}

type Runner struct {
	Type           *string
	Labels         *[]string
	OS             *string
	Arch           *string
	DockerMetadata *DockerMetadata
}
