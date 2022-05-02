package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Steps []Step

type Step struct {
	ContinueOnError  *bool                        `yaml:"continue-on-error,omitempty"`
	Env              *models.EnvironmentVariables `yaml:"env,omitempty"`
	Id               string                       `yaml:"id,omitempty"`
	If               string                       `yaml:"if,omitempty"`
	Name             string                       `yaml:"name,omitempty"`
	Run              string                       `yaml:"run,omitempty"`
	Shell            string                       `yaml:"shell,omitempty"`
	TimeoutMinutes   int                          `yaml:"timeout-minutes,omitempty"`
	Uses             string                       `yaml:"uses,omitempty"`
	With             map[string]any               `yaml:"with,omitempty"`
	WorkingDirectory string                       `yaml:"working-directory,omitempty"`
	FileLocation     *models.FileLocation
}

func (s *Steps) UnmarshalYAML(node *yaml.Node) error {
	steps := []Step{}
	for _, stepNode := range node.Content {
		var step Step
		if err := stepNode.Decode(&step); err != nil {
			return err
		}
		step.FileLocation = loadersUtils.GetFileLocation(stepNode)
		steps = append(steps, step)
	}
	*s = steps
	return nil
}
