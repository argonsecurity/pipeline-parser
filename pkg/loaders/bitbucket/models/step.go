package models

type Step struct {
	Step      *ExecutionUnitRef `yaml:"step,omitempty"`
	Parallel  []*ParallelSteps    `yaml:"parallel"`
	Variables []*Variable         `yaml:"variables"` // List of variables for the custom pipeline
}

type ParallelSteps struct {
	Step *ExecutionUnitRef `yaml:"step,omitempty"`
}
