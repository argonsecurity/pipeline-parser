package models

import (
	"strings"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Include []IncludeItem

func (i *Include) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == consts.StringTag {
		*i = Include{parseIncludeString(node)}
		return nil
	}

	var items []IncludeItem
	if err := utils.ParseSequenceOrOne(node, &items); err != nil {
		return err
	}
	*i = items
	return nil
}

type IncludeItem struct {
	Project  string `yaml:"project"`
	Ref      string `yaml:"ref"`
	Template string `yaml:"template"`
	File     string `yaml:"file"`

	Local  string `yaml:"local"`
	Remote string `yaml:"remote"`

	FileReference *models.FileReference
}

func (it *IncludeItem) UnmarshalYAML(node *yaml.Node) error {
	it.FileReference = utils.GetFileReference(node)
	if node.Tag == consts.StringTag {
		*it = parseIncludeString(node)
		return nil
	}

	return utils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		switch key {
		case "project":
			it.Project = value.Value
		case "ref":
			it.Ref = value.Value
		case "file":
			it.File = value.Value
		case "template":
			it.Template = value.Value
		case "local":
			it.Local = value.Value
		case "remote":
			it.Remote = value.Value
		}
		return nil
	}, "IncludeItem")
}

func parseIncludeString(node *yaml.Node) IncludeItem {
	it := IncludeItem{}
	it.FileReference = utils.GetFileReference(node)

	if strings.HasPrefix(node.Value, "https://") {
		it.Remote = node.Value
		return it
	}

	it.Local = node.Value
	return it
}
