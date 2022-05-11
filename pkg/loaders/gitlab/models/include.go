package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Include []IncludeItem

func (i *Include) UnmarshalYAML(node *yaml.Node) error {
	var items []IncludeItem
	if err := utils.ParseSliceOrOne(node, &items); err != nil {
		return err
	}
	*i = items
	return nil
}

type IncludeItem struct {
	Project  string   `yaml:"project"`
	Ref      string   `yaml:"ref"`
	Template string   `yaml:"template"`
	File     []string `yaml:"file"`

	Local  string `yaml:"local"`
	Remote string `yaml:"remote"`

	FileReference *models.FileReference
}

func (it *IncludeItem) UnmarshalYAML(node *yaml.Node) error {
	it.FileReference = utils.GetFileReference(node)
	return utils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		switch key {
		case "project":
			it.Project = value.Value
		case "ref":
			it.Ref = value.Value
		case "file":
			return utils.ParseSliceOrOne(value, &it.File)
		case "template":
			it.Template = value.Value
		case "local":
			it.Local = value.Value
		case "remote":
			it.Remote = value.Value
		}
		return nil
	})
}
