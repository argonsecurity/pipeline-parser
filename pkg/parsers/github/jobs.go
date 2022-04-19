package github

import (
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/parsers/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

const (
	defaultTimeoutMS int = 360 * 60 * 1000
)

func parseWorkflowJobs(workflow *githubModels.Workflow) *[]models.Job {
	return utils.GetPtr(utils.MapToSlice(workflow.Jobs.NormalJobs, parseJob))
}

func parseJob(jobName string, job *githubModels.Job) models.Job {
	parsedJob := models.Job{
		ID:                   job.ID,
		Name:                 &job.Name,
		ContinueOnError:      &job.ContinueOnError,
		EnvironmentVariables: job.Env,
	}

	if job.Name == "" {
		parsedJob.Name = job.ID
	}

	if job.TimeoutMinutes != nil && *job.TimeoutMinutes == 0 {
		timeout := int(*job.TimeoutMinutes) * 60 * 1000
		parsedJob.TimeoutMS = &timeout
	} else {
		defaultTimeout := defaultTimeoutMS
		parsedJob.TimeoutMS = &defaultTimeout
	}

	if job.If != "" {
		parsedJob.Conditions = &[]models.Condition{models.Condition(job.If)}
	}

	if job.Concurrency != nil {
		parsedJob.ConcurrencyGroup = job.Concurrency.Group
	}

	if job.Steps != nil {
		parsedJob.Steps = parseJobSteps(job.Steps)
	}

	if job.RunsOn != nil {
		parsedJob.Runner = parseRunsOnToRunner(job.RunsOn)
	}

	if job.Needs != nil {
		parsedJob.Dependencies = (*[]string)(job.Needs)
	}

	return parsedJob
}
