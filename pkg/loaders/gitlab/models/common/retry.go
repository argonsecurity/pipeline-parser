package common

import (
	"strconv"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"gopkg.in/yaml.v3"
)

type Retry struct {
	When *string `yaml:"when,omitempty"`
	Max  *int    `yaml:"max,omitempty"`
}

func (r *Retry) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == consts.IntTag {
		parsedInt, _ := strconv.Atoi(node.Value)
		r.Max = &parsedInt
		return nil
	}

	return consts.NewErrInvalidYamlTag(node.Tag)
}
