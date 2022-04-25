package github

import "github.com/argonsecurity/pipeline-parser/pkg/models"

type GitHubEnhancer struct{}

func (g *GitHubEnhancer) EnhanceJob(job models.Job) models.Job {
	return job
}

func (g *GitHubEnhancer) EnhanceStep(step models.Step) models.Step {
	return step
}
