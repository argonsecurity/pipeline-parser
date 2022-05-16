package gitlab

import (
	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func parseEnvironmentVariables(environmentVariables *gitlabModels.EnvironmentVariablesRef) *models.EnvironmentVariablesRef {
	if environmentVariables == nil {
		return nil
	}

	return &models.EnvironmentVariablesRef{
		FileReference:        environmentVariables.FileReference,
		EnvironmentVariables: environmentVariables.Variables,
	}
}
