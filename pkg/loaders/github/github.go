package github

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"gopkg.in/yaml.v3"
)

func Load(data []byte) (*models.Workflow, error) {
	workflow := &models.Workflow{}
	if err := yaml.Unmarshal(data, workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}
