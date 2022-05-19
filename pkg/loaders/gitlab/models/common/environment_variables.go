package common

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type EnvironmentVariablesRef struct {
	Variables     *Variables
	FileReference *models.FileReference
}

type Variables map[string]any
type Variable struct {
	Value       string `yaml:"value"`
	Description string `yaml:"description"`
}

func (v *Variables) UnmarshalYAML(node *yaml.Node) error {
	variables := map[string]any{}
	if err := utils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		if value.Tag == consts.MapTag {
			variable := Variable{}
			if err := value.Decode(&variable); err != nil {
				return err
			}
			variables[key] = variable.Value
		} else {
			variables[key] = value.Value
		}
		return nil
	}); err != nil {
		return err
	}

	*v = variables
	return nil
}

func (e *EnvironmentVariablesRef) UnmarshalYAML(node *yaml.Node) error {
	e.FileReference = utils.GetFileReference(node)

	e.FileReference.StartRef.Line--
	e.FileReference.StartRef.Column = 1
	e.FileReference.EndRef.Column += 2 // +2 for the ": " after the key
	return node.Decode(&e.Variables)
}
