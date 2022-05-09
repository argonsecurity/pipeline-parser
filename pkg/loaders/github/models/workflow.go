package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type Container struct {
	Credentials *Credentials  `yaml:"credentials,omitempty"`
	Env         interface{}   `yaml:"env,omitempty"`
	Image       string        `yaml:"image"`
	Options     string        `yaml:"options,omitempty"`
	Ports       []interface{} `yaml:"ports,omitempty"`
	Volumes     []string      `yaml:"volumes,omitempty"`
}

type Credentials struct {
	Password string `yaml:"password,omitempty"`
	Username string `yaml:"username,omitempty"`
}

type Defaults struct {
	Run *Run `yaml:"run,omitempty"`
}

type Environment struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url,omitempty"`
}

type Ref struct {
	Branches       []string `yaml:"branches,omitempty"`
	BranchesIgnore []string `yaml:"branches-ignore,omitempty"`
	Paths          []string `yaml:"paths,omitempty"`
	PathsIgnore    []string `yaml:"paths-ignore,omitempty"`
	Tags           []string `yaml:"tags,omitempty"`
	TagsIgnore     []string `yaml:"tags-ignore,omitempty"`
	FileReference  *models.FileReference
}

type Run struct {
	Shell            interface{} `yaml:"shell,omitempty"`
	WorkingDirectory string      `yaml:"working-directory,omitempty"`
}

type Strategy struct {
	FailFast    bool        `yaml:"fail-fast,omitempty"`
	Matrix      interface{} `yaml:"matrix"`
	MaxParallel float64     `yaml:"max-parallel,omitempty"`
}

type Workflow struct {
	Concurrency *Concurrency             `yaml:"concurrency,omitempty"`
	Defaults    *Defaults                `yaml:"defaults,omitempty"`
	Env         *EnvironmentVariablesRef `yaml:"env,omitempty"`
	Jobs        *Jobs                    `yaml:"jobs"`
	Name        string                   `yaml:"name,omitempty"`
	On          *On                      `yaml:"on"`
	Permissions *PermissionsEvent        `yaml:"permissions,omitempty"`
}
