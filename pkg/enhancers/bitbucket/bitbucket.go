package bitbucket

import "github.com/argonsecurity/pipeline-parser/pkg/models"

type BitbucketEnhancer struct{}

func (b *BitbucketEnhancer) Enhance(data *models.Pipeline, credentials *models.Credentials) (*models.Pipeline, error) {
	return data, nil
}
