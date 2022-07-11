package azure

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"gopkg.in/yaml.v3"
)

type AzureLoader struct{}

func (g *AzureLoader) Load(data []byte) (*models.Pipeline, error) {
	pipeline := &models.Pipeline{}
	err := yaml.Unmarshal(data, pipeline)
	return pipeline, err
}
