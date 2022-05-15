package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
)

type GitlabCIConfiguration struct {
	AfterScript  *common.Script `yaml:"after_script"`
	BeforeScript *common.Script `yaml:"before_script"`
	Cache        *common.Cache  `yaml:"cache"`
	Default      *Default       `yaml:"default"`
	Image        *common.Image  `yaml:"image"`
	Include      *Include       `yaml:"include"`

	Pages    any   `yaml:"pages"`
	Services []any `yaml:"services"`

	// Groups jobs into stages. All jobs in one stage must complete before next stage is executed. Defaults to ['build', 'test', 'deploy'].
	Stages    []string                        `yaml:"stages"`
	Variables *common.EnvironmentVariablesRef `yaml:"variables"`
	Workflow  *Workflow                       `yaml:"workflow"`
	Jobs      map[string]*Job                 `yaml:",inline"`
}

type Default struct {
	AfterScript   []*common.Script `yaml:"after_script"`
	Artifacts     *Artifacts       `yaml:"artifacts"`
	BeforeScript  []*common.Script `yaml:"before_script"`
	Cache         *common.Cache    `yaml:"cache"`
	Image         *common.Image    `yaml:"image"`
	Interruptible bool             `yaml:"interruptible"`
	Retry         *common.Retry    `yaml:"retry"`
	Services      []any            `yaml:"services"`
	Tags          []string         `yaml:"tags"`
	Timeout       string           `yaml:"timeout"`
}

type Workflow struct {
	Rules *common.Rules `yaml:"rules"`
}
