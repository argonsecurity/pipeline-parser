package job

import (
	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func parseDependencies(job *gitlabModels.Job) []*models.JobDependency {
	dependencies := parseJobDependencies(job)
	dependencies = append(dependencies, parseJobNeeds(job)...)
	return dependencies
}

func parseJobDependencies(job *gitlabModels.Job) []*models.JobDependency {
	var dependencies []*models.JobDependency
	if job.Dependencies != nil {
		dependencies = append(dependencies, utils.Map(job.Dependencies, func(dependency string) *models.JobDependency {
			return &models.JobDependency{
				JobID: &dependency,
			}
		})...)
	}
	return dependencies
}

func parseJobNeeds(job *gitlabModels.Job) []*models.JobDependency {
	var dependencies []*models.JobDependency
	if job.Needs != nil {
		for _, item := range *job.Needs {
			dependencies = append(dependencies, &models.JobDependency{
				JobID:    utils.GetPtrOrNil(item.Job),
				Pipeline: utils.GetPtrOrNil(item.Pipeline),
			})
		}
	}
	return dependencies
}
