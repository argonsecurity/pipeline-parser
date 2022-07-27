package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Filter struct {
	Include []string `yaml:"include,omitempty"`
	Exclude []string `yaml:"exclude,omitempty"`
}

type Trigger struct {
	Batch    bool   `yaml:"batch,omitempty"`
	Branches Filter `yaml:"branches,omitempty"`
	Paths    Filter `yaml:"paths,omitempty"`
	Tags     Filter `yaml:"tags,omitempty"`
	Stages   Filter `yaml:"stages,omitempty"`
}

type TriggerRef struct {
	Trigger       *Trigger `yaml:"trigger,omitempty"`
	FileReference *models.FileReference
}

type PR struct {
	AutoCancel bool   `yaml:"autoCancel,omitempty"`
	Branches   Filter `yaml:"branches,omitempty"`
	Paths      Filter `yaml:"paths,omitempty"`
	Drafts     bool   `yaml:"drafts,omitempty"`
}

type PRRef struct {
	PR            *PR `yaml:"pr,omitempty"`
	FileReference *models.FileReference
}

type Cron struct {
	Cron          string `yaml:"cron,omitempty"`
	DisplayName   string `yaml:"displayName,omitempty"`
	Branches      Filter `yaml:"branches,omitempty"`
	Batch         bool   `yaml:"batch,omitempty"`
	Always        bool   `yaml:"always,omitempty"`
	FileReference *models.FileReference
}

type Schedules struct {
	Crons         *[]Cron `yaml:"schedules,omitempty"`
	FileReference *models.FileReference
}

func (f *Filter) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == consts.SequenceTag {
		include, err := loadersUtils.ParseYamlStringSequenceToSlice(node)
		if err != nil {
			return err
		}
		f.Include = include
		return nil
	}

	if node.Tag == consts.MapTag {
		return loadersUtils.IterateOnMap(node, func(key string, value *yaml.Node) error {
			parsedValue, err := loadersUtils.ParseYamlStringSequenceToSlice(value)
			if err != nil {
				return err
			}

			if key == "include" {
				f.Include = parsedValue
			} else if key == "exclude" {
				f.Exclude = parsedValue
			}
			return nil
		})
	}

	return consts.NewErrInvalidYamlTag(node.Tag, "Filter")
}

func (tr *TriggerRef) UnmarshalYAML(node *yaml.Node) error {
	tr.FileReference = loadersUtils.GetFileReference(node)
	if node.Tag == consts.StringTag {
		tr.Trigger = nil
		return nil
	}

	if node.Tag == consts.SequenceTag {
		branches, err := loadersUtils.ParseYamlStringSequenceToSlice(node)
		if err != nil {
			return err
		}
		tr.Trigger = &Trigger{
			Branches: Filter{Include: branches},
		}
		return nil
	}

	tr.FileReference.StartRef.Line--      // The "trigger" node is not accessible, this is a patch
	tr.FileReference.StartRef.Column -= 2 // The "trigger" node is not accessible, this is a patch
	return node.Decode(&tr.Trigger)
}

func (prr *PRRef) UnmarshalYAML(node *yaml.Node) error {
	prr.FileReference = loadersUtils.GetFileReference(node)
	if node.Tag == consts.StringTag {
		prr.PR = nil
		return nil
	}

	if node.Tag == consts.SequenceTag {
		branches, err := loadersUtils.ParseYamlStringSequenceToSlice(node)
		if err != nil {
			return err
		}
		prr.PR = &PR{
			Branches: Filter{Include: branches},
		}
		return nil
	}

	prr.FileReference.StartRef.Line--
	prr.FileReference.StartRef.Column -= 2
	return node.Decode(&prr.PR)
}

func (s *Schedules) UnmarshalYAML(node *yaml.Node) error {
	s.FileReference = loadersUtils.GetFileReference(node)
	crons := []Cron{}
	for _, cronNode := range node.Content {
		var cron Cron
		if err := cronNode.Decode(&cron); err != nil {
			return err
		}
		cron.FileReference = loadersUtils.GetFileReference(cronNode)
		crons = append(crons, cron)
	}

	s.Crons = &crons
	return nil
}
