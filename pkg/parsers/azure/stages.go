package azure

import (
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func parseStages(stages *azureModels.Stages) []*models.Job {
	if stages == nil || stages.Stages == nil {
		return nil
	}

	var jobs []*models.Job

	for _, stage := range stages.Stages {
		if stage.Jobs != nil {
			jobs = append(jobs, parseJobs(stage.Jobs)...)
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
