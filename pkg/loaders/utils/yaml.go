package utils

import (
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

func GetFileLocation(node *yaml.Node) *models.FileLocation {
	return &models.FileLocation{
		StartRef: &models.FileRef{
			Line:   node.Line,
			Column: node.Column,
		},
		EndRef: getEndFileRef(node),
	}
}

func getEndFileRef(node *yaml.Node) *models.FileRef {
	if node.Content == nil {
		return &models.FileRef{
			Line:   node.Line,
			Column: node.Column,
		}
	}

	return getEndFileRef(node.Content[len(node.Content)-1])
}
