package consts

type Platform string

const (
	GitHubPlatform Platform = "github"
	GitLabPlatform Platform = "gitlab"
	AzurePlatform  Platform = "azure"
)

var Platforms = []Platform{
	GitHubPlatform,
	GitLabPlatform,
	AzurePlatform,
}
