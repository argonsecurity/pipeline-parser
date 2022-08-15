package job

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
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

	return utils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		switch key {
		case "include":
			t.Include = value.Value
		case "strategy":
			t.Strategy = value.Value
		case "forward":
			value.Decode(&t.Forward)
		}
		return nil
	}, "Trigger")
}
