package github

import (
	"gopkg.in/yaml.v3"
)

func Parse(data []byte) (*Workflow, error) {
	workflow := &Workflow{}
	if err := yaml.Unmarshal(data, workflow); err != nil {
		return nil, err
	}

	trigger, err := parseWorkflowTriggers(workflow)
	if err != nil {
		return nil, err
	}

	println(trigger)
	return workflow, nil
}
