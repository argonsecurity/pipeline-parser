package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type ExecutionUnitRef struct {
	ExecutionUnit *ExecutionUnit
	FileReference *models.FileReference
}

type ExecutionUnit struct {
	AfterScript []Script         `yaml:"after-script"` // Commands inside an after-script section will run when the step succeeds or fails. This; could be useful for clean up commands, test coverage, notifications, or rollbacks you; might want to run, especially if your after-script uses the value of; BITBUCKET_EXIT_CODE.; ; Note: If any commands in the after-script section fail:; ; * we won't run any more commands in that section; ; * it will not effect the reported status of the step.
	Artifacts   *Artifacts       `yaml:"artifacts"`
	Caches      []string         `yaml:"caches"` // Caches enabled for the step
	Clone       *Clone           `yaml:"clone,omitempty"`
	Deployment  string           `yaml:"deployment,omitempty"` // Sets the type of environment for your deployment step, used in the Deployments dashboard.
	Image       *Image           `yaml:"image"`
	MaxTime     *int64           `yaml:"max-time,omitempty"`
	Name        string           `yaml:"name,omitempty"` // You can add a name to a step to make displays and reports easier to read and understand.
	RunsOn      []string         `yaml:"runs-on"`        // self-hosted runner labels
	Script      []Script         `yaml:"script"`         // Commands to execute in the step
	Services    []string         `yaml:"services"`       // Services enabled for the step
	Size        *Size            `yaml:"size,omitempty"`
	Trigger     *StepTriggerType `yaml:"trigger,omitempty"` // Specifies whether a step will run automatically or only after someone manually triggers; it. You can define the trigger type as manual or automatic. If the trigger type is not; defined, the step defaults to running automatically. The first step cannot be manual. If; you want to have a whole pipeline only run from a manual trigger then use a custom; pipeline.
}

type Script struct {
	PipeToExecute *PipeToExecute
	String        string
	FileReference *models.FileReference
}

type PipeToExecute struct {
	Pipe      string                  `yaml:"pipe"`                // Pipes make complex tasks easier, by doing a lot of the work behind the scenes.; This means you can just select which pipe you want to use, and supply the necessary; variables.; You can look at the repository for the pipe to see what commands it is running.; ; Learn more about pipes: https://confluence.atlassian.com/bitbucket/pipes-958765631.html
	Variables EnvironmentVariablesRef `yaml:"variables,omitempty"` // Environment variables passed to the pipe
}

func (s *ExecutionUnitRef) UnmarshalYAML(node *yaml.Node) error {
	s.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&s.ExecutionUnit)
}

func (s *Script) UnmarshalYAML(value *yaml.Node) error {
	var pipeToExecute PipeToExecute
	if err := value.Decode(&pipeToExecute); err == nil {
		s.PipeToExecute = &pipeToExecute
		return nil
	}

	var stringToExecute string
	if err := value.Decode(&stringToExecute); err == nil {
		s.String = stringToExecute
		s.FileReference = loadersUtils.GetFileReference(value)
		return nil
	}

	return nil
}
