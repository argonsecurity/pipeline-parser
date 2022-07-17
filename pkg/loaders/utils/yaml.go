package utils

import (
	"fmt"
	"strings"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"gopkg.in/yaml.v3"
)

func GetFileReference(node *yaml.Node) *models.FileReference {
	fr := &models.FileReference{
		StartRef: &models.FileLocation{
			Line:   node.Line,
			Column: node.Column,
		},
		EndRef: GetEndFileLocation(node),
	}

	// Sometimes, a sequence node's line and column are equal to the line of the first variable (probably a bug in yaml.v3).
	if node.Tag == consts.SequenceTag && len(node.Content) > 0 && node.Content[0].Line == fr.StartRef.Line {
		if node.Content[0].Column == node.Column+2 { // Making sure that the sequence node's format is not "name: [val1, val2, ...]"
			fr.StartRef.Line--
			fr.StartRef.Column -= 2
		} else {
			fr.EndRef.Column++ // +1 for the "]" after the last value
		}
	}
	return fr
}

func GetEndFileLocation(node *yaml.Node) *models.FileLocation {
	if node.Content == nil {
		return calculateValueNodeEndFileLocation(node)
	}

	return GetEndFileLocation(node.Content[len(node.Content)-1])
}

func calculateValueNodeEndFileLocation(node *yaml.Node) *models.FileLocation {
	split := strings.Split(node.Value, "\n")
	return &models.FileLocation{
		Line:   node.Line + len(split) - 1,
		Column: node.Column + len(split[len(split)-1]),
	}
}

func GetMapKeyFileReference(keyNode, valueNode *yaml.Node) *models.FileReference {
	return &models.FileReference{
		StartRef: &models.FileLocation{
			Line:   keyNode.Line,
			Column: keyNode.Column,
		},
		EndRef: GetEndFileLocation(valueNode),
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

// A Map YAML node is very messy to iterate on
// This function wraps the messy part for cleaner code
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

func ParseSequenceOrOne[T any](node *yaml.Node, v *[]T) error {
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
