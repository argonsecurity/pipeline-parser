package utils

import (
	"fmt"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

func GetFileReference(node *yaml.Node) *models.FileReference {
	return &models.FileReference{
		StartRef: &models.FileLocation{
			Line:   node.Line,
			Column: node.Column,
		},
		EndRef: GetEndFileLocation(node),
	}
}

func GetEndFileLocation(node *yaml.Node) *models.FileLocation {
	if node.Content == nil {
		return &models.FileLocation{
			Line:   node.Line,
			Column: node.Column,
		}
	}

	return GetEndFileLocation(node.Content[len(node.Content)-1])
}

func GetMapKeyFileReference(jobIDNode, jobNode *yaml.Node) *models.FileReference {
	return &models.FileReference{
		StartRef: &models.FileLocation{
			Line:   jobIDNode.Line,
			Column: jobIDNode.Column,
		},
		EndRef: GetEndFileLocation(jobNode),
	}
}

func ParseYamlStringSequenceToSlice(node *yaml.Node) ([]string, error) {
	if node.Tag != consts.SequenceTag {
		return nil, fmt.Errorf("expected sequence tag, got %s", node.Tag)
	}

	strings := make([]string, len(node.Content))
	for i, n := range node.Content {
		if n.Tag != consts.StringTag {
			return nil, fmt.Errorf("expected string tag, got %s", n.Tag)
		}

		strings[i] = n.Value
	}
	return strings, nil
}
