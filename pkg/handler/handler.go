package handler

import (
	"fmt"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	generalEnhancer "github.com/argonsecurity/pipeline-parser/pkg/enhancers/general"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders"
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	bitbucketModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/parsers"
)

type Handler[T any] interface {
	GetPlatform() models.Platform
	GetLoader() loaders.Loader[T]
	GetParser() parsers.Parser[T]
	GetEnhancer() enhancers.Enhancer
}

func Handle(data []byte, platform models.Platform, credentials *models.Credentials) (*models.Pipeline, error) {
	var pipeline *models.Pipeline
	var err error

	if len(data) == 0 {
		return nil, consts.NewErrEmptyData()
	}

	switch platform {
	case consts.GitHubPlatform:
		pipeline, err = handle[githubModels.Workflow](data, &GitHubHandler{}, credentials)
	case consts.GitLabPlatform:
		pipeline, err = handle[gitlabModels.GitlabCIConfiguration](data, &GitLabHandler{}, credentials)
	case consts.AzurePlatform:
		pipeline, err = handle[azureModels.Pipeline](data, &AzureHandler{}, credentials)
	case consts.BitbucketPlatform:
		pipeline, err = handle[bitbucketModels.Pipeline](data, &BitbucketHandler{}, credentials)
	default:
		return nil, consts.NewErrInvalidPlatform(platform)
	}

	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

func handle[T any](data []byte, handler Handler[T], credentials *models.Credentials) (*models.Pipeline, error) {
	pipeline, err := handler.GetLoader().Load(data)
	if err != nil {
		return nil, err
	}

	parsedPipeline, err := handler.GetParser().Parse(pipeline)
	if err != nil {
		return nil, err
	}

	enhancer := handler.GetEnhancer()

	importedPipelines, err := enhancer.LoadImportedPipelines(parsedPipeline, credentials)
	if err != nil {
		fmt.Printf("Failed getting imported pipelines:\n%v", err)
	}

	for _, importedPipeline := range importedPipelines {
		parsedImportedPipeline, err := handle(importedPipeline.Data, handler, credentials)
		if err != nil {
			fmt.Printf("Failed parsing imported pipeline for job %s - %v", importedPipeline.JobName, err)
		}
		importedPipeline.Pipeline = parsedImportedPipeline
	}

	enhancedPipeline, err := handler.GetEnhancer().Enhance(parsedPipeline, importedPipelines)
	if err != nil {
		fmt.Printf("Error while enhancing pipeline:\n%v", err)
	}

	return generalEnhancer.Enhance(enhancedPipeline, handler.GetPlatform())
}
