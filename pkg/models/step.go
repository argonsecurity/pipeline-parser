package models

const (
	CommitSHA     VersionType = "commit"
	TagVersion    VersionType = "tag"
	BranchVersion VersionType = "branch"
	Latest        VersionType = "latest"
	None          VersionType = "none"

	ShellStepType StepType = "shell"
	TaskStepType  StepType = "task"

	DockerTaskType TaskType = "docker"
	CITaskType     TaskType = "ci"
)

type VersionType string
type StepType string
type TaskType string

type Shell struct {
	Type          *string        `json:"type,omitempty"`
	Script        *string        `json:"script,omitempty"`
	FileReference *FileReference `json:"file_reference,omitempty"`
}

type Step struct {
	ID                   *string                  `json:"id,omitempty"`
	Name                 *string                  `json:"name,omitempty"`
	Type                 StepType                 `json:"type,omitempty"`
	Runner               *Runner                  `json:"runner,omitempty"`
	FailsPipeline        *bool                    `json:"fails_pipeline,omitempty"`
	Disabled             *bool                    `json:"disabled,omitempty"`
	EnvironmentVariables *EnvironmentVariablesRef `json:"environment_variables,omitempty"`
	WorkingDirectory     *string                  `json:"working_directory,omitempty"`
	Timeout              *int                     `json:"timeout,omitempty"`
	Conditions           *[]Condition             `json:"conditions,omitempty"`
	Shell                *Shell                   `json:"shell,omitempty"`
	Task                 *Task                    `json:"task,omitempty"`
	Metadata             Metadata                 `json:"metadata,omitempty"`
	AfterScript          *Shell                   `json:"after_script,omitempty"`
	FileReference        *FileReference           `json:"file_reference,omitempty"`
	Imports              *Import                  `json:"imports,omitempty"`
}

type Task struct {
	ID          *string      `json:"id,omitempty"`
	Name        *string      `json:"name,omitempty"`
	Inputs      []*Parameter `json:"inputs,omitempty"`
	Version     *string      `json:"version,omitempty"`
	VersionType VersionType  `json:"version_type,omitempty"`
	Type        TaskType     `json:"type,omitempty"`
}
