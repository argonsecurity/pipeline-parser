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
	Batch    bool    `yaml:"batch,omitempty"`
	Branches *Filter `yaml:"branches,omitempty"`
	Paths    *Filter `yaml:"paths,omitempty"`
	Tags     *Filter `yaml:"tags,omitempty"`
}

type TriggerRef struct {
	Trigger       *Trigger `yaml:"trigger,omitempty"`
	FileReference *models.FileReference
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
			Branches: &Filter{Include: branches},
		}
		return nil
	}

	tr.FileReference.StartRef.Line--
	tr.FileReference.StartRef.Column -= 2
	return node.Decode(&tr.Trigger)
}

type PR struct {
	AutoCancel bool    `yaml:"autoCancel,omitempty"`
	Branches   *Filter `yaml:"branches,omitempty"`
	Paths      *Filter `yaml:"paths,omitempty"`
	Drafts     bool    `yaml:"drafts,omitempty"`
}

type PRRef struct {
	PR            *PR `yaml:"pr,omitempty"`
	FileReference *models.FileReference
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
			Branches: &Filter{Include: branches},
		}
		return nil
	}

	prr.FileReference.StartRef.Line--
	prr.FileReference.StartRef.Column -= 2
	return node.Decode(&prr.PR)
}

type Cron struct {
	Cron          string  `yaml:"cron,omitempty"`
	DisplayName   string  `yaml:"displayName,omitempty"`
	Branches      *Filter `yaml:"branches,omitempty"`
	Batch         bool    `yaml:"batch,omitempty"`
	Always        bool    `yaml:"always,omitempty"`
	FileReference *models.FileReference
}

type Schedules struct {
	Crons         *[]Cron `yaml:"schedules,omitempty"`
	FileReference *models.FileReference
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
