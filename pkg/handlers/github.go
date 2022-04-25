package handlers

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders"
	githubLoader "github.com/argonsecurity/pipeline-parser/pkg/loaders/github"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/parsers"
	githubParser "github.com/argonsecurity/pipeline-parser/pkg/parsers/github"
)

type GitHubHandler struct{}

func (g *GitHubHandler) GetLoader() loaders.Loader[githubModels.Workflow] {
	return &githubLoader.GitHubLoader{}
}

func (g *GitHubHandler) GetParser() parsers.Parser[githubModels.Workflow] {
	return &githubParser.GitHubParser{}
}
