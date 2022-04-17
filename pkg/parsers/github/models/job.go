package models

import "github.com/argonsecurity/pipeline-parser/pkg/utils"

type Jobs struct {
	NormalJobs               map[string]*Job
	ReusableWorkflowCallJobs map[string]*ReusableWorkflowCallJob
}

type Concurrency struct {
	CancelInProgress *bool   `mapstructure:"cancel-in-progress,omitempty" yaml:"cancel-in-progress,omitempty"`
	Group            *string `mapstructure:"group" yaml:"group"`
}

func (c *Concurrency) UnmarshalText(text []byte) error {
	c.Group = utils.GetPtr(string(text))
	return nil
}

type Job struct {
	Concurrency     *Concurrency          `mapstructure:"concurrency,omitempty" yaml:"concurrency,omitempty"`
	Container       interface{}           `mapstructure:"container,omitempty" yaml:"container,omitempty"`
	ContinueOnError bool                  `mapstructure:"continue-on-error,omitempty" yaml:"continue-on-error,omitempty"`
	Defaults        *Defaults             `mapstructure:"defaults,omitempty" yaml:"defaults,omitempty"`
	Env             interface{}           `mapstructure:"env,omitempty" yaml:"env,omitempty"`
	Environment     interface{}           `mapstructure:"environment,omitempty" yaml:"environment,omitempty"`
	If              string                `mapstructure:"if,omitempty" yaml:"if,omitempty"`
	Name            string                `mapstructure:"name,omitempty" yaml:"name,omitempty"`
	Needs           interface{}           `mapstructure:"needs,omitempty" yaml:"needs,omitempty"`
	Outputs         map[string]string     `mapstructure:"outputs,omitempty" yaml:"outputs,omitempty"`
	Permissions     *PermissionsEvent     `mapstructure:"permissions,omitempty" yaml:"permissions,omitempty"`
	RunsOn          *RunsOn               `mapstructure:"runs-on" yaml:"runs-on"`
	Services        map[string]*Container `mapstructure:"services,omitempty" yaml:"services,omitempty"`
	Steps           *[]Step               `mapstructure:"steps,omitempty" yaml:"steps,omitempty"`
	Strategy        *Strategy             `mapstructure:"strategy,omitempty" yaml:"strategy,omitempty"`
	TimeoutMinutes  *float64              `mapstructure:"timeout-minutes,omitempty" yaml:"timeout-minutes,omitempty"`
}

type ReusableWorkflowCallJob struct {
	If          string            `mapstructure:"if,omitempty" yaml:"if,omitempty"`
	Name        string            `mapstructure:"name,omitempty" yaml:"name,omitempty"`
	Needs       interface{}       `mapstructure:"needs,omitempty" yaml:"needs,omitempty"`
	Permissions *PermissionsEvent `mapstructure:"permissions,omitempty" yaml:"permissions,omitempty"`
	Secrets     interface{}       `mapstructure:"secrets,omitempty" yaml:"secrets,omitempty"`
	Uses        string            `mapstructure:"uses" yaml:"uses"`
	With        interface{}       `mapstructure:"with,omitempty"`
}
