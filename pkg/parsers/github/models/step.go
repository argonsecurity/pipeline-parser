package models

type Step struct {
	ContinueOnError  *bool          `mapstructure:"continue-on-error,omitempty" yaml:"continue-on-error,omitempty"`
	Env              any            `mapstructure:"env,omitempty" yaml:"env,omitempty"`
	Id               string         `mapstructure:"id,omitempty" yaml:"id,omitempty"`
	If               string         `mapstructure:"if,omitempty" yaml:"if,omitempty"`
	Name             string         `mapstructure:"name,omitempty" yaml:"name,omitempty"`
	Run              string         `mapstructure:"run,omitempty" yaml:"run,omitempty"`
	Shell            any            `mapstructure:"shell,omitempty" yaml:"shell,omitempty"`
	TimeoutMinutes   int            `mapstructure:"timeout-minutes,omitempty" yaml:"timeout-minutes,omitempty"`
	Uses             string         `mapstructure:"uses,omitempty" yaml:"uses,omitempty"`
	With             map[string]any `mapstructure:"with,omitempty" yaml:"with,omitempty"`
	WorkingDirectory string         `mapstructure:"working-directory,omitempty" yaml:"working-directory,omitempty"`
}
