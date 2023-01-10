package models

import "gopkg.in/yaml.v3"

type Services struct {
	Image     *Image                  `yaml:"image"`
	Memory    *int64                  `yaml:"memory,omitempty"`    // Memory limit for the service container, in megabytes
	Variables EnvironmentVariablesRef `yaml:"variables,omitempty"` // Environment variables passed to the service container
}

func (s *Services) UnmarshalYAML(node *yaml.Node) error {
	return node.Decode(&s)
}
