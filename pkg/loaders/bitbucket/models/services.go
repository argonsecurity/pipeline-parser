package models

type Service struct {
	Image     *Image                  `yaml:"image"`
	Memory    *int64                  `yaml:"memory,omitempty"`    // Memory limit for the service container, in megabytes
	Variables *EnvironmentVariablesRef `yaml:"variables,omitempty"` // Environment variables passed to the service container
}
