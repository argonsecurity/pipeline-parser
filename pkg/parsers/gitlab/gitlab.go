package gitlab

import (
	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type GitLabParser struct{}

func (g *GitLabParser) Parse(gitlabCIConfiguration *gitlabModels.GitlabCIConfiguration) (*models.Pipeline, error) {
	pipeline := &models.Pipeline{}
	var err error

	pipeline.Triggers, err = parseTriggers(gitlabCIConfiguration)
	if err != nil {
		return nil, err
	}

	pipeline.Jobs, err = parseJobs(gitlabCIConfiguration)
	if err != nil {
		return nil, err
	}

	pipeline.Defaults, err = parseDefaults(gitlabCIConfiguration)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

func parseDefaults(gitlabCIConfiguration *gitlabModels.GitlabCIConfiguration) (*models.Defaults, error) {
	defaults := &models.Defaults{
		EnvironmentVariables: parseEnvironmentVariables(*gitlabCIConfiguration.Variables),
	}
	return defaults, nil
}
