package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type EnvironmentVariablesRef struct {
	models.EnvironmentVariables
	FileReference *models.FileReference
}

type Template struct {
	Template   string         `yaml:"template,omitempty"`
	Parameters map[string]any `yaml:"parameters,omitempty"`
}

type Extends struct {
	Template      `yaml:"inline"`
	FileReference *models.FileReference
}

type DependsOn []string

func (e *EnvironmentVariablesRef) UnmarshalYAML(node *yaml.Node) error {
	var env models.EnvironmentVariables
	if err := node.Decode(&env); err != nil {
		return err
	}

	e.EnvironmentVariables = env
	e.FileReference = loadersUtils.GetFileReference(node)
	e.FileReference.StartRef.Line-- // The "env" node is not accessible, this is a patch
	return nil
}

func (e *Extends) UnmarshalYAML(node *yaml.Node) error {
	e.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&e.Template)
}

func (n *DependsOn) UnmarshalYAML(node *yaml.Node) error {
	var tags []string
	if err := loadersUtils.ParseSequenceOrOne(node, &tags); err != nil {
		return err
	}

	*n = tags
	return nil
}
