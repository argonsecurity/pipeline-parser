package azure

import (
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func parseParameters(parameters *azureModels.Parameters) []*models.Parameter {
	if parameters == nil || len(*parameters) == 0 {
		return nil
	}

	var parsedParameters []*models.Parameter
	for _, param := range *parameters {
		name := param.Name

		parsedParameters = append(parsedParameters, &models.Parameter{
			Name:          &name,
			Default:       param.Default,
			Options:       param.Values,
			FileReference: param.FileReference,
		})
	}
	return parsedParameters
}
