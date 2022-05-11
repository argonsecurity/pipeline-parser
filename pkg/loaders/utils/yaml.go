package utils

import (
	"fmt"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
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
		return nil, consts.NewErrInvalidYamlTag(node.Tag)
	}

	strings := make([]string, len(node.Content))
	for i, n := range node.Content {
		if n.Tag != consts.StringTag {
			return nil, consts.NewErrInvalidYamlTag(node.Tag)
		}

		strings[i] = n.Value
	}
	return strings, nil
}

func MustParseYamlBooleanValue(node *yaml.Node) *bool {
	if node.Value == "true" {
		return utils.GetPtr(true)
	}

	if node.Value == "false" {
		return utils.GetPtr(false)
	}

	panic(fmt.Sprintf("invalid boolean value: %s", node.Value))
}

func IterateOnMap(node *yaml.Node, cb func(key string, value *yaml.Node) error) error {
	if node.Tag != consts.MapTag {
		return consts.NewErrInvalidYamlTag(node.Tag)
	}

	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]
		value.Line = key.Line
		value.Column = key.Column
		if err := cb(key.Value, value); err != nil {
			return err
		}
	}

	return nil
}

func ParseSliceOrOne[T any](node *yaml.Node, v *[]T) error {
	var t T
	if node.Tag == consts.SequenceTag {
		return node.Decode(&v)
	}

	if err := node.Decode(&t); err != nil {
		return err
	}

	*v = []T{t}
	return nil
}
