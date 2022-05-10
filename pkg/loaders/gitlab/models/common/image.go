package common

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"gopkg.in/yaml.v3"
)

type Image struct {
	Name       string   `yaml:"name"`
	Entrypoint []string `yaml:"entrypoint"`
}

func (im *Image) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == consts.StringTag {
		im.Name = node.Value
		return nil
	}

	if node.Tag == consts.MapTag {
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i].Value
			value := node.Content[i+1]
			switch key {
			case "name":
				im.Name = value.Value
			case "entrypoint":
				entrypoints, err := utils.ParseYamlStringSequenceToSlice(value)
				if err != nil {
					return err
				}
				im.Entrypoint = entrypoints
			}
		}
		return nil
	}

	return consts.NewErrInvalidYamlTag(node.Tag)
}
