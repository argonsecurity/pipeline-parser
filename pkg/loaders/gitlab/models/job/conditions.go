package job

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"gopkg.in/yaml.v3"
)

type Conditions struct {
	Refs       []string
	Variables  []string
	Changes    []string
	Kubernetes string
}

func (c *Conditions) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == consts.SequenceTag {
		refs, err := utils.ParseYamlStringSequenceToSlice(node)
		if err != nil {
			return err
		}
		c.Refs = refs
		return nil
	}

	if node.Tag == consts.MapTag {
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i].Value
			value := node.Content[i+1]

			switch key {
			case "refs":
				refs, err := utils.ParseYamlStringSequenceToSlice(value)
				if err != nil {
					return err
				}
				c.Refs = refs
			case "variables":
				variables, err := utils.ParseYamlStringSequenceToSlice(value)
				if err != nil {
					return err
				}
				c.Variables = variables
			case "changes":
				changes, err := utils.ParseYamlStringSequenceToSlice(value)
				if err != nil {
					return err
				}
				c.Changes = changes
			case "kubernetes":
				c.Kubernetes = value.Value
			}
		}
	}
	return consts.NewErrInvalidYamlTag(node.Tag)
}
