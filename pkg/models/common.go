package models

type Platform string

type EnvironmentVariables map[string]any
type EnvironmentVariablesRef struct {
	EnvironmentVariables `json:"environment_variables,omitempty"`
	FileReference        *FileReference `json:"file_reference,omitempty"`
	Imports              *Import        `json:"imports,omitempty"`
}

type Condition struct {
	Statement string            `json:"statement,omitempty"`
	Allow     *bool             `json:"allow,omitempty"`
	Paths     *Filter           `json:"paths,omitempty"`
	Exists    *Filter           `json:"exists,omitempty"`
	Branches  *Filter           `json:"branches,omitempty"`
	Events    []EventType       `json:"events,omitempty"`
	Variables map[string]string `json:"variables,omitempty"`
}

type Filter struct {
	AllowList []string `json:"allow_list,omitempty"`
	DenyList  []string `json:"deny_list,omitempty"`
}

type Variable struct {
	Context   *string `json:"context,omitempty"`
	Name      *string `json:"name,omitempty"`
	Value     *string `json:"value,omitempty"`
	Parent    *Entity `json:"parent,omitempty"`
	Masked    *bool   `json:"masked,omitempty"`
	Protected *bool   `json:"protected,omitempty"`
}

type Parameter struct {
	Name          *string        `json:"name,omitempty"`
	Value         any            `json:"value,omitempty"`
	Description   *string        `json:"description,omitempty"`
	Default       any            `json:"default,omitempty"`
	Options       []string       `json:"options,omitempty"`
	FileReference *FileReference `json:"file_reference,omitempty"`
}

type Entity struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type FileLocation struct {
	Line   int `json:"line,omitempty"`
	Column int `json:"column,omitempty"`
}

type FileReference struct {
	StartRef *FileLocation `json:"start_ref,omitempty"`
	EndRef   *FileLocation `json:"end_ref,omitempty"`
	IsAlias  bool          `json:"is_alias,omitempty"`
}
