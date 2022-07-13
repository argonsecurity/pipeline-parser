package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Steps []Step

type Target struct {
	Container         string   `yaml:"container,omitempty"`
	Commands          string   `yaml:"commands,omitempty"`
	SettableVariables []string `yaml:"settableVariables,omitempty"`
}

type Step struct {
	Name                    string                   `yaml:"name,omitempty"`
	Condition               string                   `yaml:"condition,omitempty"`
	ContinueOnError         bool                     `yaml:"continueOnError,omitempty"`
	DisplayName             string                   `yaml:"displayName,omitempty"`
	Target                  *Target                  `yaml:"target,omitempty"`
	Enabled                 bool                     `yaml:"enabled,omitempty"`
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
	SaveCache               string                   `yaml:"saveCache,omitempty"`
	Script                  string                   `yaml:"script,omitempty"`
	Task                    string                   `yaml:"task,omitempty"`
	Inputs                  map[string]any           `yaml:"inputs,omitempty"`
	Template                string                   `yaml:"template,omitempty"`
	Parameters              map[string]any           `yaml:"parameters,omitempty"`

	FileReference *models.FileReference
}

type Checkout struct {
	Checkout           string `yaml:"checkout,omitempty"`
	Clean              bool   `yaml:"clean,omitempty"`
	FetchDepth         int    `yaml:"fetchDepth,omitempty"`
	Lfs                bool   `yaml:"lfs,omitempty"`
	PersistCredentials bool   `yaml:"persistCredentials,omitempty"`
	Submodules         string `yaml:"submodules,omitempty"`
	Path               string `yaml:"path,omitempty"`
}

type Download struct {
	Download string `yaml:"download,omitempty"`
	Artifact string `yaml:"artifact,omitempty"`
	Patterns string `yaml:"patterns,omitempty"`
}

type DownloadBuild struct {
	Download string `yaml:"download,omitempty"`
	Artifact string `yaml:"artifact,omitempty"`
	Path     string `yaml:"path,omitempty"`
	Patterns string `yaml:"patterns,omitempty"`
}

type GetPackage struct {
	GetPackage string `yaml:"getPackage,omitempty"`
	Path       string `yaml:"path,omitempty"`
}

type Powershell struct {
	Powershell            string `yaml:"powershell,omitempty"`
	ErrorActionPreference string `yaml:"errorActionPreference,omitempty"`
	FailOnStderr          bool   `yaml:"failOnStderr,omitempty"`
	IgnoreLASTEXITCODE    bool   `yaml:"ignoreLASTEXITCODE,omitempty"`
	WorkingDirectory      string `yaml:"workingDirectory,omitempty"`
}

type Publish struct {
	Publish  string `yaml:"publish,omitempty"`
	Artifact string `yaml:"artifact,omitempty"`
}

type Pwsh struct {
	Pwsh                  string `yaml:"pwsh,omitempty"`
	ErrorActionPreference string `yaml:"errorActionPreference,omitempty"`
	FailOnStderr          bool   `yaml:"failOnStderr,omitempty"`
	IgnoreLASTEXITCODE    bool   `yaml:"ignoreLASTEXITCODE,omitempty"`
	WorkingDirectory      string `yaml:"workingDirectory,omitempty"`
}

type RestoreCache struct {
	RestoreCache string `yaml:"restoreCache,omitempty"`
	Path         string `yaml:"path,omitempty"`
}

type SaveCache struct {
	SaveCache string `yaml:"saveCache,omitempty"`
	Path      string `yaml:"path,omitempty"`
}
type Script struct {
	Script           string `yaml:"script,omitempty"`
	FailOnStderr     bool   `yaml:"failOnStderr,omitempty"`
	WorkingDirectory string `yaml:"workingDirectory,omitempty"`
}
type Task struct {
	Task   string         `yaml:"task,omitempty"`
	Inputs map[string]any `yaml:"inputs,omitempty"`
}

type Template struct {
	Template   string         `yaml:"template,omitempty"`
	Parameters map[string]any `yaml:"parameters,omitempty"`
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
