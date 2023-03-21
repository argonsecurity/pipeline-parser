package github

import (
	"errors"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type GitHubEnhancer struct{}

func (g *GitHubEnhancer) Enhance(data *models.Pipeline, credentials *models.Credentials) (*models.Pipeline, error) {
	var errs error
	data, err := enhanceReusableWorkflows(data, credentials)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	return data, errs
}
