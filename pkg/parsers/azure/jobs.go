package azure

import (
	"strconv"

	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

var (
	defaultTimeoutMS int = 60 * 60 * 1000
)

func parseJobs(jobs *azureModels.Jobs) []*models.Job {
	if jobs == nil {
		return nil
	}

	var parsedJobs []*models.Job

	if jobs.CIJobs != nil {
		parsedJobs = utils.Map(jobs.CIJobs, parseCIJob)
	}

	if jobs.DeploymentJobs != nil {
		parsedJobs = append(parsedJobs, utils.Map(jobs.DeploymentJobs, parseDeploymentJob)...)
	}

	if jobs.TemplateJobs != nil {
		parsedJobs = append(parsedJobs, utils.Map(jobs.TemplateJobs, parseTemplateJob)...)
	}

	return parsedJobs
}

func parseTemplateJob(job *azureModels.TemplateJob) *models.Job {
	if job == nil {
		return nil
	}
	path, alias := parseTemplateString(job.Template.Template)
	parsedJob := &models.Job{
		ID: &job.Template.Template,
		Imports: &models.Import{
			Source: &models.ImportSource{
				Path:            &path,
				Type:            calculateSourceType(alias),
				RepositoryAlias: &alias,
			},
			FileReference: job.FileReference,
			Parameters:    job.Parameters,
		},
		FileReference: job.FileReference,
	}

	return parsedJob
}

func parseCIJob(job *azureModels.CIJob) *models.Job {
	if job == nil {
		return nil
	}

	parsedJob := parseBaseJob(&job.BaseJob)

	parsedJob.ID = &job.Job
	parsedJob.FileReference = job.FileReference

	return parsedJob
}

func parseDeploymentJob(job *azureModels.DeploymentJob) *models.Job {
	if job == nil {
		return nil
	}

	parsedJob := parseBaseJob(&job.BaseJob)

	parsedJob.ID = &job.Deployment
	parsedJob.FileReference = job.FileReference

	return parsedJob
}

func parseBaseJob(job *azureModels.BaseJob) *models.Job {
	if job == nil {
		return nil
	}

	parsedJob := &models.Job{
		Name:            &job.DisplayName,
		ContinueOnError: utils.GetPtr(strconv.FormatBool(job.ContinueOnError)),
		Runner:          parseRunner(job),
	}

	if job.Variables != nil {
		parsedJob.EnvironmentVariables = parseVariables(job.Variables)
	}

	if job.TimeoutInMinutes != 0 {
		timeout := int(job.TimeoutInMinutes) * 60 * 1000
		parsedJob.TimeoutMS = &timeout
	} else {
		parsedJob.TimeoutMS = &defaultTimeoutMS
	}

	if job.Condition != "" {
		parsedJob.Conditions = []*models.Condition{{Statement: job.Condition}}
	}

	if job.Steps != nil {
		parsedJob.Steps = parseSteps(job.Steps)
	}

	if job.DependsOn != nil {
		parsedJob.Dependencies = parseDependencies(job.DependsOn)
	}

	return parsedJob
}

func parseDependencies(dependsOn *azureModels.DependsOn) []*models.JobDependency {
	if dependsOn == nil {
		return nil
	}

	return utils.Map(([]string)(*dependsOn), func(dependency string) *models.JobDependency {
		return &models.JobDependency{
			JobID: &dependency,
		}
	})
}
