package enhancers

import (
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers/config"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func enhanceJob(job models.Job, config config.EnhancementConfiguration) models.Job {
	if utils.AnyMatch(config.Build.Names, *job.Name) {
		job.Metadata.Build = true
	}

	if utils.AnyMatch(config.Test.Names, *job.Name) {
		job.Metadata.Test = true
	}

	if utils.AnyMatch(config.Deploy.Names, *job.Name) {
		job.Metadata.Deploy = true
	}

	for _, step := range *job.Steps {
		if step.Metadata.Build {
			job.Metadata.Build = true
		}

		if step.Metadata.Test {
			job.Metadata.Test = true
		}

		if step.Metadata.Deploy {
			job.Metadata.Deploy = true
		}
	}

	return job
}
