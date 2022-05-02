package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type EnvironmentVariablesRef struct {
	models.EnvironmentVariables
	FileLocation *models.FileLocation
}

func (e *EnvironmentVariablesRef) UnmarshalYAML(node *yaml.Node) error {
	var env models.EnvironmentVariables
	if err := node.Decode(&env); err != nil {
		return err
	}

	e.EnvironmentVariables = env
	e.FileLocation = loadersUtils.GetFileLocation(node)
	return nil
}
