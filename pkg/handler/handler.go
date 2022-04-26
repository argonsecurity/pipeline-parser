package handler

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
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

	switch platform {
	case consts.GitHubPlatform:
		pipeline, err = handle[githubModels.Workflow](data, &GitHubHandler{})
	}

	if err != nil {
		return nil, err
	}

	return enhancers.Enhance(pipeline, platform)
}

func handle[T any](data []byte, handler Handler[T]) (*models.Pipeline, error) {
	workflow, err := handler.GetLoader().Load(data)
	if err != nil {
		return nil, err
	}

	pipeline, err := handler.GetParser().Parse(workflow)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}
