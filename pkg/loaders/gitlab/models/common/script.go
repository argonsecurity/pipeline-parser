package common

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Script struct {
	Commands      []string
	FileReference *models.FileReference
}

func (s *Script) UnmarshalYAML(node *yaml.Node) error {
	s.FileReference = utils.GetFileReference(node)

	if node.Tag == consts.StringTag {
		s.Commands = []string{node.Value}
		s.FileReference.EndRef.Column += len(node.Value)
		return nil
	}

	if node.Tag == consts.SequenceTag {
		// We don't have access to the key in the YAML, this is a patch
		s.FileReference.StartRef.Line--

		commands, err := utils.ParseYamlStringSequenceToSlice(node)
		if err != nil {
			return err
		}
		s.Commands = commands
		return nil
	}

	return consts.NewErrInvalidYamlTag(node.Tag)
}
