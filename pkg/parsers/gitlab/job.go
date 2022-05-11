package gitlab

import (
	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func parseJobs(gitlabCIConfiguration *gitlabModels.GitlabCIConfiguration) (*[]models.Job, error) {
	jobs, err := utils.MapToSliceErr(gitlabCIConfiguration.Jobs, parseJob)
	if err != nil {
		return nil, err
	}

	return &jobs, nil
}

func parseJob(jobID string, job *gitlabModels.Job) (models.Job, error) {
	parsedJob := models.Job{
		ID:                   &jobID,
		Name:                 &jobID,
		ContinueOnError:      job.AllowFailure.Enabled,
		ConcurrencyGroup:     &job.Stage,
		PreSteps:             parseScript(job.BeforeScript),
		PostSteps:            parseScript(job.AfterScript),
		Steps:                parseScript(job.Script),
		TimeoutMS:            utils.GetPtr(parseTimeoutString(job.Timeout)),
		EnvironmentVariables: parseEnvironmentVariables(*job.Variables),
		Tags:                 &job.Tags,
		Runner:               parseRunner(job.Image),
		// Conditions:       job.Rules,
	}
	return parsedJob, nil
}

func parseScript(script *common.Script) *[]models.Step {
	return utils.GetPtr(
		utils.MapWithIndex(script.Commands, func(command string, index int) models.Step {
			return models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: &command,
				},
				FileReference: parseCommandFileReference(script, index),
			}
		}),
	)
}

func parseCommandFileReference(script *common.Script, commandIndex int) *models.FileReference {
	return &models.FileReference{
		StartRef: &models.FileLocation{
			Line:   script.FileReference.StartRef.Line + commandIndex + 1,
			Column: script.FileReference.StartRef.Column,
		},
		EndRef: &models.FileLocation{
			Line:   script.FileReference.EndRef.Line + commandIndex + 1,
			Column: script.FileReference.EndRef.Column + len(script.Commands[commandIndex]),
		},
	}

}
