package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Pipeline struct {
	Extends         *Extends      `yaml:"extends,omitempty"`
	Jobs            *Jobs         `yaml:"jobs,omitempty"`
	Stages          *Stages       `yaml:"stages,omitempty"`
	ContinueOnError *bool         `yaml:"continueOnError,omitempty"`
	Pool            *Pool         `yaml:"pool,omitempty"`
	Container       *JobContainer `yaml:"container,omitempty"`
	Name            string        `yaml:"name,omitempty"`
	Trigger         *TriggerRef   `yaml:"trigger,omitempty"`
	Parameters      *Parameters   `yaml:"parameters,omitempty"`
	PR              *PRRef        `yaml:"pr,omitempty"`
	Schedules       *Schedules    `yaml:"schedules,omitempty"`
	Resources       *Resources    `yaml:"resources,omitempty"`
	Steps           *Steps        `yaml:"steps,omitempty"`
	Variables       *Variables    `yaml:"variables,omitempty"`
	LockBehavior    string        `yaml:"lockBehavior,omitempty"`
}

type Parameter struct {
	Name          string   `yaml:"name,omitempty"`
	DisplayName   string   `yaml:"displayName,omitempty"`
	Type          string   `yaml:"type,omitempty"`
	Default       any      `yaml:"default,omitempty"`
	Values        []string `yaml:"values,omitempty"`
	FileReference *models.FileReference
}

type Parameters []Parameter

func (p *Parameters) UnmarshalYAML(node *yaml.Node) error {
	var parameters Parameters
	for _, parameterNode := range node.Content {
		parameter, err := parseParameter(parameterNode)
		if err != nil {
			return err
		}
		parameters = append(parameters, parameter)
	}
	*p = parameters
	return nil
}

func parseParameter(node *yaml.Node) (Parameter, error) {
	var parameter Parameter
	if err := node.Decode(&parameter); err != nil {
		return parameter, err
	}
	parameter.FileReference = loadersUtils.GetFileReference(node)
	return parameter, nil
}
