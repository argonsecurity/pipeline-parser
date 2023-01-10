package handler

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders"
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	bitbucketModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/parsers"
)

type Handler[T any] interface {
	GetLoader() loaders.Loader[T]
	GetParser() parsers.Parser[T]
}

func Handle(data []byte, platform consts.Platform) (*models.Pipeline, error) {
	var pipeline *models.Pipeline
	var err error

	if len(data) == 0 {
		return nil, consts.NewErrEmptyData()
	}

	switch platform {
	case consts.GitHubPlatform:
		pipeline, err = handle[githubModels.Workflow](data, &GitHubHandler{})
	case consts.GitLabPlatform:
		pipeline, err = handle[gitlabModels.GitlabCIConfiguration](data, &GitLabHandler{})
	case consts.AzurePlatform:
		pipeline, err = handle[azureModels.Pipeline](data, &AzureHandler{})
	case consts.BitBucketPlatform:
		pipeline, err = handle[bitbucketModels.Pipeline](data, &BitBucketHandler{})
	default:
		return nil, consts.NewErrInvalidPlatform(platform)
	}

	if err != nil {
		return nil, err
	}

	return enhancers.Enhance(pipeline, platform)
}

func handle[T any](data []byte, handler Handler[T]) (*models.Pipeline, error) {
	pipeline, err := handler.GetLoader().Load(data)
	if err != nil {
		return nil, err
	}

	parsedPipeline, err := handler.GetParser().Parse(pipeline)
	if err != nil {
		return nil, err
	}

	return parsedPipeline, nil
}
