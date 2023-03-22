package github

import (
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/pkg/errors"
)

type GitHubEnhancer struct{}

func (g *GitHubEnhancer) LoadImportedPipelines(data *models.Pipeline, credentials *models.Credentials) ([]*enhancers.ImportedPipeline, error) {
	var errs error
	importedPipelines, err := getReusableWorkflows(data, credentials)
	if err != nil {
		errs = errors.Wrap(errs, err.Error())
	}

	return importedPipelines, errs
}

func (g *GitHubEnhancer) Enhance(data *models.Pipeline, importedPipelines []*enhancers.ImportedPipeline, credentials *models.Credentials) (*models.Pipeline, error) {
	for _, importedPipeline := range importedPipelines {
		data = mergePipelines(data, importedPipeline)
	}

	return data, nil
}

func mergePipelines(pipeline *models.Pipeline, importedPipeline *enhancers.ImportedPipeline) *models.Pipeline {
	if pipeline == nil || pipeline.Jobs == nil {
		return pipeline
	}

	for _, job := range pipeline.Jobs {
		if *job.Name == importedPipeline.JobName {
			job.Imports.Pipeline = importedPipeline.Pipeline
		}
	}

	return pipeline
}
