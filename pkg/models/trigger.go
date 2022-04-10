package models

const (
	PushEvent        EventType = "push"
	PullRequestEvent EventType = "pull_request"
	ForkEvent        EventType = "fork"
	ManualEvent      EventType = "manual"
)

type EventType string

type Trigger struct {
	Branches  *Filter
	Paths     *Filter
	Paramters []Parameter
	Filters   map[string]string
	Event     EventType
	Disabled  *bool
	Scheduled *string
}
