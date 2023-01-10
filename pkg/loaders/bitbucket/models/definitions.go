package models

type Definitions struct {
	Caches   *Caches             `yaml:"caches,omitempty"`
	Services map[string]Services `yaml:"services,omitempty"`
	Steps    []Step              `yaml:"steps,omitempty"`
}
