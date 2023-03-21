package handler

import (
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	azureEnhancer "github.com/argonsecurity/pipeline-parser/pkg/enhancers/azure"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders"
	azureLoader "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure"
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/parsers"
	azureParser "github.com/argonsecurity/pipeline-parser/pkg/parsers/azure"
)

type AzureHandler struct{}

func (g *AzureHandler) GetLoader() loaders.Loader[azureModels.Pipeline] {
	return &azureLoader.AzureLoader{}
}

func (g *AzureHandler) GetParser() parsers.Parser[azureModels.Pipeline] {
	return &azureParser.AzureParser{}
}

func (g *AzureHandler) GetEnhancer() enhancers.Enhancer {
	return &azureEnhancer.AzureEnhancer{}
}
