package azure

import (
	"strings"

	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/mitchellh/mapstructure"
)

func parseExtends(extends *azureModels.Extends) []*models.Import {
	if extends == nil {
		return nil
	}

	path, alias := parseTemplateString(extends.Template.Template)
	parameters, paramImports := parseExtendParameters(extends.Parameters, extends.FileReference)
	imports := []*models.Import{{
		FileReference: extends.FileReference,
		Parameters:    parameters,
		Source: &models.ImportSource{
			Path:            &path,
			RepositoryAlias: &alias,
			Type:            calculateSourceType(alias),
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

func parseExtendParameters(params map[string]any, rootFileRef *models.FileReference) (parameters map[string]any, imports []*models.Import) {
	if params == nil {
		return nil, nil
	}

	for key, param := range params {
		items := []any{param}
		if utils.IsArray(param) {
			items = append(items, param.([]any)...)
		}
		for _, item := range items {
			value, ok := tryToParseTemplate(item)
			if ok {
				path, alias := parseTemplateString(value.Template)
				parameters, paramImports := parseExtendParameters(value.Parameters, rootFileRef)

				imports = append(imports, &models.Import{
					Parameters: parameters,
					Source: &models.ImportSource{
						Path:            &path,
						RepositoryAlias: &alias,
						Type:            calculateSourceType(alias),
					},
					FileReference: rootFileRef,
				})
				imports = append(imports, paramImports...)
				continue
			}
			if parameters == nil {
				parameters = make(map[string]any)
			}
			parameters[key] = item
		}
	}

	return parameters, imports
}

func tryToParseTemplate(input any) (azureModels.Template, bool) {
	var azureTemplate azureModels.Template
	return azureTemplate, mapstructure.Decode(input, &azureTemplate) == nil
}

func calculateSourceType(alias string) models.SourceType {
	if alias == "self" || alias == "" {
		return models.SourceTypeLocal
	}
	return models.SourceTypeRemote
}
