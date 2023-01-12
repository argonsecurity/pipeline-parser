package models

type Pipeline struct {
	Image       *Image          `yaml:"image"`
	Clone       *Clone          `yaml:"clone,omitempty"`
	Options     *GlobalSettings `yaml:"options,omitempty"`
	Definitions *Definitions    `yaml:"definitions,omitempty"`
	Pipelines   *BuildPipelines `yaml:"pipelines"`
}
