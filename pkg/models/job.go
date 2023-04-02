package models

type TokenPermissions struct {
	Permissions   map[string]Permission
	FileReference *FileReference
}

type ConcurrencyGroup string

type SecretsRef struct {
	Secrets       map[string]any `json:"secrets,omitempty"`
	Inherit       bool           `json:"inherit,omitempty"`
	FileReference *FileReference `json:"file_reference,omitempty"`
}

type Job struct {
	ID                   *string                  `json:"id,omitempty"`
	Name                 *string                  `json:"name,omitempty"`
	Steps                []*Step                  `json:"steps,omitempty"`
	ContinueOnError      *string                  `json:"continue_on_error,omitempty"`
	PreSteps             []*Step                  `json:"pre_steps,omitempty"`
	PostSteps            []*Step                  `json:"post_steps,omitempty"`
	EnvironmentVariables *EnvironmentVariablesRef `json:"environment_variables,omitempty"`
	Runner               *Runner                  `json:"runner,omitempty"`
	Conditions           []*Condition             `json:"conditions,omitempty"`
	ConcurrencyGroup     *ConcurrencyGroup        `json:"concurrency_group,omitempty"`
	Inputs               []*Parameter             `json:"inputs,omitempty"`
	TimeoutMS            *int                     `json:"timeout_ms,omitempty"`
	Tags                 []string                 `json:"tags,omitempty"`
	TokenPermissions     *TokenPermissions        `json:"token_permissions,omitempty"`
	Dependencies         []*JobDependency         `json:"dependencies,omitempty"`
	Metadata             Metadata                 `json:"metadata,omitempty"`
	Matrix               *Matrix                  `json:"matrix,omitempty"`
	FileReference        *FileReference           `json:"file_reference,omitempty"`
	Imports              *Import                  `json:"imports,omitempty"`
}

type Matrix struct {
	Matrix        map[string]any
	Include       []map[string]any
	Exclude       []map[string]any
	FileReference *FileReference
}

type JobDependency struct {
	JobID            *string           `json:"job_id,omitempty"`
	ConcurrencyGroup *ConcurrencyGroup `json:"concurrency_group,omitempty"`
	Pipeline         *string           `json:"pipeline,omitempty"`
}
