package enhancers

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers/config"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

var (
	platformToEnhancerMapping = map[consts.Platform]config.EnhancementConfiguration{
		consts.GitHubPlatform: config.GithubConfiguration,
	}
)

func Enhance(pipeline *models.Pipeline, platform consts.Platform) (*models.Pipeline, error) {
	platformConfig, ok := platformToEnhancerMapping[platform]
	if !ok {
		return pipeline, &consts.ErrInvalidPlatform{Platform: platform}
	}

	if pipeline.Jobs != nil {
		jobs := make([]models.Job, len(*pipeline.Jobs))
		for i, job := range *pipeline.Jobs {
			job = enhanceJob(job, config.CommonConfiguration)
			job = enhanceJob(job, platformConfig)
			if job.Steps != nil {
				steps := make([]models.Step, len(*job.Steps))
				for i, step := range *job.Steps {
					step = enhanceStep(step, config.CommonConfiguration)
					step = enhanceStep(step, platformConfig)
					steps[i] = step
				}
				job.Steps = &steps
			}
			jobs[i] = job
		}

		pipeline.Jobs = &jobs
	}
	return pipeline, nil
}
