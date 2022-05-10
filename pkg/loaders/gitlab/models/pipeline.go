package models

import "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"

type GitlabCIConfiguration struct {
	AdditionalProperties map[string]interface{} `yaml:"-"`
	AfterScript          []interface{}          `yaml:"after_script"`
	BeforeScript         []interface{}          `yaml:"before_script"`
	Cache                *common.Cache          `yaml:"cache"`
	Default              *Default               `yaml:"default"`
	Image                *common.Image          `yaml:"image"`

	// Can be `IncludeItem` or `IncludeItem[]`. Each `IncludeItem` will be a string, or an object with properties for the method if including external YAML file. The external content will be fetched, included and evaluated along the `.gitlab-ci.yml`.
	Include interface{} `yaml:"include"`

	// A special job used to upload static sites to Gitlab pages. Requires a `public/` directory with `artifacts.path` pointing to it.
	Pages    interface{}   `yaml:"pages"`
	Schema   string        `yaml:"$schema"`
	Services []interface{} `yaml:"services"`

	// Groups jobs into stages. All jobs in one stage must complete before next stage is executed. Defaults to ['build', 'test', 'deploy'].
	Stages    []string         `yaml:"stages"`
	Variables *GlobalVariables `yaml:"variables"`
	Workflow  *Workflow        `yaml:"workflow"`
}

type Default struct {
	AfterScript   []*common.Script `yaml:"after_script"`
	Artifacts     *Artifacts       `yaml:"artifacts"`
	BeforeScript  []*common.Script `yaml:"before_script"`
	Cache         *common.Cache    `yaml:"cache"`
	Image         *common.Image    `yaml:"image"`
	Interruptible bool             `yaml:"interruptible"`
	Retry         interface{}      `yaml:"retry"`
	Services      []interface{}    `yaml:"services"`
	Tags          []string         `yaml:"tags"`
	Timeout       string           `yaml:"timeout"`
}

type Workflow struct {
	Rules []*common.RulesItems `yaml:"rules"`
}

type GlobalVariables struct {
	AdditionalProperties map[string]interface{} `yaml:"-"`
}
