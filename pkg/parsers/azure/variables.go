package azure

import (
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func parseVariables(variables *azureModels.Variables) *models.EnvironmentVariablesRef {
	if variables == nil || len(*variables) == 0 {
		return nil
	}

	env := make(models.EnvironmentVariables)

	for _, variable := range *variables {
		if variable.Name != "" {
			env[variable.Name] = variable.Value
		}
	}

	return &models.EnvironmentVariablesRef{
		EnvironmentVariables: env,
		FileReference: &models.FileReference{
			StartRef: &models.FileLocation{
				Line:   (*variables)[0].FileReference.StartRef.Line,
				Column: (*variables)[0].FileReference.StartRef.Column,
			},
			EndRef: &models.FileLocation{
				Line:   (*variables)[len(*variables)-1].FileReference.EndRef.Line,
				Column: (*variables)[len(*variables)-1].FileReference.EndRef.Column,
			},
		},
	}

}
