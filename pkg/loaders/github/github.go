package github

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"gopkg.in/yaml.v3"
)

type GitHubLoader struct{}

func (g *GitHubLoader) Load(data []byte) (*models.Workflow, error) {
	var workflow *models.Workflow
	err := yaml.Unmarshal(data, workflow)
	return workflow, err
}
