package models

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
	Types     []string `mapstructure:"types"`
	Workflows []string `mapstructure:"workflows"`
	Ref       `mapstructure:"ref,squash"`
}

type Events map[string]*struct {
	Types []string `mapstructure:"types"`
}

type On struct {
	Push              *Ref              `mapstructure:"push"`
	PullRequest       *Ref              `mapstructure:"pull_request"`
	PullRequestTarget *Ref              `mapstructure:"pull_request_target"`
	WorkflowCall      *WorkflowCall     `mapstructure:"workflow_call"`
	Schedule          *[]Cron           `mapstructure:"schedule"`
	WorkflowRun       *WorkflowRun      `mapstructure:"workflow_run"`
	WorkflowDispatch  *WorkflowDispatch `mapstructure:"workflow_dispatch"`
	Events
}

type Cron struct {
	Cron string `mapstructure:"cron" yarn:"cron"`
}
