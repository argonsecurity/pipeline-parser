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
	Branches      *Filter        `json:"branches,omitempty"`
	Paths         *Filter        `json:"paths,omitempty"`
	Tags          *Filter        `json:"tags,omitempty"`
	Exists        *Filter        `json:"exists,omitempty"`
	Parameters    []Parameter    `json:"parameters,omitempty"`
	Pipelines     []string       `json:"pipelines,omitempty"`
	Filters       map[string]any `json:"filters,omitempty"`
	Event         EventType      `json:"event,omitempty"`
	Disabled      *bool          `json:"disabled,omitempty"`
	Schedules     *[]string      `json:"schedules,omitempty"`
	FileReference *FileReference `json:"file_reference,omitempty"`
}

type Triggers struct {
	Triggers      []*Trigger     `json:"triggers,omitempty"`
	FileReference *FileReference `json:"file_reference,omitempty"`
}
