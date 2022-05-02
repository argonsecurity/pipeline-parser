package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
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
	Inputs       Inputs `mapstructure:"inputs"`
	FileLocation *models.FileLocation
}

type WorkflowCall struct {
	Inputs  Inputs  `mapstructure:"inputs"`
	Outputs Outputs `mapstructure:"outputs"`
	Secrets map[string]*struct {
		Description string `mapstructure:"description"`
		Required    bool   `mapstructure:"required"`
	}
	FileLocation *models.FileLocation
}

type WorkflowRun struct {
	Types        []string `mapstructure:"types"`
	Workflows    []string `mapstructure:"workflows"`
	Ref          `mapstructure:"ref,squash"`
	FileLocation *models.FileLocation
}

type Events map[string]Event

type Event struct {
	Types []string `mapstructure:"types"`
}

type Cron struct {
	Cron         string `mapstructure:"cron" yarn:"cron"`
	FileLocation *models.FileLocation
}

type On struct {
	Push              *Ref
	PullRequest       *Ref
	PullRequestTarget *Ref
	WorkflowCall      *WorkflowCall
	Schedule          *[]Cron
	WorkflowRun       *WorkflowRun
	WorkflowDispatch  *WorkflowDispatch `mapstructure:"workflow_dispatch"`
	FileLocation      *models.FileLocation
	Events
}

func (on *On) UnmarshalYAML(node *yaml.Node) error {
	events := make([]string, 0)
	if err := node.Decode(&events); err == nil {
		for _, event := range events {
			on.unmarshalKey(event, yaml.Node{})
		}
		return nil
	}

	var triggersMap map[string]yaml.Node
	if err := node.Decode(&triggersMap); err != nil {
		return err
	}

	for k, v := range triggersMap {
		if err := on.unmarshalKey(k, v); err != nil {
			return err
		}
	}
	on.FileLocation = loadersUtils.GetFileLocation(node)
	return nil
}

func (on *On) unmarshalKey(key string, node yaml.Node) error {
	var err error
	fileLocation := loadersUtils.GetFileLocation(&node)
	switch key {
	case "schedule":
		mapstructure.Decode(node, &on.Schedule)
	case "push":
		on.Push = &Ref{FileLocation: fileLocation}
		if !node.IsZero() {
			err = node.Decode(&on.Push)
		}
	case "pull_request":
		on.PullRequest = &Ref{FileLocation: fileLocation}
		if !node.IsZero() {
			err = node.Decode(&on.PullRequest)
		}
	case "pull_request_target":
		on.PullRequestTarget = &Ref{FileLocation: fileLocation}
		if !node.IsZero() {
			err = node.Decode(&on.PullRequestTarget)
		}
	case "workflow_call":
		on.WorkflowCall = &WorkflowCall{FileLocation: fileLocation}
		if !node.IsZero() {
			err = node.Decode(&on.WorkflowCall)
		}
	case "workflow_run":
		on.WorkflowRun = &WorkflowRun{FileLocation: fileLocation}
		if !node.IsZero() {
			err = node.Decode(&on.WorkflowRun)
		}
	case "workflow_dispatch":
		on.WorkflowDispatch = &WorkflowDispatch{FileLocation: fileLocation}
		if !node.IsZero() {
			err = node.Decode(&on.WorkflowDispatch)
		}
	default:
		if on.Events == nil {
			on.Events = make(Events)
		}
		var event Event
		if err := node.Decode(&event); err != nil {
			return err
		}
		on.Events[key] = event
	}
	return err
}
