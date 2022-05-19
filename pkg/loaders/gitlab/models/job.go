package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/job"
	jobModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/job"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Job struct {
	AfterScript  *common.Script `yaml:"after_script"`
	BeforeScript *common.Script `yaml:"before_script"`

	AllowFailure  *jobModels.AllowFailure         `yaml:"allow_failure"`
	Artifacts     *Artifacts                      `yaml:"artifacts"`
	Cache         *common.Cache                   `yaml:"cache"`
	Coverage      string                          `yaml:"coverage"`
	Dependencies  []string                        `yaml:"dependencies"`
	Environment   any                             `yaml:"environment"` // TODO: implement
	Extends       any                             `yaml:"extends"`
	Image         *common.Image                   `yaml:"image"`
	Inherit       *jobModels.Inherit              `yaml:"inherit"`
	Interruptible bool                            `yaml:"interruptible"`
	Needs         *job.Needs                      `yaml:"needs"`
	Parallel      *jobModels.Parallel             `yaml:"parallel"`
	Release       *Release                        `yaml:"release"`
	ResourceGroup string                          `yaml:"resource_group"`
	Retry         *common.Retry                   `yaml:"retry"`
	Rules         *common.Rules                   `yaml:"rules"`
	Script        *common.Script                  `yaml:"script"`
	Secrets       *Secrets                        `yaml:"secrets"`
	Services      []any                           `yaml:"services"` // TODO: implement
	Stage         string                          `yaml:"stage"`
	StartIn       string                          `yaml:"start_in"`
	Tags          []string                        `yaml:"tags"`
	Timeout       string                          `yaml:"timeout"`
	Trigger       *jobModels.Trigger              `yaml:"trigger"`
	Variables     *common.EnvironmentVariablesRef `yaml:"variables"`
	When          string                          `yaml:"when"`

	Except *jobModels.Controls `yaml:"except"`
	Only   *jobModels.Controls `yaml:"only"`

	FileReference *models.FileReference
}

type Release struct {
	Assets      *Assets  `yaml:"assets"`
	Description string   `yaml:"description"`
	Milestones  []string `yaml:"milestones"`
	Name        string   `yaml:"name"`
	Ref         string   `yaml:"ref"`
	ReleasedAt  string   `yaml:"released_at"`
	TagName     string   `yaml:"tag_name"`
}

type Assets struct {
	Links []*LinksItems `yaml:"links"`
}

type LinksItems struct {
	Filepath string `yaml:"filepath"`
	LinkType string `yaml:"link_type"`
	Name     string `yaml:"name"`
	Url      string `yaml:"url"`
}

type Secrets struct {
	AdditionalProperties map[string]*SecretsItem `yaml:"-"`
}

type SecretsItem struct {
	Vault any `yaml:"vault"`
}

// There's a bug in go-yaml and this is the only way we can currently have the jobs inside the ci configuration
// while keeping the job's file reference.
// Without overcomplicating, the bug won't allow us to both implement UnmarshalYAML and parse Job inline (as in, without a separate internal field)
func (j *Job) UnmarshalYAML(node *yaml.Node) error {
	j.FileReference = &models.FileReference{
		StartRef: &models.FileLocation{
			Line:   node.Line - 1, // We don't have access to the key node, and therefor we assume it's one line higher
			Column: 1,             // All jobs start from column 1
		},
		EndRef: utils.GetEndFileLocation(node),
	}
	return utils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		switch key {
		case "after_script":
			return value.Decode(&j.AfterScript)
		case "before_script":
			return value.Decode(&j.BeforeScript)
		case "allow_failure":
			return value.Decode(&j.AllowFailure)
		case "artifacts":
			return value.Decode(&j.Artifacts)
		case "cache":
			return value.Decode(&j.Cache)
		case "coverage":
			return value.Decode(&j.Coverage)
		case "dependencies":
			return value.Decode(&j.Dependencies)
		case "environment":
			return value.Decode(&j.Environment)
		case "extends":
			return value.Decode(&j.Extends)
		case "image":
			return value.Decode(&j.Image)
		case "inherit":
			return value.Decode(&j.Inherit)
		case "interruptible":
			return value.Decode(&j.Interruptible)
		case "needs":
			return value.Decode(&j.Needs)
		case "parallel":
			return value.Decode(&j.Parallel)
		case "release":
			return value.Decode(&j.Release)
		case "resource_group":
			return value.Decode(&j.ResourceGroup)
		case "retry":
			return value.Decode(&j.Retry)
		case "rules":
			return value.Decode(&j.Rules)
		case "script":
			return value.Decode(&j.Script)
		case "secrets":
			return value.Decode(&j.Secrets)
		case "services":
			return value.Decode(&j.Services)
		case "stage":
			return value.Decode(&j.Stage)
		case "start_in":
			return value.Decode(&j.StartIn)
		case "tags":
			return value.Decode(&j.Tags)
		case "timeout":
			return value.Decode(&j.Timeout)
		case "trigger":
			return value.Decode(&j.Trigger)
		case "variables":
			return value.Decode(&j.Variables)
		case "when":
			return value.Decode(&j.When)
		case "except":
			return value.Decode(&j.Except)
		case "only":
			return value.Decode(&j.Only)
		}
		return nil
	})
}
