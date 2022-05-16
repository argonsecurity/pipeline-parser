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

type Rule struct {
	Changes       []string                 `yaml:"changes,omitempty"`
	Exists        []string                 `yaml:"exists,omitempty"`
	If            string                   `yaml:"if,omitempty"`
	Variables     *EnvironmentVariablesRef `yaml:"variables,omitempty"`
	When          string                   `yaml:"when,omitempty"`
	FileReference *models.FileReference
}

type Rules struct {
	RulesList     []*Rule
	FileReference *models.FileReference
}

func (r *Rules) UnmarshalYAML(node *yaml.Node) error {
	rules := make([]*Rule, 0)
	for _, item := range node.Content {
		rule := &Rule{}
		if err := item.Decode(&rule); err != nil {
			return err
		}
		rule.FileReference = utils.GetFileReference(item)
		rules = append(rules, rule)
	}

	*r = Rules{
		RulesList:     rules,
		FileReference: utils.GetFileReference(node),
	}
	return nil
}
