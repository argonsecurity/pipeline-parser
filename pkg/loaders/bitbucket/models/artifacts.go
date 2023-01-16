package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"gopkg.in/yaml.v3"
)

type Artifacts struct {
	SharedStepFiles *SharedStepFiles
	Paths           []*string
}

type SharedStepFiles struct {
	Download *bool    `yaml:"download,omitempty"` // Indicates whether to download artifact in the step
	Paths    []*string `yaml:"paths"`
}

func (a *Artifacts) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.SequenceNode {
		var paths []*string
		if err := loadersUtils.ParseSequenceOrOne(node, &paths); err != nil {
			return err
		}
		a.Paths = paths
		return nil
	}
	if node.Kind == yaml.MappingNode {
		var sharedStepFiles SharedStepFiles
		if err := node.Decode(&sharedStepFiles); err != nil {
			return err
		}
		a.SharedStepFiles = &sharedStepFiles
		return nil
	}
	return nil
}
