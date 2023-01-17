package bitbucket

import (
	bitbucketModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func parsePipelineDefaults(pipeline *bitbucketModels.Pipeline) *models.Defaults {
	if pipeline == nil {
		return nil
	}

	var defaults models.Defaults
	defaults.Runner = parseRunner(pipeline)

	if pipeline.Options != nil {
		defaults.Settings = &map[string]any{
			"docker":   pipeline.Options.Docker,
			"max-time": pipeline.Options.MaxTime,
			"size":     pipeline.Options.Size,
		}
	}

	return &defaults
}

func parseRunner(pipeline *bitbucketModels.Pipeline) *models.Runner {
	if pipeline == nil {
		return nil
	}

	if pipeline.Image != nil && pipeline.Image.ImageData != nil {
		runner := &models.Runner{
			DockerMetadata: &models.DockerMetadata{
				Image: pipeline.Image.ImageData.Name,
			},
		}

		return runner
	}
	return nil
}
