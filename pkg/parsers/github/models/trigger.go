package models

import (
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

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

type Events map[string]Event

type Event struct {
	Types []string `mapstructure:"types"`
}

type Cron struct {
	Cron string `mapstructure:"cron" yarn:"cron"`
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

func (on *On) UnmarshalYAML(node *yaml.Node) error {
	events := make([]string, 0)
	if err := node.Decode(&events); err == nil {
		for _, event := range events {
			on.unmarshalKey(event, nil)
		}
		return nil
	}

	var triggersMap map[string]any
	if err := node.Decode(&triggersMap); err != nil {
		return err
	}

	for k, v := range triggersMap {
		on.unmarshalKey(k, v)
	}
	return nil
}

func (on *On) unmarshalKey(key string, value any) {
	switch key {
	case "schedule":
		mapstructure.Decode(value, &on.Schedule)
	case "push":
		if value == nil {
			on.Push = &Ref{}
		} else {
			mapstructure.Decode(value, &on.Push)
		}
	case "pull_request":
		if value == nil {
			on.PullRequest = &Ref{}
		} else {
			mapstructure.Decode(value, &on.PullRequest)
		}
	case "pull_request_target":
		if value == nil {
			on.PullRequestTarget = &Ref{}
		} else {
			mapstructure.Decode(value, &on.PullRequestTarget)
		}
	case "workflow_call":
		if value == nil {
			on.WorkflowCall = &WorkflowCall{}
		} else {
			mapstructure.Decode(value, &on.WorkflowCall)
		}
	case "workflow_run":
		if value == nil {
			on.WorkflowRun = &WorkflowRun{}
		} else {
			mapstructure.Decode(value, &on.WorkflowRun)
		}
	case "workflow_dispatch":
		if value == nil {
			on.WorkflowDispatch = &WorkflowDispatch{}
		} else {
			mapstructure.Decode(value, &on.WorkflowDispatch)
		}
	default:
		on.Events[key] = value.(Event)
	}
}
