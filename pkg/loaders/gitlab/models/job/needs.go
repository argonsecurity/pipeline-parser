package job

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"gopkg.in/yaml.v3"
)

type NeedsItem struct {
	Artifacts bool   `yaml:"artifacts"`
	Project   string `yaml:"project"`
	Pipeline  string `yaml:"pipeline"`
	Job       string `yaml:"job"`
	Ref       string `yaml:"ref"`
}

type Needs []*NeedsItem

func (n *Needs) UnmarshalYAML(node *yaml.Node) error {
	needs := []*NeedsItem{}
	for _, item := range node.Content {
		if item.Tag == consts.StringTag {
			needs = append(needs,
				&NeedsItem{
					Job: item.Value,
				},
			)
		} else if item.Tag == consts.MapTag {
			needsItem := &NeedsItem{}
			if err := item.Decode(&needsItem); err != nil {
				return err
			}
			needs = append(needs, needsItem)
		}
	}

	*n = needs
	return nil
}
