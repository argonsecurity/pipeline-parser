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
	if pipeline.Image != nil {
		runner := &models.Runner{
			DockerMetadata: &models.DockerMetadata{
				Image: &pipeline.Image.Name,
			},
		}
		if pipeline.Image.ImageWithCustomUser != nil {
			runner.DockerMetadata = &models.DockerMetadata{
				Image: pipeline.Image.ImageWithCustomUser.Name,
			}
		}

		return runner
	}
	return nil
}
