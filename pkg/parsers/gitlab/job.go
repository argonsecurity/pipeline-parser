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

func parseJob(jobID string, job *gitlabModels.Job) (*models.Job, error) {
	conditions := triggers.ParseRulesConditions(job.Rules)
	conditions = append(conditions, triggers.ParseControls(job.Except, true))
	conditions = append(conditions, triggers.ParseControls(job.Only, false))

	var continueOnError *bool
	if job.AllowFailure != nil {
		continueOnError = job.AllowFailure.Enabled
	}

	parsedJob := &models.Job{
		ID:                   &jobID,
		Name:                 &jobID,
		ContinueOnError:      continueOnError,
		ConcurrencyGroup:     &job.Stage,
		PreSteps:             parseScript(job.BeforeScript),
		PostSteps:            parseScript(job.AfterScript),
		Steps:                parseScript(job.Script),
		TimeoutMS:            utils.GetPtr(parseTimeoutString(job.Timeout)),
		EnvironmentVariables: parseEnvironmentVariables(job.Variables),
		Tags:                 job.Tags,
		Runner:               parseRunner(job.Image),
		Conditions:           conditions,
	}
	return parsedJob, nil
}

func parseScript(script *common.Script) []*models.Step {
	if script == nil {
		return nil
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
			Line:   script.FileReference.EndRef.Line + commandIndex + 1,
			Column: script.FileReference.EndRef.Column + len(script.Commands[commandIndex]),
		},
	}

}
