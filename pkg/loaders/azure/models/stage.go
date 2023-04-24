package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Stage struct {
	Stage           string         `yaml:"stage,omitempty"`
	DisplayName     string         `yaml:"displayName,omitempty"`
	Pool            *Pool          `yaml:"pool,omitempty"`
	DependsOn       *DependsOn     `yaml:"dependsOn,omitempty"`
	Condition       string         `yaml:"condition,omitempty"`
	Variables       *Variables     `yaml:"variables,omitempty"`
	Jobs            *Jobs          `yaml:"jobs,omitempty"`
	LockBehavior    string         `yaml:"lockBehavior,omitempty"`
	TemplateContext map[string]any `yaml:"templateContext,omitempty"`
	FileReference   *models.FileReference
}

type TemplateStage struct {
	Template      `yaml:",inline"`
	FileReference *models.FileReference
}

type Stages struct {
	Stages         []*Stage         `yaml:"stages,omitempty"`
	TemplateStages []*TemplateStage `yaml:"templateStages,omitempty"`
	FileReference  *models.FileReference
}

func (s *Stages) UnmarshalYAML(node *yaml.Node) error {
	var stages []*Stage
	var templateStages []*TemplateStage

	for _, stageNode := range node.Content {
		if stageNode.Tag == consts.StringTag {
			templateStages = append(templateStages, &TemplateStage{
				Template: Template{
					Template: stageNode.Value,
				},
				FileReference: loadersUtils.GetFileReference(stageNode),
			})
			continue
		}

		if isTemplateStage(stageNode) {
			stage, err := parseTemplateStage(stageNode)
			if err != nil {
				return err
			}
			templateStages = append(templateStages, &stage)
			continue
		}

		stage, err := parseStage(stageNode)
		if err != nil {
			return err
		}
		stages = append(stages, &stage)
	}

	*s = Stages{
		Stages:         stages,
		TemplateStages: templateStages,
		FileReference:  loadersUtils.GetFileReference(node),
	}
	return nil
}

func parseStage(node *yaml.Node) (Stage, error) {
	var stage Stage
	if err := node.Decode(&stage); err != nil {
		return stage, err
	}
	stage.FileReference = loadersUtils.GetFileReference(node)
	return stage, nil
}

func parseTemplateStage(node *yaml.Node) (TemplateStage, error) {
	var templateStage TemplateStage
	if err := node.Decode(&templateStage); err != nil {
		return templateStage, err
	}
	templateStage.FileReference = loadersUtils.GetFileReference(node)
	return templateStage, nil
}

func isTemplateStage(job *yaml.Node) bool {
	for _, node := range job.Content {
		if node.Tag == consts.StringTag && node.Value == "template" {
			return true
		}
	}
	return false
}
