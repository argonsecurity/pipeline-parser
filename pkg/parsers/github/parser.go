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

	pipeline := &models.Pipeline{}
	pipeline.Jobs = parseWorkflowJobs(workflow)
	pipeline.Triggers, err = parseWorkflowTriggers(workflow)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}
