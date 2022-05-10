package job

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"gopkg.in/yaml.v3"
)

type Trigger struct {
	Include  string
	Strategy string
	Forward  *TriggerForward
}

type TriggerForward struct {
	YAMLVariables     bool
	PipelineVariables bool
}

func (t *Trigger) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == consts.StringTag {
		t.Include = node.Value
		return nil
	}

	if node.Tag == consts.MapTag {
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i].Value
			value := node.Content[i+1]
			switch key {
			case "include":
				t.Include = value.Value
			case "strategy":
				t.Strategy = value.Value
			case "forward":
				value.Decode(&t.Forward)
			}
		}
		return nil
	}

	return consts.NewErrInvalidYamlTag(node.Tag)
}
