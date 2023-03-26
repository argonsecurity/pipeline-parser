package azure

import (
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type AzureEnhancer struct{}

func (a *AzureEnhancer) LoadImportedPipelines(data *models.Pipeline, credentials *models.Credentials) ([]*enhancers.ImportedPipeline, error) {
	return nil, nil
}

func (a *AzureEnhancer) Enhance(data *models.Pipeline, importedPipelines []*enhancers.ImportedPipeline) (*models.Pipeline, error) {
	return data, nil
}
