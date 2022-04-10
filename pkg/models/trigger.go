package models

const (
	PushEvent            EventType = "push"
	PullRequestEvent     EventType = "pull_request"
	ForkEvent            EventType = "fork"
	ManualEvent          EventType = "manual"
	PipelineTriggerEvent EventType = "pipeline_trigger"
	PipelineRunEvent     EventType = "pipeline_run"
)

type EventType string

type Trigger struct {
	Branches  *Filter
	Paths     *Filter
	Paramters []Parameter
	Pipelines []string
	Filters   map[string]any
	Event     EventType
	Disabled  *bool
	Scheduled *string
}