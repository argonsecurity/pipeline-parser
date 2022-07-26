package handler

import (
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
