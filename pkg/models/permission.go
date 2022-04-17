package models

const (
	PullRequestPermission = "pull-request"
	PushPermission        = "push"
	RunPipelinePermission = "run-pipeline"
)

type Permission struct {
	Read  bool
	Write bool
	Admin bool
}
