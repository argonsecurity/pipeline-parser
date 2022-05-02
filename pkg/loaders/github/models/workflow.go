package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type Container struct {
	Credentials *Credentials  `mapstructure:"credentials,omitempty" yaml:"credentials,omitempty"`
	Env         interface{}   `mapstructure:"env,omitempty" yaml:"env,omitempty"`
	Image       string        `mapstructure:"image" yaml:"image"`
	Options     string        `mapstructure:"options,omitempty" yaml:"options,omitempty"`
	Ports       []interface{} `mapstructure:"ports,omitempty" yaml:"ports,omitempty"`
	Volumes     []string      `mapstructure:"volumes,omitempty" yaml:"volumes,omitempty"`
}

type Credentials struct {
	Password string `mapstructure:"password,omitempty" yaml:"password,omitempty"`
	Username string `mapstructure:"username,omitempty" yaml:"username,omitempty"`
}

type Defaults struct {
	Run *Run `mapstructure:"run,omitempty" yaml:"run,omitempty"`
}

type Environment struct {
	Name string `mapstructure:"name" yaml:"name"`
	Url  string `mapstructure:"url,omitempty" yaml:"url,omitempty"`
}

type Ref struct {
	Branches       []string `mapstructure:"branches,omitempty" yaml:"branches,omitempty"`
	BranchesIgnore []string `mapstructure:"branches-ignore,omitempty" yaml:"branches-ignore,omitempty"`
	Paths          []string `mapstructure:"paths,omitempty" yaml:"paths,omitempty"`
	PathsIgnore    []string `mapstructure:"paths-ignore,omitempty" yaml:"paths-ignore,omitempty"`
	Tags           []string `mapstructure:"tags,omitempty" yaml:"tags,omitempty"`
	TagsIgnore     []string `mapstructure:"tags-ignore,omitempty" yaml:"tags-ignore,omitempty"`
	FileLocation   *models.FileLocation
}

type Workflow struct {
	Concurrency *Concurrency                 `mapstructure:"concurrency,omitempty" yaml:"concurrency,omitempty"`
	Defaults    *Defaults                    `mapstructure:"defaults,omitempty" yaml:"defaults,omitempty"`
	Env         *models.EnvironmentVariables `mapstructure:"env,omitempty" yaml:"env,omitempty"`
	Jobs        *Jobs                        `mapstructure:"jobs" yaml:"jobs"`
	Name        string                       `mapstructure:"name,omitempty" yaml:"name,omitempty"`
	On          *On                          `mapstructure:"on" yaml:"on"`
	Permissions *PermissionsEvent            `mapstructure:"permissions,omitempty" yaml:"permissions,omitempty"`
}

type Run struct {
	Shell            interface{} `mapstructure:"shell,omitempty" yaml:"shell,omitempty"`
	WorkingDirectory string      `mapstructure:"working-directory,omitempty" yaml:"working-directory,omitempty"`
}

type Strategy struct {
	FailFast    bool        `mapstructure:"fail-fast,omitempty" yaml:"fail-fast,omitempty"`
	Matrix      interface{} `mapstructure:"matrix" yaml:"matrix"`
	MaxParallel float64     `mapstructure:"max-parallel,omitempty" yaml:"max-parallel,omitempty"`
}
