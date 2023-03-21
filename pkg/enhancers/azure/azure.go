package azure

import "github.com/argonsecurity/pipeline-parser/pkg/models"

type AzureEnhancer struct{}

func (a *AzureEnhancer) Enhance(data *models.Pipeline, credentials *models.Credentials) (*models.Pipeline, error) {
	return data, nil
}
