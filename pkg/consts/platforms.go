package consts

import "github.com/argonsecurity/pipeline-parser/pkg/models"

const (
	GitHubPlatform    models.Platform = "github"
	GitLabPlatform    models.Platform = "gitlab"
	AzurePlatform     models.Platform = "azure"
	BitbucketPlatform models.Platform = "bitbucket"
)

var Platforms = []models.Platform{
	GitHubPlatform,
	GitLabPlatform,
	AzurePlatform,
	BitbucketPlatform,
}
