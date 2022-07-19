package azure

import (
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func parseEnvironmentVariablesRef(envRef *azureModels.EnvironmentVariablesRef) *models.EnvironmentVariablesRef {
	if envRef == nil {
		return nil
	}

	return &models.EnvironmentVariablesRef{
		EnvironmentVariables: envRef.EnvironmentVariables,
		FileReference:        envRef.FileReference,
	}
}

func parseExtends(extends *azureModels.Extends) []string {
	if extends == nil {
		return nil
	}

	return []string{extends.Template.Template}
}
