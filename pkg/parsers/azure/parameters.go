package azure

import (
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func parseParameters(pipeline *azureModels.Pipeline) []*models.Parameter {
	if pipeline == nil || pipeline.Parameters == nil {
		return nil
	}

	var parameters []*models.Parameter
	for _, param := range *pipeline.Parameters {
		name := param.Name

		parameters = append(parameters, &models.Parameter{
			Name: &name,
			// Value:         &param.Values, // TODO: check what value means
			Default:       param.Default,
			FileReference: param.FileReference,
		})
	}
	return parameters
}
