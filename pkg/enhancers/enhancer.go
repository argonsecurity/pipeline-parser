package enhancers

import "github.com/argonsecurity/pipeline-parser/pkg/models"

type ImportedPipeline struct {
	JobName  string
	Data     []byte
	Pipeline *models.Pipeline
}

type Enhancer interface {
	LoadImportedPipelines(data *models.Pipeline, credentials *models.Credentials) ([]*ImportedPipeline, error)
	Enhance(data *models.Pipeline, importedPipelines []*ImportedPipeline, credentials *models.Credentials) (*models.Pipeline, error)
}
