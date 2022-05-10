package common

type Cache struct {
	When      string      `yaml:"when,omitempty"`
	Key       interface{} `yaml:"key,omitempty"`
	Paths     []string    `yaml:"paths,omitempty"`
	Policy    string      `yaml:"policy,omitempty"`
	Untracked bool        `yaml:"untracked,omitempty"`
}

type RulesItems struct {
	Changes   []string   `yaml:"changes,omitempty"`
	Exists    []string   `yaml:"exists,omitempty"`
	If        string     `yaml:"if,omitempty"`
	Variables *Variables `yaml:"variables,omitempty"`
	When      string     `yaml:"when,omitempty"`
}

type Variables struct {
	AdditionalProperties map[string]interface{} `yaml:"-"`
}
