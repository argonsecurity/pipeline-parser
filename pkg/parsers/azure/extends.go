package azure

import (
	"strings"

	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func parseExtends(extends *azureModels.Extends) []*models.Import {
	if extends == nil {
		return nil
	}

	path, alias := parseTemplateString(extends.Template.Template)
	parameters, paramImports := parseExtendParameters(extends.Parameters)

	imports := []*models.Import{{
		FileReference: extends.FileReference,
		Parameters:    parameters,
		Source: &models.ImportSource{
			Path:            &path,
			RepositoryAlias: &alias,
		},
	}}

	imports = append(imports, paramImports...)

	return imports
}

func parseTemplateString(template string) (path string, alias string) {
	// common.yml@templates  # Template reference
	parts := strings.Split(template, "@")
	path = parts[0]
	alias = ""
	if len(parts) == 2 {
		alias = parts[1]
	}
	return path, alias
}

func parseExtendParameters(params map[string]any) (parameters map[string]any, imports []*models.Import) {
	if params == nil {
		return nil, nil
	}

	for key, param := range params {
		value, ok := param.(azureModels.Template)
		if ok {
			path, alias := parseTemplateString(value.Template)
			parameters, paramImports := parseExtendParameters(value.Parameters)

			imports = append(imports, &models.Import{
				Parameters: parameters,
				Source: &models.ImportSource{
					Path:            &path,
					RepositoryAlias: &alias,
				},
			})
			imports = append(imports, paramImports...)
			continue
		}
		if parameters == nil {
			parameters = make(map[string]any)
		}
		parameters[key] = param
	}

	return parameters, imports
}
