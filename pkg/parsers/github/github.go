package github

import (
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/parsers/github/models"
	"gopkg.in/yaml.v3"
)

func Parse(data []byte) (*models.Pipeline, error) {
	var err error
	workflow := &githubModels.Workflow{}
	if err := yaml.Unmarshal(data, workflow); err != nil {
		return nil, err
	}

	pipeline := &models.Pipeline{
		Name: &workflow.Name,
	}

	if pipeline.Triggers, err = parseWorkflowTriggers(workflow); err != nil {
		return nil, err
	}

	if workflow.Jobs != nil {
		if pipeline.Jobs, err = parseWorkflowJobs(workflow); err != nil {
			return nil, err
		}
	}

	if pipeline.Defaults, err = parseWorkflowDefaults(workflow); err != nil {
		return nil, err
	}

	return pipeline, nil
}

func parseWorkflowDefaults(workflow *githubModels.Workflow) (*models.Defaults, error) {
	if workflow.Permissions == nil && workflow.Env == nil {
		return nil, nil
	}

	defaults := &models.Defaults{}
	if workflow.Permissions != nil {
		permissions, err := parseTokenPermissions(workflow.Permissions)
		if err != nil {
			return nil, err
		}
		defaults.TokenPermissions = permissions
	}

	if workflow.Env != nil {
		defaults.EnvironmentVariables = workflow.Env
	}

	return defaults, nil
}
