package azure

import (
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type AzureParser struct{}

func (g *AzureParser) Parse(azurePipeline *azureModels.Pipeline) (*models.Pipeline, error) {
	// var err error
	pipeline := &models.Pipeline{
		Name: &azurePipeline.Name,
	}

	pipeline.Triggers = parsePipelineTriggers(azurePipeline)
	pipeline.Parameters = parseParameters(azurePipeline)
	pipeline.Imports = parseExtends(azurePipeline.Extends)

	if azurePipeline.ContinueOnError != nil {
		pipeline.Defaults = &models.Defaults{
			ContinueOnError: azurePipeline.ContinueOnError,
		}
	}

	if azurePipeline.Jobs != nil {
		pipeline.Jobs = parseJobs(azurePipeline.Jobs)
	}

	return pipeline, nil
}
