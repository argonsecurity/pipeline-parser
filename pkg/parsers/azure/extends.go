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

	// common.yml@templates  # Template reference
	parts := strings.Split(extends.Template.Template, "@")
	path := parts[0]
	alias := ""
	if len(parts) == 2 {
		alias = parts[1]
	}

	imports := []*models.Import{{
		FileReference: extends.FileReference,
		Parameters:    extends.Parameters,
		Source: &models.ImportSource{
			Path:            &path,
			RepositoryAlias: &alias,
		},
	}}

	return imports
}
