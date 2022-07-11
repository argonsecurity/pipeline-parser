package models

import "github.com/argonsecurity/pipeline-parser/pkg/models"

type Steps []Step

type ShellCommand struct {
	Script        string
	FileReference *models.FileReference
}

type Target struct {
	Container         string   `yaml:"container,omitempty"`
	Commands          string   `yaml:"commands,omitempty"`
	SettableVariables []string `yaml:"settableVariables,omitempty"`
}

type Step struct {
	Name                    string                   `yaml:"name,omitempty"`
	Condition               string                   `yaml:"condition,omitempty"`
	ContinueOnError         bool                     `yaml:"continueOnError,omitempty"`
	DisplayName             string                   `yaml:"displayName,omitempty"`
	Target                  *Target                  `yaml:"target,omitempty"`
	Enabled                 bool                     `yaml:"enabled,omitempty"`
	Env                     *EnvironmentVariablesRef `yaml:"env,omitempty"`
	TimeoutInMinutes        int                      `yaml:"timeoutInMinutes,omitempty"`
	RetryCountOnTaskFailure int                      `yaml:"retryCountOnTaskFailure,omitempty"`
}

type Task struct {
	Step
	Task   string         `yaml:"task,omitempty"`
	Inputs map[string]any `yaml:"inputs,omitempty"`
}
