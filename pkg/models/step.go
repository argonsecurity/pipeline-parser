package models

const (
	CommitSHA     VersionType = "commit"
	TagVersion    VersionType = "tag"
	BranchVersion VersionType = "branch"
	Latest        VersionType = "latest"
	None          VersionType = "none"

	ShellStepType StepType = "shell"
	TaskStepType  StepType = "task"
)

type VersionType string
type StepType string

type Shell struct {
	Type   *string
	Script *string
}

type Step struct {
	ID                   *string
	Name                 *string
	Type                 StepType
	FailsPipeline        *bool
	Disabled             *bool
	EnvironmentVariables *EnvironmentVariables
	WorkingDirectory     *string
	Timeout              *int
	Conditions           *[]Condition
	Shell                *Shell
	Task                 *Task
}

type Task struct {
	ID          *string
	Name        *string
	Inputs      *[]Parameter
	Version     *string
	VersionType VersionType
}
