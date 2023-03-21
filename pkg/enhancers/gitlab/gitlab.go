package gitlab

import "github.com/argonsecurity/pipeline-parser/pkg/models"

type GitLabEnhancer struct{}

func (g *GitLabEnhancer) Enhance(data *models.Pipeline, credentials *models.Credentials) (*models.Pipeline, error) {
	return data, nil
}
