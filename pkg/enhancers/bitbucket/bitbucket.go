package bitbucket

import (
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type BitbucketEnhancer struct{}

func (b *BitbucketEnhancer) LoadImportedPipelines(data *models.Pipeline, credentials *models.Credentials) ([]*enhancers.ImportedPipeline, error) {
	return nil, nil
}

func (b *BitbucketEnhancer) Enhance(data *models.Pipeline, importedPipelines []*enhancers.ImportedPipeline, credentials *models.Credentials) (*models.Pipeline, error) {
	return data, nil
}
