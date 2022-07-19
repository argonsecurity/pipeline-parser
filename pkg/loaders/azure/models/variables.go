package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Variables []Variable

type Variable struct {
	Name          string `yaml:"name,omitempty"`
	Value         string `yaml:"value,omitempty"`
	Readonly      bool   `yaml:"readonly,omitempty"`
	Group         string `yaml:"group,omitempty"`
	Template      `yaml:",inline"`
	FileReference *models.FileReference
}

func (v *Variables) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == consts.MapTag {
		*v = parseVariablesMap(node)
		return nil
	}

	variables, err := parseVariablesList(node)
	if err != nil {
		return err
	}

	*v = variables
	return nil
}

func parseVariablesMap(node *yaml.Node) []Variable {
	variables := []Variable{}

	loadersUtils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		variable := Variable{
			Name:  key,
			Value: value.Value,
		}

		variable.FileReference = loadersUtils.GetFileReference(value)

		variables = append(variables, variable)
		return nil
	})

	return variables
}

func parseVariablesList(node *yaml.Node) ([]Variable, error) {
	variables := []Variable{}
	for _, variableNode := range node.Content {
		var variable Variable
		if err := variableNode.Decode(&variable); err != nil {
			return nil, err
		}
		variable.FileReference = loadersUtils.GetFileReference(variableNode)
		variables = append(variables, variable)
	}
	return variables, nil
}
