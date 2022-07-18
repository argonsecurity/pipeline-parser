package azure

import (
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type AzureParser struct{}

func (g *AzureParser) Parse(workflow *azureModels.Pipeline) (*models.Pipeline, error) {
	// var err error
	pipeline := &models.Pipeline{
		Name: &workflow.Name,
	}

	pipeline.Triggers = parsePipelineTriggers(workflow)

	return pipeline, nil
}
