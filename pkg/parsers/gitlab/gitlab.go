package gitlab

import (
	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/parsers/gitlab/triggers"
)

type GitLabParser struct{}

func (g *GitLabParser) Parse(gitlabCIConfiguration *gitlabModels.GitlabCIConfiguration) (*models.Pipeline, error) {
	var err error
	pipeline := &models.Pipeline{}

	pipeline.Defaults = parseDefaults(gitlabCIConfiguration)

	if gitlabCIConfiguration.Workflow != nil {
		pipeline.Triggers, pipeline.Defaults.Conditions = triggers.ParseRules(gitlabCIConfiguration.Workflow.Rules)
	}
	pipeline.Jobs, err = parseJobs(gitlabCIConfiguration)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

func parseDefaults(gitlabCIConfiguration *gitlabModels.GitlabCIConfiguration) *models.Defaults {
	defaults := &models.Defaults{
		EnvironmentVariables: parseEnvironmentVariables(gitlabCIConfiguration.Variables),
		Runner:               parseRunner(gitlabCIConfiguration.Image),
	}
	return defaults
}
