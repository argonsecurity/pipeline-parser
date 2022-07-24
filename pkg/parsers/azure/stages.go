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

	return parseJobs(stage.Jobs)
}
