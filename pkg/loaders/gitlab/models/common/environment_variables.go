package common

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type EnvironmentVariablesRef struct {
	Variables     map[string]any
	FileReference *models.FileReference
}

func (e *EnvironmentVariablesRef) UnmarshalYAML(node *yaml.Node) error {
	e.FileReference = utils.GetFileReference(node)

	e.FileReference.StartRef.Line--
	e.FileReference.StartRef.Column = 0
	return node.Decode(&e.Variables)
}
