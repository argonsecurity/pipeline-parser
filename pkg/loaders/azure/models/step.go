package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Steps []Step

type StepTarget struct {
	Container         string   `yaml:"container,omitempty"`
	Commands          string   `yaml:"commands,omitempty"`
	SettableVariables []string `yaml:"settableVariables,omitempty"`
	FileReference     *models.FileReference
}

type TaskInputs struct {
	Inputs        map[string]any
	FileReference *models.FileReference
}

type Step struct {
	Name                    string                   `yaml:"name,omitempty"`
	Condition               string                   `yaml:"condition,omitempty"`
	ContinueOnError         *bool                    `yaml:"continueOnError,omitempty"`
	DisplayName             string                   `yaml:"displayName,omitempty"`
	Target                  *StepTarget              `yaml:"target,omitempty"`
	Enabled                 *bool                    `yaml:"enabled,omitempty"`
	Env                     *EnvironmentVariablesRef `yaml:"env,omitempty"`
	TimeoutInMinutes        int                      `yaml:"timeoutInMinutes,omitempty"`
	RetryCountOnTaskFailure int                      `yaml:"retryCountOnTaskFailure,omitempty"`
	Bash                    string                   `yaml:"bash,omitempty"`
	Checkout                string                   `yaml:"checkout,omitempty"`
	Clean                   bool                     `yaml:"clean,omitempty"`
	FetchDepth              int                      `yaml:"fetchDepth,omitempty"`
	Lfs                     bool                     `yaml:"lfs,omitempty"`
	PersistCredentials      bool                     `yaml:"persistCredentials,omitempty"`
	Submodules              string                   `yaml:"submodules,omitempty"`
	Path                    string                   `yaml:"path,omitempty"`
	Download                string                   `yaml:"download,omitempty"`
	DownloadBuild           string                   `yaml:"downloadBuild,omitempty"`
	Artifact                string                   `yaml:"artifact,omitempty"`
	Patterns                string                   `yaml:"patterns,omitempty"`
	GetPackage              string                   `yaml:"getPackage,omitempty"`
	Powershell              string                   `yaml:"powershell,omitempty"`
	ErrorActionPreference   string                   `yaml:"errorActionPreference,omitempty"`
	FailOnStderr            bool                     `yaml:"failOnStderr,omitempty"`
	IgnoreLASTEXITCODE      bool                     `yaml:"ignoreLASTEXITCODE,omitempty"`
	WorkingDirectory        string                   `yaml:"workingDirectory,omitempty"`
	Publish                 string                   `yaml:"publish,omitempty"`
	Pwsh                    string                   `yaml:"pwsh,omitempty"`
	RestoreCache            string                   `yaml:"restoreCache,omitempty"`
	ReviewApp               string                   `yaml:"reviewApp,omitempty"`
	SaveCache               string                   `yaml:"saveCache,omitempty"`
	Script                  string                   `yaml:"script,omitempty"`
	Task                    string                   `yaml:"task,omitempty"`
	Inputs                  *TaskInputs              `yaml:"inputs,omitempty"`
	Template                `yaml:",inline"`

	FileReference *models.FileReference
}

func (t *StepTarget) UnmarshalYAML(node *yaml.Node) error {
	t.FileReference = loadersUtils.GetFileReference(node)
	if node.Tag == consts.StringTag {
		t.Container = node.Value
		return nil
	}

	return loadersUtils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		switch key {
		case "container":
			t.Container = value.Value
		case "commands":
			t.Commands = value.Value
		case "settableVariables":
			var settableVariables []string
			if err := loadersUtils.ParseSequenceOrOne(node, &settableVariables); err != nil {
				return err
			}
			t.SettableVariables = settableVariables
		}

		return nil
	})
}

func (ti *TaskInputs) UnmarshalYAML(node *yaml.Node) error {
	ti.FileReference = loadersUtils.GetFileReference(node)
	// The with block looks like this
	// with:
	//   key1: value1
	//   key2: value2
	// The node refers to the first input
	// We want to include the "with" in the File Reference so we have to subtract 1 from the line number and 2 from the column number
	ti.FileReference.StartRef.Line--
	ti.FileReference.StartRef.Column -= 2
	return node.Decode(&ti.Inputs)
}

func (s *Steps) UnmarshalYAML(node *yaml.Node) error {
	var steps []Step

	for _, stepNode := range node.Content {
		step, err := parseStep(stepNode)
		if err != nil {
			return err
		}
		steps = append(steps, step)
	}

	*s = steps
	return nil
}

func parseStep(node *yaml.Node) (Step, error) {
	var step Step
	if err := node.Decode(&step); err != nil {
		return step, err
	}
	step.FileReference = loadersUtils.GetFileReference(node)
	return step, nil
}
