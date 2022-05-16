package enhancers

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers/config"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

var (
	platformToEnhancerMapping = map[consts.Platform]*config.EnhancementConfiguration{}
)

func Enhance(pipeline *models.Pipeline, platform consts.Platform) (*models.Pipeline, error) {
	platformConfig := platformToEnhancerMapping[platform]

	if pipeline.Jobs != nil {
		jobs := make([]*models.Job, len(pipeline.Jobs))
		for i, job := range pipeline.Jobs {
			if job.Steps != nil {
				steps := make([]*models.Step, len(job.Steps))
				for i, step := range job.Steps {
					step = enhanceStep(step, config.CommonConfiguration)
					if platformConfig != nil {
						step = enhanceStep(step, platformConfig)
					}
					steps[i] = step
				}
				job.Steps = steps
			}
			job = enhanceJob(job, config.CommonConfiguration)
			if platformConfig != nil {
				job = enhanceJob(job, platformConfig)
			}
			jobs[i] = job
		}

		pipeline.Jobs = jobs
	}
	return pipeline, nil
}
