package handler

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders"
	bitbucketLoader "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket"
	bitbucketModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	"github.com/argonsecurity/pipeline-parser/pkg/parsers"
	bitbucketParser "github.com/argonsecurity/pipeline-parser/pkg/parsers/bitbucket"
)

type BitBucketHandler struct{}

func (g *BitBucketHandler) GetLoader() loaders.Loader[bitbucketModels.Pipeline] {
	return &bitbucketLoader.BitbucketLoader{}
}

func (g *BitBucketHandler) GetParser() parsers.Parser[bitbucketModels.Pipeline] {
	return &bitbucketParser.BitbucketParser{}
}
