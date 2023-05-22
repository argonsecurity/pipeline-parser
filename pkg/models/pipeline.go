package models

type Pipeline struct {
	Id         *string      `json:"id,omitempty"`
	Name       *string      `json:"name,omitempty"`
	Triggers   *Triggers    `json:"triggers,omitempty"`
	Jobs       []*Job       `json:"jobs,omitempty"`
	Imports    []*Import    `json:"imports,omitempty"`
	Parameters []*Parameter `json:"parameters,omitempty"`
	Defaults   *Defaults    `json:"defaults,omitempty"`
	Platform   Platform     `json:"platform,omitempty"`
}

type Scans struct {
	Secrets      *bool `json:"secrets,omitempty"`
	Iac          *bool `json:"iac,omitempty"`
	Pipelines    *bool `json:"pipelines,omitempty"`
	SAST         *bool `json:"sast,omitempty"`
	Dependencies *bool `json:"dependencies,omitempty"`
	License      *bool `json:"license,omitempty"`
}

type Resources struct {
	Repositories  []*ImportSource `json:"repositories,omitempty"`
	FileReference *FileReference  `json:"file_reference,omitempty"`
}

type Defaults struct {
	EnvironmentVariables *EnvironmentVariablesRef `json:"environment_variables,omitempty"`
	Scans                *Scans                   `json:"scans,omitempty"`
	Runner               *Runner                  `json:"runner,omitempty"`
	Conditions           []*Condition             `json:"conditions,omitempty"`
	ContinueOnError      *bool                    `json:"continue_on_error,omitempty"`
	TokenPermissions     *TokenPermissions        `json:"token_permissions,omitempty"`
	Settings             *map[string]any          `json:"settings,omitempty"`
	FileReference        *FileReference           `json:"file_reference,omitempty"`
	PostSteps            []*Step                  `json:"post_steps,omitempty"`
	PreSteps             []*Step                  `json:"pre_steps,omitempty"`
	Resources            *Resources               `json:"resources,omitempty"`
}
