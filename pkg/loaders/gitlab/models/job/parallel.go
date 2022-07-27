package job

import (
	"strconv"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"gopkg.in/yaml.v3"
)

type Parallel struct {
	Max    *int    `yaml:"max,omitempty"`
	Matrix *Matrix `yaml:"matrix,omitempty"`
}

type Matrix map[string][]string

func (p *Parallel) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == consts.IntTag {
		intValue, err := strconv.Atoi(node.Value)
		if err != nil {
			return err
		}
		p.Max = &intValue
		return nil
	}

	return utils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		switch key {
		case "max":
			v, _ := strconv.Atoi(value.Value)
			p.Max = &v
		case "matrix":
			p.Matrix = &Matrix{}
			if err := p.Matrix.UnmarshalYAML(value); err != nil {
				return err
			}
		}
		return nil
	}, "Parallel")
}

func (m *Matrix) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag != consts.MapTag {
		return consts.NewErrInvalidYamlTag(node.Tag, "Matrix")
	}

	return utils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		if value.Tag == consts.StringTag {
			(*m)[key] = []string{value.Value}
		}

		if value.Tag == consts.SequenceTag {
			parsedStrings, err := utils.ParseYamlStringSequenceToSlice(value, "Matrix")
			if err != nil {
				return err
			}
			(*m)[key] = parsedStrings
		}
		return nil
	}, "Matrix")
}
