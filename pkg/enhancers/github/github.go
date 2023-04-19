package github

import (
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type GitHubEnhancer struct{}

func (g *GitHubEnhancer) LoadImportedPipelines(data *models.Pipeline, credentials *models.Credentials, _ string) ([]*enhancers.ImportedPipeline, error) {
	importedPipelines, err := getReusableWorkflows(data, credentials)
	if err != nil {
		return importedPipelines, err
	}

	return importedPipelines, nil
}

func (g *GitHubEnhancer) Enhance(data *models.Pipeline, importedPipelines []*enhancers.ImportedPipeline) (*models.Pipeline, error) {
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

func (g *GitHubEnhancer) InheritParentPipelineData(parent, child *models.Pipeline) *models.Pipeline {
	return child
}
