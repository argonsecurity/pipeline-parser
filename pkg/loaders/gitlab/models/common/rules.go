package common

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

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
	rules := make([]*Rule, len(node.Content))
	for i, item := range node.Content {
		rule := &Rule{}
		if err := item.Decode(&rule); err != nil {
			return err
		}
		rule.FileReference = utils.GetFileReference(item)
		rules[i] = rule
	}

	*r = Rules{
		RulesList:     rules,
		FileReference: utils.GetFileReference(node),
	}
	return nil
}
