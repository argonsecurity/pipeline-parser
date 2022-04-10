package github

import (
	"github.com/argonsecurity/pipeline-parser/pkg/parsers/github/models"
	"gopkg.in/yaml.v3"
)

func Parse(data []byte) (*models.Workflow, error) {
	workflow := &models.Workflow{}
	if err := yaml.Unmarshal(data, workflow); err != nil {
		return nil, err
	}

	triggers, err := parseWorkflowTriggers(workflow)
	if err != nil {
		return nil, err
	}

	println(triggers)
	return workflow, nil
}
