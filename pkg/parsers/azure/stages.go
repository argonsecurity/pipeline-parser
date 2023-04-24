package azure

import (
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func parseStages(stages *azureModels.Stages) []*models.Job {
	if stages == nil || (stages.Stages == nil && stages.TemplateStages == nil) {
		return nil
	}

	var jobs []*models.Job

	for _, stage := range stages.Stages {
		if stage.Jobs != nil {
			jobs = append(jobs, parseStage(stage)...)
		}
	}

	for _, stage := range stages.TemplateStages {
		if stage.Template.Template != "" {
			jobs = append(jobs, parseTemplateStage(stage))
		}
	}

	return jobs
}

func parseStage(stage *azureModels.Stage) []*models.Job {
	if stage == nil || stage.Jobs == nil {
		return nil
	}

	parsedJobs := parseJobs(stage.Jobs)

	if stage.Variables == nil {
		return parsedJobs
	}
	envs := parseVariables(stage.Variables)

	for _, job := range parsedJobs {
		if job.EnvironmentVariables == nil {
			job.EnvironmentVariables = envs
			continue
		}

		for k, v := range envs.EnvironmentVariables {
			job.EnvironmentVariables.EnvironmentVariables[k] = v
		}
	}

	return parsedJobs
}

func parseTemplateStage(stage *azureModels.TemplateStage) *models.Job {
	path, alias := parseTemplateString(stage.Template.Template)
	return &models.Job{
		ID: &stage.Template.Template,
		Imports: &models.Import{
			Source: &models.ImportSource{
				Path:            &path,
				Type:            calculateSourceType(alias),
				RepositoryAlias: &alias,
			},
			Parameters:    stage.Parameters,
			FileReference: stage.FileReference,
		},
		FileReference: stage.FileReference,
	}
}
