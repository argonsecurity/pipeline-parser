package job

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"gopkg.in/yaml.v3"
)

type Trigger struct {
	Include  *common.Include
	Strategy string
	Forward  *TriggerForward
}

type TriggerForward struct {
	YAMLVariables     bool
	PipelineVariables bool
}

func (t *Trigger) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == consts.StringTag {
		t.Include = &common.Include{
			common.ParseIncludeString(node),
		}
		return nil
	}

	return utils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		switch key {
		case "include":
			t.Include = &common.Include{}
			value.Decode(t.Include)
		case "strategy":
			t.Strategy = value.Value
		case "forward":
			value.Decode(&t.Forward)
		}
		return nil
	}, "Trigger")
}
