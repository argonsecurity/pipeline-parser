package models

type GlobalSettings struct {
	Docker  *bool  `yaml:"docker,omitempty"` // A flag to add Docker to all build steps in all pipelines
	MaxTime *int64 `yaml:"max-time,omitempty"`
	Size    *Size  `yaml:"size,omitempty"`
}

type Size string

const (
	X1 Size = "1x"
	X2 Size = "2x"
)
