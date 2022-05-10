package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	jobModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/job"
)

type Job struct {
	AfterScript  []*common.Script `yaml:"after_script"`
	BeforeScript []*common.Script `yaml:"before_script"`

	AllowFailure  *jobModels.AllowFailure `yaml:"allow_failure"`
	Artifacts     *Artifacts              `yaml:"artifacts"`
	Cache         *common.Cache           `yaml:"cache"`
	Coverage      string                  `yaml:"coverage"`
	Dependencies  []string                `yaml:"dependencies"`
	Environment   interface{}             `yaml:"environment"` // TODO: implement
	Extends       interface{}             `yaml:"extends"`
	Image         *common.Image           `yaml:"image"`
	Inherit       *jobModels.Inherit      `yaml:"inherit"`
	Interruptible bool                    `yaml:"interruptible"`
	Needs         []interface{}           `yaml:"needs"`
	Parallel      *jobModels.Parallel     `yaml:"parallel"`
	Release       *Release                `yaml:"release"`
	ResourceGroup string                  `yaml:"resource_group"`
	Retry         interface{}             `yaml:"retry"`
	Rules         []*common.RulesItems    `yaml:"rules"`
	Script        *common.Script          `yaml:"script"`
	Secrets       *Secrets                `yaml:"secrets"`
	Services      []interface{}           `yaml:"services"` // TODO: implement
	Stage         string                  `yaml:"stage"`
	StartIn       string                  `yaml:"start_in"`
	Tags          []string                `yaml:"tags"`
	Timeout       string                  `yaml:"timeout"`
	Trigger       *jobModels.Trigger      `yaml:"trigger"`
	Variables     *common.Variables       `yaml:"variables"`
	When          string                  `yaml:"when"`

	Except *jobModels.Conditions `yaml:"except"`
	Only   *jobModels.Conditions `yaml:"only"`
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
	Vault interface{} `yaml:"vault"`
}
