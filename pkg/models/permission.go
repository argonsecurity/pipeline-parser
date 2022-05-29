package models

const (
	PullRequestPermission = "pull-request"
	PushPermission        = "push"
	RunPipelinePermission = "run-pipeline"
)

type Permission struct {
	Read  bool `json:"read,omitempty"`
	Write bool `json:"write,omitempty"`
	Admin bool `json:"admin,omitempty"`
}
