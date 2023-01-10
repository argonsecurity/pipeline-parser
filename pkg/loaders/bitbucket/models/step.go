package models

type Step struct {
	Step      *BuildExecutionUnit `yaml:"step,omitempty"`
	Parallel  []ParallelSteps     `yaml:"parallel"`
	Variables []Variable          `yaml:"variables"` // List of variables for the custom pipeline
}

type ParallelSteps struct {
	Step *BuildExecutionUnit `yaml:"step,omitempty"`
}
