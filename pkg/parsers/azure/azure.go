package azure

import (
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

type AzureParser struct{}

func (g *AzureParser) Parse(azurePipeline *azureModels.Pipeline) *models.Pipeline {
	if azurePipeline == nil {
		return nil
	}

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
	} else {
		pipeline.Jobs = []*models.Job{generateDefaultJob()}
		if azurePipeline.Steps != nil {
			pipeline.Jobs[0].Steps = parseSteps(azurePipeline.Steps)
		}
	}

	return pipeline
}

func generateDefaultJob() *models.Job {
	return &models.Job{
		Name: utils.GetPtr("default"),
	}
}
