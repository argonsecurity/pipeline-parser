package models

import "gopkg.in/yaml.v3"

type Jobs []Job

type Matrix map[string]map[string]string

type Strategy struct {
	Matrix      *Matrix `yaml:"matrix,omitempty"`
	MaxParallel int     `yaml:"maxParallel,omitempty"`
	Parallel    string  `yaml:"parallel,omitempty"`
}

// type JobContainer struct {
// 	Image           string                   `yaml:"image,omitempty"`
// 	Endpoint        string                   `yaml:"endpoint,omitempty"`
// 	Env             *EnvironmentVariablesRef `yaml:"env,omitempty"`
// 	MapDockerSocket bool                     `yaml:"mapDockerSocket,omitempty"`
// 	Options         string                 `yaml:"options,omitempty"`
// 	Ports		   []string                 `yaml:"ports,omitempty"`
// 	Volumes		   []string                 `yaml:"volumes,omitempty"`
// 	MountReadOnly   *MountReadOnly           `yaml:"mountReadOnly,omitempty"`

// }

type Job struct {
	Job                    string     `yaml:"job,omitempty"`
	DisplayName            string     `yaml:"displayName,omitempty"`
	DependsOn              []string   `yaml:"dependsOn,omitempty"`
	Condition              string     `yaml:"condition,omitempty"`
	ContinueOnError        bool       `yaml:"continueOnError,omitempty"`
	TimeoutInMinutes       int        `yaml:"timeoutInMinutes,omitempty"`
	CancelTimeoutInMinutes int        `yaml:"cancelTimeoutInMinutes,omitempty"`
	Variables              *Variables `yaml:"variables,omitempty"`
	Strategy               *Strategy  `yaml:"strategy,omitempty"`
	Pool                   *Pool      `yaml:"pool,omitempty"`
}

func (j *Jobs) UnmarshalYAML(node *yaml.Node) error {
	return nil
}
