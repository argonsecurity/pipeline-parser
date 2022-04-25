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

	for label, labelConfig := range config.LabelMapping {
		if utils.AnyMatch(labelConfig.Names, *job.Name) {
			job.Metadata.Labels = append(job.Metadata.Labels, label)
		}
	}
	return job
}