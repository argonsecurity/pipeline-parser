package job

import (
	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/parsers/gitlab/common"
	"github.com/argonsecurity/pipeline-parser/pkg/parsers/gitlab/triggers"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func ParseJobs(gitlabCIConfiguration *gitlabModels.GitlabCIConfiguration) ([]*models.Job, error) {
	jobs, err := utils.MapToSliceErr(gitlabCIConfiguration.Jobs, parseJob)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func parseJob(jobID string, job *gitlabModels.Job) (*models.Job, error) {
	parsedJob := &models.Job{
		ID:                   &jobID,
		Name:                 &jobID,
		ContinueOnError:      getJobContinueOnError(job),
		ConcurrencyGroup:     getJobConcurrencyGroup(job),
		Dependencies:         parseDependencies(job),
		PreSteps:             common.ParseScript(job.BeforeScript),
		PostSteps:            common.ParseScript(job.AfterScript),
		Steps:                common.ParseScript(job.Script),
		EnvironmentVariables: common.ParseEnvironmentVariables(job.Variables),
		Tags:                 job.Tags,
		Runner:               common.ParseRunner(job.Image),
		Conditions:           getJobConditions(job),
		FileReference:        job.FileReference,
	}
	return parsedJob, nil
}

func getJobContinueOnError(job *gitlabModels.Job) *bool {
	if job.AllowFailure != nil {
		return job.AllowFailure.Enabled
	}
	return nil
}

func getJobConcurrencyGroup(job *gitlabModels.Job) *models.ConcurrencyGroup {
	if job.Stage == "" {
		return nil
	}

	return utils.GetPtr(models.ConcurrencyGroup(job.Stage))
}

func getJobConditions(job *gitlabModels.Job) []*models.Condition {
	conditions := triggers.ParseConditionRules(job.Rules)
	if parsedExcept := triggers.ParseControls(job.Except, true); parsedExcept != nil {
		conditions = append(conditions, parsedExcept)
	}
	if parsedOnly := triggers.ParseControls(job.Only, false); parsedOnly != nil {
		conditions = append(conditions, parsedOnly)
	}
	return conditions
}
