package consts

type Platform string

const (
	GitHubPlatform    Platform = "github"
	GitLabPlatform    Platform = "gitlab"
	AzurePlatform     Platform = "azure"
	BitBucketPlatform Platform = "bitbucket"
)

var Platforms = []Platform{
	GitHubPlatform,
	GitLabPlatform,
	AzurePlatform,
	BitBucketPlatform,
}
