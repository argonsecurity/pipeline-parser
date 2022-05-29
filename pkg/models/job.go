package models

type TokenPermissions struct {
	Permissions   map[string]Permission
	FileReference *FileReference
}

type ConcurrencyGroup string

type Job struct {
	ID                   *string
	Name                 *string
	Steps                []*Step
	ContinueOnError      *bool
	PreSteps             []*Step
	PostSteps            []*Step
	EnvironmentVariables *EnvironmentVariablesRef
	Runner               *Runner
	Conditions           []*Condition
	ConcurrencyGroup     *ConcurrencyGroup
	Inputs               []*Parameter
	TimeoutMS            *int
	Tags                 []string
	TokenPermissions     *TokenPermissions
	Dependencies         []*JobDependency
	Metadata             Metadata
	FileReference        *FileReference
}

type JobDependency struct {
	JobID            *string
	ConcurrencyGroup *ConcurrencyGroup
	Pipeline         *string
}
