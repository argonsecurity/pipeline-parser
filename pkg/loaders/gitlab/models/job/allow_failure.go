package job

import (
	"errors"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	commonUtils "github.com/argonsecurity/pipeline-parser/pkg/utils"
	"gopkg.in/yaml.v3"
)

type AllowFailure struct {
	Enabled   *bool
	ExitCodes []int `yaml:"exit_codes,omitempty"`
}

func (a *AllowFailure) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == consts.BooleanTag {
		a.Enabled = utils.MustParseYamlBooleanValue(node)
		return nil
	}

	if node.Tag == consts.MapTag {
		a.Enabled = commonUtils.GetPtr(true)
		if len(node.Content) != 2 {
			return errors.New("AllowFailure node must have exactly one key and value")
		}

		exitCodesNode := node.Content[1]
		if exitCodesNode.Tag == consts.IntTag {
			var exitCode int
			if err := exitCodesNode.Decode(&exitCode); err != nil {
				return err
			}
			a.ExitCodes = []int{exitCode}
			return nil
		}

		if exitCodesNode.Tag == consts.SequenceTag {
			var exitCodes []int
			if err := exitCodesNode.Decode(&exitCodes); err != nil {
				return err
			}
			return nil
		}
		return consts.NewErrInvalidYamlTag(exitCodesNode.Tag, "ExitCode")
	}
	return consts.NewErrInvalidYamlTag(node.Tag, "AllowFailure")
}
