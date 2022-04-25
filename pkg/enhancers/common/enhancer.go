package common

import (
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type CommonEnhancer struct{}

func (c *CommonEnhancer) EnhanceJob(job models.Job) models.Job {
	return job
}

func (c *CommonEnhancer) EnhanceStep(step models.Step) models.Step {
	return step
}
