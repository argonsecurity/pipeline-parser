package gitlab

import (
	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/parsers/gitlab/triggers"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func parseJobs(gitlabCIConfiguration *gitlabModels.GitlabCIConfiguration) ([]*models.Job, error) {
	jobs, err := utils.MapToSliceErr(gitlabCIConfiguration.Jobs, parseJob)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func getJobConditions(job *gitlabModels.Job) []*models.Condition {
	conditions := triggers.ParseRulesConditions(job.Rules)
	if parsedExcept := triggers.ParseControls(job.Except, true); parsedExcept != nil {
		conditions = append(conditions, parsedExcept)
	}
	if parsedOnly := triggers.ParseControls(job.Only, false); parsedOnly != nil {
		conditions = append(conditions, parsedOnly)
	}
	return conditions
}

func getJobContinueOnError(job *gitlabModels.Job) *bool {
	if job.AllowFailure != nil {
		return job.AllowFailure.Enabled
	}
	return nil
}

func parseJob(jobID string, job *gitlabModels.Job) (*models.Job, error) {
	parsedJob := &models.Job{
		ID:                   &jobID,
		Name:                 &jobID,
		ContinueOnError:      getJobContinueOnError(job),
		ConcurrencyGroup:     &job.Stage,
		PreSteps:             parseScript(job.BeforeScript),
		PostSteps:            parseScript(job.AfterScript),
		Steps:                parseScript(job.Script),
		EnvironmentVariables: parseEnvironmentVariables(job.Variables),
		Tags:                 job.Tags,
		Runner:               parseRunner(job.Image),
		Conditions:           getJobConditions(job),
		// TimeoutMS:            utils.GetPtr(parseTimeoutString(job.Timeout)),
	}
	return parsedJob, nil
}

func parseScript(script *common.Script) []*models.Step {
	if script == nil {
		return nil
	}

	if len(script.Commands) == 1 {
		return []*models.Step{
			{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: &script.Commands[0],
				},
				FileReference: script.FileReference,
			},
		}
	}

	return utils.MapWithIndex(script.Commands, func(command string, index int) *models.Step {
		return &models.Step{
			Type: models.ShellStepType,
			Shell: &models.Shell{
				Script: &command,
			},
			FileReference: parseCommandFileReference(script, index),
		}
	})

}

func parseCommandFileReference(script *common.Script, commandIndex int) *models.FileReference {
	return &models.FileReference{
		StartRef: &models.FileLocation{
			Line:   script.FileReference.StartRef.Line + commandIndex + 1,
			Column: script.FileReference.StartRef.Column,
		},
		EndRef: &models.FileLocation{
			Line:   script.FileReference.StartRef.Line + commandIndex + 1,
			Column: script.FileReference.EndRef.Column + len(script.Commands[commandIndex]),
		},
	}

}
