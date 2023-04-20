package azure

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func parseResources(resources *azureModels.Resources) *models.Resources {
	if resources == nil || len(resources.Repositories) == 0 {
		return nil
	}

	parsedResources := &models.Resources{
		Repositories: make([]*models.ImportSource, 0),
	}

	for _, repo := range resources.Repositories {
		parsedResources.Repositories = append(parsedResources.Repositories, &models.ImportSource{
			RepositoryAlias: &repo.Repository.Repository,
			Reference:       &repo.Repository.Ref,
			Type:            parseRepoType(repo.Repository.Type),
			SCM:             parseRepoSCM(repo.Repository.Type),
			Repository:      &repo.Repository.Name,
		})
	}

	parsedResources.FileReference = resources.FileReference
	return parsedResources
}

func parseRepoType(repoType string) models.SourceType {
	switch repoType {
	case "github", "git":
		return models.SourceTypeRemote
	default:
		return models.SourceTypeLocal
	}
}

func parseRepoSCM(repoType string) models.Platform {
	switch repoType {
	case "github":
		return consts.GitHubPlatform
	default:
		return consts.AzurePlatform
	}
}
