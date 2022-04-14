package models

type Jobs struct {
	NormalJobs               map[string]*NormalJob
	ReusableWorkflowCallJobs map[string]*ReusableWorkflowCallJob
}

type Concurrency struct {
	CancelInProgress *bool   `mapstructure:"cancel-in-progress,omitempty" yaml:"cancel-in-progress,omitempty"`
	Group            *string `mapstructure:"group" yaml:"group"`
}

func (c *Concurrency) UnmarshalText(text []byte) error {
	s := string(text)
	c.Group = &s
	return nil
}

type NormalJob struct {
	Concurrency     *Concurrency          `mapstructure:"concurrency,omitempty" yaml:"concurrency,omitempty" json:"concurrency,omitempty"`
	Container       interface{}           `mapstructure:"container,omitempty" yaml:"container,omitempty" json:"container,omitempty"`
	ContinueOnError bool                  `mapstructure:"continue-on-error,omitempty" yaml:"continue-on-error,omitempty" json:"continue-on-error,omitempty"`
	Defaults        *Defaults             `mapstructure:"defaults,omitempty" yaml:"defaults,omitempty" json:"defaults,omitempty"`
	Env             interface{}           `mapstructure:"env,omitempty" yaml:"env,omitempty" json:"env,omitempty"`
	Environment     interface{}           `mapstructure:"environment,omitempty" yaml:"environment,omitempty" json:"environment,omitempty"`
	If              string                `mapstructure:"if,omitempty" yaml:"if,omitempty" json:"if,omitempty"`
	Name            string                `mapstructure:"name,omitempty" yaml:"name,omitempty" json:"name,omitempty"`
	Needs           interface{}           `mapstructure:"needs,omitempty" yaml:"needs,omitempty" json:"needs,omitempty"`
	Outputs         map[string]string     `mapstructure:"outputs,omitempty" yaml:"outputs,omitempty" json:"outputs,omitempty"`
	Permissions     *PermissionsEvent     `mapstructure:"permissions,omitempty" yaml:"permissions,omitempty" json:"permissions,omitempty"`
	RunsOn          interface{}           `mapstructure:"runs-on" yaml:"runs-on" json:"runs-on"`
	Services        map[string]*Container `mapstructure:"services,omitempty" yaml:"services,omitempty" json:"services,omitempty"`
	Steps           *[]Step               `mapstructure:"steps,omitempty" yaml:"steps,omitempty" json:"steps,omitempty"`
	Strategy        *Strategy             `mapstructure:"strategy,omitempty" yaml:"strategy,omitempty" json:"strategy,omitempty"`
	TimeoutMinutes  *float64              `mapstructure:"timeout-minutes,omitempty" yaml:"timeout-minutes,omitempty" json:"timeout-minutes,omitempty"`
}

type ReusableWorkflowCallJob struct {
	If          string            `mapstructure:"if,omitempty" yaml:"if,omitempty" json:"if,omitempty"`
	Name        string            `mapstructure:"name,omitempty" yaml:"name,omitempty" json:"name,omitempty"`
	Needs       interface{}       `mapstructure:"needs,omitempty" yaml:"needs,omitempty" json:"needs,omitempty"`
	Permissions *PermissionsEvent `mapstructure:"permissions,omitempty" yaml:"permissions,omitempty" json:"permissions,omitempty"`
	Secrets     interface{}       `mapstructure:"secrets,omitempty" yaml:"secrets,omitempty" json:"secrets,omitempty"`
	Uses        string            `mapstructure:"uses" yaml:"uses" json:"uses"`
	With        interface{}       `json:"with,omitempty"`
}
