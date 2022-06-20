package consts

type Platform string

const (
	GitHubPlatform Platform = "github"
	GitLabPlatform Platform = "gitlab"
)

var Platforms = []Platform{
	GitHubPlatform,
	GitLabPlatform,
}
