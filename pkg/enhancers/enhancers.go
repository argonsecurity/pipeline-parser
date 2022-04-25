package enhancers

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers/common"
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers/github"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type Enhancer interface {
	EnhanceJob(job models.Job) models.Job
	EnhanceStep(step models.Step) models.Step
}

var (
	commonEnhancer Enhancer = &common.CommonEnhancer{}

	platformToEnhancerMapping = map[consts.Platform]Enhancer{
		consts.GitHubPlatform: &github.GitHubEnhancer{},
	}
)

func Enhance(pipeline models.Pipeline, platform consts.Platform) (models.Pipeline, error) {
	enhancer, ok := platformToEnhancerMapping[platform]
	if !ok {
		return pipeline, &consts.ErrInvalidPlatform{Platform: platform}
	}

	jobs := make([]models.Job, len(*pipeline.Jobs))
	for _, job := range *pipeline.Jobs {
		job = enhancer.EnhanceJob(job)
		job = commonEnhancer.EnhanceJob(job)

		steps := make([]models.Step, len(*job.Steps))
		for _, step := range *job.Steps {
			step = enhancer.EnhanceStep(step)
			step = commonEnhancer.EnhanceStep(step)
			steps = append(steps, step)
		}
		job.Steps = &steps
		jobs = append(jobs, job)
	}

	pipeline.Jobs = &jobs

	return pipeline, nil
}
