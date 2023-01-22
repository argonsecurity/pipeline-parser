package bitbucket

import (
	bitbucketModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type BitbucketParser struct{}

func (g *BitbucketParser) Parse(bitbucketPipeline *bitbucketModels.Pipeline) (*models.Pipeline, error) {
	if bitbucketPipeline == nil {
		return nil, nil
	}

	var pipeline models.Pipeline

	pipeline.Defaults = parsePipelineDefaults(bitbucketPipeline)
	pipeline.Jobs = parseJobs(bitbucketPipeline)

	return &pipeline, nil
}
