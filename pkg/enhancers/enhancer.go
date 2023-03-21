package enhancers

import "github.com/argonsecurity/pipeline-parser/pkg/models"

type Enhancer interface {
	Enhance(data *models.Pipeline, credentials *models.Credentials) (*models.Pipeline, error)
}
