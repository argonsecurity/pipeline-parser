package models

import (
	commonModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/common/models"
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Steps []Step

type ShellCommand struct {
	Script        string
	FileReference *models.FileReference
}

type With commonModels.Map

type Step struct {
	ContinueOnError  *string                  `yaml:"continue-on-error,omitempty"`
	Env              *EnvironmentVariablesRef `yaml:"env,omitempty"`
	Id               string                   `yaml:"id,omitempty"`
	If               string                   `yaml:"if,omitempty"`
	Name             string                   `yaml:"name,omitempty"`
	Run              *ShellCommand            `yaml:"run,omitempty"`
	Shell            string                   `yaml:"shell,omitempty"`
	TimeoutMinutes   int                      `yaml:"timeout-minutes,omitempty"`
	Uses             string                   `yaml:"uses,omitempty"`
	With             *With                    `yaml:"with,omitempty"`
	WorkingDirectory string                   `yaml:"working-directory,omitempty"`
	FileReference    *models.FileReference
}

func (w *With) UnmarshalYAML(value *yaml.Node) error {
	var m commonModels.Map
	if err := value.Decode(&m); err != nil {
		return err
	}
	*w = With(m)
	return nil
}

func (s *Steps) UnmarshalYAML(node *yaml.Node) error {
	steps := []Step{}
	for _, stepNode := range node.Content {
		var step Step
		if err := stepNode.Decode(&step); err != nil {
			return err
		}
		step.FileReference = loadersUtils.GetFileReference(stepNode)
		steps = append(steps, step)
	}
	*s = steps
	return nil
}

func (s *ShellCommand) UnmarshalYAML(node *yaml.Node) error {
	s.FileReference = loadersUtils.GetFileReference(node)
	s.Script = node.Value
	return nil
}
