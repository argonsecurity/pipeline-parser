package handler

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	gitlabEnhancer "github.com/argonsecurity/pipeline-parser/pkg/enhancers/gitlab"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders"
	gitlabLoader "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab"
	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/parsers"
	gitlabParser "github.com/argonsecurity/pipeline-parser/pkg/parsers/gitlab"
)

type GitLabHandler struct{}

func (g *GitLabHandler) GetPlatform() models.Platform {
	return consts.GitLabPlatform
}

func (g *GitLabHandler) GetLoader() loaders.Loader[gitlabModels.GitlabCIConfiguration] {
	return &gitlabLoader.GitLabLoader{}
}

func (g *GitLabHandler) GetParser() parsers.Parser[gitlabModels.GitlabCIConfiguration] {
	return &gitlabParser.GitLabParser{}
}

func (g *GitLabHandler) GetEnhancer() enhancers.Enhancer {
	return &gitlabEnhancer.GitLabEnhancer{}
}
