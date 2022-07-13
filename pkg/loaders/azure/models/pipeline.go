package models

type Pipeline struct {
	Extends *Extends `yaml:"extends,omitempty"`
	// Jobs       *Jobs       `yaml:"jobs,omitempty"`
	continueOnError bool        `yaml:"continueOnError,omitempty"`
	Pool            *Pool       `yaml:"pool,omitempty"`
	Name            string      `yaml:"name,omitempty"`
	Trigger         *TriggerRef `yaml:"trigger,omitempty"`
	Parameters      *Parameters `yaml:"parameters,omitempty"`
	PR              *PRRef      `yaml:"pr,omitempty"`
	Schedules       *Schedules  `yaml:"schedules,omitempty"`
	// Resources *Resources  `yaml:"resources,omitempty"`
	Steps        *Steps     `yaml:"steps,omitempty"`
	Variables    *Variables `yaml:"variables,omitempty"`
	LockBehavior string     `yaml:"lockBehavior,omitempty"`
}
