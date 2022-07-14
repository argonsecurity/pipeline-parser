package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type EnvironmentVariablesRef struct {
	models.EnvironmentVariables
	FileReference *models.FileReference
}

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

type Template struct {
	Template   string         `yaml:"template,omitempty"`
	Parameters map[string]any `yaml:"parameters,omitempty"`
}

type Extends struct {
	Template      `yaml:"inline"`
	FileReference *models.FileReference
}

func (e *Extends) UnmarshalYAML(node *yaml.Node) error {
	e.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(e)
}

type DependsOn []string

func (n *DependsOn) UnmarshalYAML(node *yaml.Node) error {
	var tags []string
	var err error

	if node.Tag == consts.SequenceTag {
		if tags, err = loadersUtils.ParseYamlStringSequenceToSlice(node); err != nil {
			return err
		}
	} else if node.Tag == consts.StringTag {
		tags = []string{node.Value}
	} else {
		return consts.NewErrInvalidYamlTag(node.Tag)
	}

	*n = tags
	return nil
}
