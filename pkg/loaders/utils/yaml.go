package utils

import (
	"fmt"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
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
