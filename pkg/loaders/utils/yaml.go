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
		EndRef: GetEndFileRef(node),
	}
}

func GetEndFileRef(node *yaml.Node) *models.FileRef {
	if node.Content == nil {
		return &models.FileRef{
			Line:   node.Line,
			Column: node.Column,
		}
	}

	return GetEndFileRef(node.Content[len(node.Content)-1])
}

func GetMapKeyFileLocation(jobIDNode, jobNode *yaml.Node) *models.FileLocation {
	return &models.FileLocation{
		StartRef: &models.FileRef{
			Line:   jobIDNode.Line,
			Column: jobIDNode.Column,
		},
		EndRef: GetEndFileRef(jobNode),
	}
}
