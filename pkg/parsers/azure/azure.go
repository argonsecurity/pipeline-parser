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

	pipeline.Defaults = parsePipelineDefaults(azurePipeline)
	pipeline.Triggers = parsePipelineTriggers(azurePipeline)
	pipeline.Parameters = parseParameters(azurePipeline.Parameters)
	pipeline.Imports = parseExtends(azurePipeline.Extends)

	var jobs []*models.Job

	if azurePipeline.Stages != nil {
		jobs = append(jobs, parseStages(azurePipeline.Stages)...)
	}

	if azurePipeline.Jobs != nil {
		jobs = append(pipeline.Jobs, parseJobs(azurePipeline.Jobs)...)
	}

	if len(jobs) == 0 {
		jobs = []*models.Job{generateDefaultJob()}
		if azurePipeline.Pool != nil {
			jobs[0].Runner = parsePool(azurePipeline.Pool, jobs[0].Runner)
		}

		if azurePipeline.Container != nil {
			jobs[0].Runner = parseContainer(azurePipeline.Container, jobs[0].Runner)
		}
	}

	pipeline.Jobs = jobs

	if azurePipeline.Steps != nil {
		pipeline.Jobs[0].Steps = parseSteps(azurePipeline.Steps)
	}

	// todo: flatten imports ?

	return pipeline, nil
}

func parsePipelineDefaults(pipeline *azureModels.Pipeline) *models.Defaults {
	if pipeline == nil {
		return nil
	}

	defaults := &models.Defaults{
		ContinueOnError: pipeline.ContinueOnError,
	}

	if pipeline.Variables != nil {
		defaults.EnvironmentVariables = parseVariables(pipeline.Variables)
	}

	if pipeline.Resources != nil {
		defaults.Resources = parseResources(pipeline.Resources)
	}

	return defaults
}

func generateDefaultJob() *models.Job {
	return &models.Job{
		Name:   utils.GetPtr("default"),
		Runner: &models.Runner{},
	}
}
