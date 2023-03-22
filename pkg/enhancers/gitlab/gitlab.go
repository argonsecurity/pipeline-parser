package gitlab

import (
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type GitLabEnhancer struct{}

func (g *GitLabEnhancer) LoadImportedPipelines(data *models.Pipeline, credentials *models.Credentials) ([]*enhancers.ImportedPipeline, error) {
	return nil, nil
}

func (g *GitLabEnhancer) Enhance(data *models.Pipeline, importedPipelines []*enhancers.ImportedPipeline, credentials *models.Credentials) (*models.Pipeline, error) {
	return data, nil
}
