package models

type Artifacts struct {
	SharedStepFiles *SharedStepFiles
	Paths           []string
}

type SharedStepFiles struct {
	Download *bool    `yaml:"download,omitempty"` // Indicates whether to download artifact in the step
	Paths    []string `yaml:"paths"`
}
