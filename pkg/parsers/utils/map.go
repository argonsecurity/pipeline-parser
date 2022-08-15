package utils

import (
	commonLoaderModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/common/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func ParseMapToParameters(mapNode commonLoaderModels.Map) []*models.Parameter {
	parameters := make([]*models.Parameter, 0)
	for _, entry := range mapNode.Values {
		var key = entry.Key // define key here so the pointer won't change in the loop
		parameters = append(parameters, &models.Parameter{
			Name:          &key,
			Value:         entry.Value,
			FileReference: entry.FileReference,
		})
	}

	return parameters
}
