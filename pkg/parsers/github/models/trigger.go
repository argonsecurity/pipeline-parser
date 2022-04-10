package models

type Gitevent struct {
	Paths          []string `mapstructure:"paths"`
	PathsIgnore    []string `mapstructure:"paths-ignore"`
	Branches       []string `mapstructure:"branches"`
	BranchesIgnore []string `mapstructure:"branches-ignore"`
}

type Inputs map[string]struct {
	Description string      `mapstructure:"description"`
	Default     interface{} `mapstructure:"default"`
	Required    bool        `mapstructure:"required"`
	Type        string      `mapstructure:"type"`
	Options     []string    `mapstructure:"options,omitempty"`
}

type Outputs map[string]*struct {
	Description string `mapstructure:"description"`
	Value       string `mapstructure:"value"`
}

type WorkflowDispatch struct {
	Inputs Inputs `mapstructure:"inputs"`
}

type WorkflowCall struct {
	Inputs  Inputs  `mapstructure:"inputs"`
	Outputs Outputs `mapstructure:"outputs"`
	Secrets map[string]*struct {
		Description string `mapstructure:"description"`
		Required    bool   `mapstructure:"required"`
	}
}

type WorkflowRun struct {
	Types          []string `mapstructure:"types"`
	Workflows      []string `mapstructure:"workflows"`
	Branches       []string `mapstructure:"branches"`
	BranchesIgnore []string `mapstructure:"branches-ignore"`
}

type Events map[string]*struct {
	Types []string `mapstructure:"types"`
}

type On struct {
	Push              *Gitevent         `mapstructure:"push"`
	PullRequest       *Gitevent         `mapstructure:"pull_request"`
	PullRequestTarget *Gitevent         `mapstructure:"pull_request_target"`
	WorkflowCall      *WorkflowCall     `mapstructure:"workflow_call"`
	Schedule          map[string]string `mapstructure:"schedule"`
	WorkflowRun       *WorkflowRun      `mapstructure:"workflow_run"`
	WorkflowDispatch  *WorkflowDispatch `mapstructure:"workflow_dispatch"`
	Events
}
