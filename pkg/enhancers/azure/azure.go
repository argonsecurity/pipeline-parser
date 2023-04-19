package azure

import (
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type AzureEnhancer struct{}

func (a *AzureEnhancer) LoadImportedPipelines(data *models.Pipeline, credentials *models.Credentials, organization string) ([]*enhancers.ImportedPipeline, error) {
	importedPipelines, err := getTemplates(data, credentials, organization)
	if err != nil {
		return importedPipelines, err
	}

	return importedPipelines, nil
}

func (a *AzureEnhancer) Enhance(data *models.Pipeline, importedPipelines []*enhancers.ImportedPipeline) (*models.Pipeline, error) {
	for _, importedPipeline := range importedPipelines {
		data = mergePipelines(data, importedPipeline)
	}

	return data, nil
}

func mergePipelines(pipeline *models.Pipeline, importedPipeline *enhancers.ImportedPipeline) *models.Pipeline {
	if pipeline == nil || pipeline.Jobs == nil {
		return pipeline
	}

	for _, imported := range pipeline.Imports {
		if *imported.Source.Path == importedPipeline.JobName {
			imported.Pipeline = importedPipeline.Pipeline
		}
	}

	return pipeline
}
