package common

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Cache struct {
	When      string   `yaml:"when,omitempty"`
	Key       any      `yaml:"key,omitempty"`
	Paths     []string `yaml:"paths,omitempty"`
	Policy    string   `yaml:"policy,omitempty"`
	Untracked bool     `yaml:"untracked,omitempty"`
}

type RulesItems struct {
	Changes       []string                 `yaml:"changes,omitempty"`
	Exists        []string                 `yaml:"exists,omitempty"`
	If            string                   `yaml:"if,omitempty"`
	Variables     *EnvironmentVariablesRef `yaml:"variables,omitempty"`
	When          string                   `yaml:"when,omitempty"`
	FileReference *models.FileReference
}

type Rules []*RulesItems

func (r *Rules) UnmarshalYAML(node *yaml.Node) error {
	for _, item := range node.Content {
		rule := &RulesItems{}
		if err := item.Decode(&rule); err != nil {
			return err
		}
		rule.FileReference = utils.GetFileReference(item)
		*r = append(*r, rule)
	}
	return nil
}
