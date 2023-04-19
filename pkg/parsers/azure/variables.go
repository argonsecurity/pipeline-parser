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
	var imports *models.Import

	for _, variable := range *variables {
		if variable.Name != "" {
			env[variable.Name] = variable.Value
		}

		path, alias := parseTemplateString(variable.Template.Template)
		if variable.Template.Template != "" {
			imports = &models.Import{
				Source: &models.ImportSource{
					Path:            &path,
					Type:            calculateSourceType(alias),
					RepositoryAlias: &alias,
				},
				Parameters:    variable.Parameters,
				FileReference: variable.FileReference,
			}
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
		Imports: imports,
	}

}
