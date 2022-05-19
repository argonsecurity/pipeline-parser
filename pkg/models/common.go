package models

type EnvironmentVariables map[string]any
type EnvironmentVariablesRef struct {
	EnvironmentVariables
	FileReference *FileReference
}

type Condition struct {
	Statement string
	Allow     *bool
	Paths     *Filter
	Exists    *Filter
	Branches  *Filter
	Events    []EventType
	Variables map[string]string
}

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

type FileLocation struct {
	Line   int
	Column int
}

type FileReference struct {
	StartRef *FileLocation
	EndRef   *FileLocation
}
