package models

type TokenPermissions struct {
	Permissions   map[string]Permission
	FileReference *FileReference
}

type Job struct {
	ID                   *string
	Name                 *string
	Steps                *[]Step
	ContinueOnError      *bool
	PreSteps             *[]Step
	PostSteps            *[]Step
	EnvironmentVariables *EnvironmentVariablesRef
	Runner               *Runner
	Conditions           *[]Condition
	ConcurrencyGroup     *string
	Inputs               *[]Parameter
	TimeoutMS            *int
	Tags                 *[]string
	TokenPermissions     *TokenPermissions
	Dependencies         *[]string
	Metadata             Metadata
	FileReference        *FileReference
}
