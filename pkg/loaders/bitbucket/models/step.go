package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"gopkg.in/yaml.v3"
)

type StepMap map[string][]*Step

type Step struct {
	Step      *ExecutionUnitRef     `yaml:"step,omitempty"`
	Parallel  []*ParallelSteps      `yaml:"parallel"`
	Variables []*CustomStepVariable `yaml:"variables"` // List of variables for the custom pipeline
}

type ParallelSteps struct {
	Step *ExecutionUnitRef `yaml:"step,omitempty"`
}

func (s *Step) UnmarshalYAML(node *yaml.Node) error {
	return loadersUtils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		switch key {
		case "parallel":
			var parallel []*ParallelSteps
			if err := loadersUtils.ParseSequenceOrOne(value, &parallel); err != nil {
				return err
			}
			s.Parallel = parallel
			return nil
		case "step":
			var step *ExecutionUnitRef
			if err := value.Decode(&step); err != nil {
				return err
			}
			if isAliasStep(value) {
				step.FileReference.IsAlias = true
			}
			s.Step = step
			return nil
		case "variables":
			vars, err := parseVariables(value)
			if err != nil {
				return err
			}
			s.Variables = vars
			return nil
		}
		return nil
	}, "Step")
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

func parseVariables(node *yaml.Node) ([]*CustomStepVariable, error) {
	var vars []*CustomStepVariable
	if err := loadersUtils.ParseSequenceOrOne(node, &vars); err != nil {
		return nil, err
	}
	return vars, nil
}

func isAliasStep(node *yaml.Node) bool {
	if node.Kind == yaml.AliasNode {
		return true
	}

	var res = false
	loadersUtils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		if value.Kind == yaml.AliasNode {
			res = true
		}
		return nil
	}, "Step")
	return res
}
