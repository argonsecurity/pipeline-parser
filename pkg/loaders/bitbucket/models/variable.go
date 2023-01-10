package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Variable struct {
	Name          *string `json:"name,omitempty"` // Name of a variable for the custom pipeline
	FileReference *models.FileReference
}

func (v *Variable) UnmarshalYAML(node *yaml.Node) error {
	v.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&v.Name)
}
