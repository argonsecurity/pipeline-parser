package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"gopkg.in/yaml.v3"
)

type StepMap map[string][]*Step

type Step struct {
	Step      *ExecutionUnitRef `yaml:"step,omitempty"`
	Parallel  []*ParallelSteps  `yaml:"parallel"`
	Variables []*Variable       `yaml:"variables"` // List of variables for the custom pipeline
}

type ParallelSteps struct {
	Step *ExecutionUnitRef `yaml:"step,omitempty"`
}

func (s *Step) UnmarshalYAML(node *yaml.Node) error {
	var step *ExecutionUnitRef
	if err := node.Decode(&step); err != nil {
		return err
	}
	*s = Step{Step: step}
	return nil
}

func (sm *StepMap) UnmarshalYAML(node *yaml.Node) error {
	var stepMap = make(map[string][]*Step)
	return loadersUtils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		var steps []*Step
		if err := loadersUtils.ParseSequenceOrOne(value, &steps); err != nil {
			return err
		}
		stepMap[key] = steps
		*sm = stepMap
		return nil
	}, "StepMap")
}
