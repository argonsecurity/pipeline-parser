package models

const (
	PushEvent            EventType = "push"
	PullRequestEvent     EventType = "pull_request"
	ForkEvent            EventType = "fork"
	ManualEvent          EventType = "manual"
	PipelineTriggerEvent EventType = "pipeline_trigger"
	PipelineRunEvent     EventType = "pipeline_run"
	ScheduledEvent       EventType = "scheduled"
)

type EventType string

type Trigger struct {
	Branches      *Filter
	Paths         *Filter
	Exists        *Filter
	Parameters    []Parameter
	Pipelines     []string
	Filters       map[string]any
	Event         EventType
	Disabled      *bool
	Schedules     *[]string
	FileReference *FileReference
}

type Triggers struct {
	Triggers      []*Trigger
	FileReference *FileReference
}
