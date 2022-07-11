package models

type Pipeline struct {
	Name       string      `yaml:"name,omitempty"`
	Pool       *Pool       `yaml:"pool,omitempty"`
	Parameters *Parameters `yaml:"parameters,omitempty"`
	Trigger    *TriggerRef `yaml:"trigger,omitempty"`
	PR         *PRRef      `yaml:"pr,omitempty"`
	Variables  *Variables  `yaml:"variables,omitempty"`
}
