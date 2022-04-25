package models

import "github.com/argonsecurity/pipeline-parser/pkg/models"

type Enhancer interface {
	EnhanceJob(job models.Job) models.Job
	EnhanceStep(step models.Step) models.Step
}
