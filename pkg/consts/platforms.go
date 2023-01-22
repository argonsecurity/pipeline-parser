package consts

type Platform string

const (
	GitHubPlatform    Platform = "github"
	GitLabPlatform    Platform = "gitlab"
	AzurePlatform     Platform = "azure"
	BitbucketPlatform Platform = "bitbucket"
)

var Platforms = []Platform{
	GitHubPlatform,
	GitLabPlatform,
	AzurePlatform,
	BitbucketPlatform,
}
