package github

import (
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/parsers/github/models"
	"gopkg.in/yaml.v3"
)

func Parse(data []byte) (*models.Pipeline, error) {
	workflow := &githubModels.Workflow{}
	if err := yaml.Unmarshal(data, workflow); err != nil {
		return nil, err
	}

	pipeline := &models.Pipeline{
		Name: &workflow.Name,
	}
	triggers, err := parseWorkflowTriggers(workflow)
	if err != nil {
		return nil, err
	}
	pipeline.Triggers = triggers
	if workflow.Jobs != nil {
		pipeline.Jobs = parseWorkflowJobs(workflow)
	}
	pipeline.Defaults = parseWorkflowDefaults(workflow)

	return pipeline, nil
}

func parseWorkflowDefaults(workflow *githubModels.Workflow) *models.Defaults {
	if workflow.Permissions == nil && workflow.Env == nil {
		return nil
	}

	return &models.Defaults{
		TokenPermissions:     parseTokenPermissions(workflow.Permissions),
		EnvironmentVariables: workflow.Env,
	}
}
