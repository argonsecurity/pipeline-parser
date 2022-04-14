package github

import (
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/parsers/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

const (
	defaultTimeoutMS int = 360 * 60 * 1000
)

func parseWorkflowJobs(workflow *githubModels.Workflow) []models.Job {
	return utils.MapToSlice(workflow.Jobs.NormalJobs, parseNormalJob)
}

// type Job struct {
// 	Steps                *[]Step
// 	Runner               *Runner
// 	Inputs               *[]Parameter
// 	Tags                 *[]string
// }

func parseNormalJob(jobName string, normalJob *githubModels.NormalJob) models.Job {
	job := models.Job{
		ID:                   &jobName,
		Name:                 &jobName,
		ContinueOnError:      &normalJob.ContinueOnError,
		EnvironmentVariables: parseEnvironmentVariables(normalJob.Env),
	}

	if normalJob.TimeoutMinutes != nil {
		timeout := int(*normalJob.TimeoutMinutes) * 60 * 1000
		job.TimeoutMS = &timeout
	} else {
		defaultTimeout := defaultTimeoutMS
		job.TimeoutMS = &defaultTimeout
	}

	if normalJob.If != "" {
		job.Conditions = &[]models.Condition{models.Condition(normalJob.If)}
	}

	if normalJob.Concurrency != nil {
		job.ConcurrencyGroup = normalJob.Concurrency.Group
	}

	if normalJob.Steps != nil {
		steps := parseJobSteps(normalJob.Steps)
		job.Steps = &steps
	}

	return job
}

func parseJobConcurrency(concurrency interface{}) *string {
	if concurrency == nil {
		return nil
	}

	concurrencyAsString, ok := concurrency.(string)
	if ok {
		return &concurrencyAsString
	}

	return nil
}
