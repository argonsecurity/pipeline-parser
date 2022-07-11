package models

type Parameters []Parameter

type Parameter struct {
	Name        string   `yaml:"name,omitempty"`
	DisplayName string   `yaml:"displayName,omitempty"`
	Type        string   `yaml:"type,omitempty"`
	Default     any      `yaml:"default,omitempty"`
	Values      []string `yaml:"values,omitempty"`
}
