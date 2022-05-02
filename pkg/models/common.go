package models

type EnvironmentVariables map[string]any
type EnvironmentVariablesRef struct {
	EnvironmentVariables
	FileLocation *FileLocation
}

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
	Value       any
	Description *string
	Default     any
}

type Entity struct {
	ID   *string
	Name *string
}

type FileRef struct {
	Line   int
	Column int
}

type FileLocation struct {
	StartRef *FileRef
	EndRef   *FileRef
}
