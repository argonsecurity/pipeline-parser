package common

import (
	"strconv"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"gopkg.in/yaml.v3"
)

type Retry struct {
	When *[]string `yaml:"when,omitempty"`
	Max  *int      `yaml:"max,omitempty"`
}

func (r *Retry) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == consts.IntTag { // format: "retry: 3"
		parsedInt, _ := strconv.Atoi(node.Value)
		r.Max = &parsedInt
		return nil
	}

	return utils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		switch key {
		case "when":
			if value.Tag == consts.SequenceTag {
				parsedStrings, _ := utils.ParseYamlStringSequenceToSlice(value)
				r.When = &parsedStrings
			}
			if value.Tag == consts.StringTag {
				r.When = &[]string{value.Value}
			}
		case "max":
			parsedInt, _ := strconv.Atoi(value.Value)
			r.Max = &parsedInt
		}
		return nil
	})
}
