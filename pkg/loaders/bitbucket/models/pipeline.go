package models

type Pipeline struct {
	Clone       *Clone          `yaml:"clone,omitempty"`
	Definitions *Definitions    `yaml:"definitions,omitempty"`
	Image       *Image          `yaml:"image"`
	Options     *GlobalSettings `yaml:"options,omitempty"`
	Pipelines   *BuildPipelines `yaml:"pipelines"`
}
