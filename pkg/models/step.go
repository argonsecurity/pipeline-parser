package models

const (
	CommitSHA     VersionType = "commit"
	TagVersion    VersionType = "tag"
	BranchVersion VersionType = "branch"
	Latest        VersionType = "latest"
	None          VersionType = "none"

	ShellType StepType = "shell"
	TaskType  StepType = "task"
)

type VersionType string
type StepType string

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
	Shell                *string
	Task                 *Task
}

type Task struct {
	ID          *string
	Name        *string
	Inputs      *[]Parameter
	Version     *string
	VersionType VersionType
}
