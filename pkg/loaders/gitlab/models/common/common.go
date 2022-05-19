package common

type Cache struct {
	When      string   `yaml:"when,omitempty"`
	Key       any      `yaml:"key,omitempty"`
	Paths     []string `yaml:"paths,omitempty"`
	Policy    string   `yaml:"policy,omitempty"`
	Untracked bool     `yaml:"untracked,omitempty"`
}
