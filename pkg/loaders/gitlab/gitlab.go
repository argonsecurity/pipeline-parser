package gitlab

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"gopkg.in/yaml.v3"
)

type GitLabLoader struct{}

func (g *GitLabLoader) Load(data []byte) (*models.GitlabCIConfiguration, error) {
	gitlabCIConfig := &models.GitlabCIConfiguration{}
	err := yaml.Unmarshal(data, gitlabCIConfig)
	return gitlabCIConfig, err
}
