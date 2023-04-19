package enhancers

import "github.com/argonsecurity/pipeline-parser/pkg/models"

type ImportedPipeline struct {
	JobName             string
	OriginFileReference *models.FileReference
	Data                []byte
	Pipeline            *models.Pipeline
}

type Enhancer interface {
	LoadImportedPipelines(data *models.Pipeline, credentials *models.Credentials, organization string) ([]*ImportedPipeline, error)
	Enhance(data *models.Pipeline, importedPipelines []*ImportedPipeline) (*models.Pipeline, error)
}
