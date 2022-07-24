package azure

import (
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

type AzureParser struct{}

func (g *AzureParser) Parse(azurePipeline *azureModels.Pipeline) (*models.Pipeline, error) {
	if azurePipeline == nil {
		return nil, nil
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

	var jobs []*models.Job

	if azurePipeline.Stages != nil {
		jobs = append(jobs, parseStages(azurePipeline.Stages)...)
	}

	if azurePipeline.Jobs != nil {
		jobs = append(pipeline.Jobs, parseJobs(azurePipeline.Jobs)...)
	}

	if len(jobs) == 0 {
		jobs = []*models.Job{generateDefaultJob()}
	}

	pipeline.Jobs = jobs

	if azurePipeline.Steps != nil {
		pipeline.Jobs[0].Steps = parseSteps(azurePipeline.Steps)
	}

	return pipeline, nil
}

func generateDefaultJob() *models.Job {
	return &models.Job{
		Name: utils.GetPtr("default"),
	}
}
