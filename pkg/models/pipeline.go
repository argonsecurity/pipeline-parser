package models

type Pipeline struct {
	Id         *string
	Name       *string
	Triggers   *Triggers
	Jobs       []*Job
	Imports    []string
	Parameters []*Parameter
	Defaults   *Defaults
}

type Scans struct {
	Secrets      *bool
	Iac          *bool
	Pipelines    *bool
	SAST         *bool
	Dependencies *bool
	License      *bool
}

type Defaults struct {
	EnvironmentVariables *EnvironmentVariablesRef
	Scans                *Scans
	Runner               *Runner
	Conditions           []*Condition
	TokenPermissions     *TokenPermissions
	Settings             *map[string]any
	FileReference        *FileReference
	PostSteps            []*Step
	PreSteps             []*Step
}
