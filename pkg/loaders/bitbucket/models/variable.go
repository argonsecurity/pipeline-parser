package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type CustomStepVariable struct {
	Name          *string  `yaml:"name,omitempty"` // Name of a variable for the custom pipeline
	Default       *string  `yaml:"default,omitempty"`
	AllowedValues []*string `yaml:"allowed-values,omitempty"`
	FileReference *models.FileReference
}

func (v *CustomStepVariable) UnmarshalYAML(node *yaml.Node) error {
	v.FileReference = loadersUtils.GetFileReference(node)
	return loadersUtils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		switch key {
		case "allowed-values":
			var items []*string
			if err := loadersUtils.ParseSequenceOrOne(value, &items); err != nil {
				return err
			}
			v.AllowedValues = items
			return nil
		case "default":
			str, err := decodeValue(value)
			if err != nil {
				return err
			}
			v.Default = str
		case "name":
			str, err := decodeValue(value)
			if err != nil {
				return err
			}
			v.Name = str
		}
		return nil
	}, "Variable")
}

func decodeValue(value *yaml.Node) (*string, error) {
	var def *string
	if err := value.Decode(&def); err != nil {
		return nil, err
	}
	return def, nil
}
