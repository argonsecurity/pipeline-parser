package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type Container struct {
	Credentials *Credentials  `mapstructure:"credentials,omitempty" yaml:"credentials,omitempty"`
	Env         interface{}   `mapstructure:"env,omitempty" yaml:"env,omitempty"`
	Image       string        `mapstructure:"image" yaml:"image"`
	Options     string        `mapstructure:"options,omitempty" yaml:"options,omitempty"`
	Ports       []interface{} `mapstructure:"ports,omitempty" yaml:"ports,omitempty"`
	Volumes     []string      `mapstructure:"volumes,omitempty" yaml:"volumes,omitempty"`
}

type Credentials struct {
	Password string `mapstructure:"password,omitempty" yaml:"password,omitempty"`
	Username string `mapstructure:"username,omitempty" yaml:"username,omitempty"`
}

type Defaults struct {
	Run *Run `mapstructure:"run,omitempty" yaml:"run,omitempty"`
}

type Environment struct {
	Name string `mapstructure:"name" yaml:"name"`
	Url  string `mapstructure:"url,omitempty" yaml:"url,omitempty"`
}

type PermissionsEvent struct {
	Actions            string `mapstructure:"actions,omitempty" yaml:"actions,omitempty"`
	Checks             string `mapstructure:"checks,omitempty" yaml:"checks,omitempty"`
	Contents           string `mapstructure:"contents,omitempty" yaml:"contents,omitempty"`
	Deployments        string `mapstructure:"deployments,omitempty" yaml:"deployments,omitempty"`
	Discussions        string `mapstructure:"discussions,omitempty" yaml:"discussions,omitempty"`
	IdToken            string `mapstructure:"id-token,omitempty" yaml:"id-token,omitempty"`
	Issues             string `mapstructure:"issues,omitempty" yaml:"issues,omitempty"`
	Packages           string `mapstructure:"packages,omitempty" yaml:"packages,omitempty"`
	Pages              string `mapstructure:"pages,omitempty" yaml:"pages,omitempty"`
	PullRequests       string `mapstructure:"pull-requests,omitempty" yaml:"pull-requests,omitempty"`
	RepositoryProjects string `mapstructure:"repository-projects,omitempty" yaml:"repository-projects,omitempty"`
	SecurityEvents     string `mapstructure:"security-events,omitempty" yaml:"security-events,omitempty"`
	Statuses           string `mapstructure:"statuses,omitempty" yaml:"statuses,omitempty"`
}

// Ref
type Ref struct {
	Branches       []string `mapstructure:"branches,omitempty" yaml:"branches,omitempty"`
	BranchesIgnore []string `mapstructure:"branches-ignore,omitempty" yaml:"branches-ignore,omitempty"`
	Paths          []string `mapstructure:"paths,omitempty" yaml:"paths,omitempty"`
	PathsIgnore    []string `mapstructure:"paths-ignore,omitempty" yaml:"paths-ignore,omitempty"`
	Tags           []string `mapstructure:"tags,omitempty" yaml:"tags,omitempty"`
	TagsIgnore     []string `mapstructure:"tags-ignore,omitempty" yaml:"tags-ignore,omitempty"`
}

// Workflow
type Workflow struct {
	Concurrency *Concurrency      `mapstructure:"concurrency,omitempty" yaml:"concurrency,omitempty"`
	Defaults    *Defaults         `mapstructure:"defaults,omitempty" yaml:"defaults,omitempty"`
	Env         interface{}       `mapstructure:"env,omitempty" yaml:"env,omitempty"`
	Jobs        *Jobs             `mapstructure:"jobs" yaml:"jobs"`
	Name        string            `mapstructure:"name,omitempty" yaml:"name,omitempty"`
	On          interface{}       `mapstructure:"on" yaml:"on"`
	Permissions *PermissionsEvent `mapstructure:"permissions,omitempty" yaml:"permissions,omitempty"`
}

type Run struct {
	Shell            interface{} `mapstructure:"shell,omitempty" yaml:"shell,omitempty"`
	WorkingDirectory string      `mapstructure:"working-directory,omitempty" yaml:"working-directory,omitempty"`
}

type Strategy struct {
	FailFast    bool        `mapstructure:"fail-fast,omitempty" yaml:"fail-fast,omitempty"`
	Matrix      interface{} `mapstructure:"matrix" yaml:"matrix"`
	MaxParallel float64     `mapstructure:"max-parallel,omitempty" yaml:"max-parallel,omitempty"`
}

func (j *Jobs) UnmarshalYAML(node *yaml.Node) error {
	var v map[string]any
	if err := node.Decode(&v); err != nil {
		return err
	}

	normalJobs := make(map[string]*NormalJob, 0)
	reusableWorkflowCallJobs := make(map[string]*ReusableWorkflowCallJob, 0)

	for k, v := range v {
		var normalJob *NormalJob
		var reusableWorkflowCallJob *ReusableWorkflowCallJob
		dc := &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.TextUnmarshallerHookFunc(),
			),
			Result: &normalJob,
		}
		decoder, err := mapstructure.NewDecoder(dc)
		if err != nil {
			return err
		}
		if err := decoder.Decode(v); err == nil {
			normalJobs[k] = normalJob
			continue
		} else if err := mapstructure.Decode(v, &reusableWorkflowCallJob); err == nil {
			reusableWorkflowCallJobs[k] = reusableWorkflowCallJob
			continue
		} else {
			return errors.New("unable to unmarshal jobs")
		}

	}
	*j = Jobs{
		NormalJobs:               normalJobs,
		ReusableWorkflowCallJobs: reusableWorkflowCallJobs,
	}
	return nil
}

func (strct *Container) UnmarshalJSON(b []byte) error {
	imageReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "credentials":
			if err := json.Unmarshal([]byte(v), &strct.Credentials); err != nil {
				return err
			}
		case "env":
			if err := json.Unmarshal([]byte(v), &strct.Env); err != nil {
				return err
			}
		case "image":
			if err := json.Unmarshal([]byte(v), &strct.Image); err != nil {
				return err
			}
			imageReceived = true
		case "options":
			if err := json.Unmarshal([]byte(v), &strct.Options); err != nil {
				return err
			}
		case "ports":
			if err := json.Unmarshal([]byte(v), &strct.Ports); err != nil {
				return err
			}
		case "volumes":
			if err := json.Unmarshal([]byte(v), &strct.Volumes); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if image (a required property) was received
	if !imageReceived {
		return errors.New("\"image\" is required but was not present")
	}
	return nil
}

func (strct *Defaults) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "run":
			if err := json.Unmarshal([]byte(v), &strct.Run); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

func (strct *Environment) UnmarshalJSON(b []byte) error {
	nameReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
			nameReceived = true
		case "url":
			if err := json.Unmarshal([]byte(v), &strct.Url); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if name (a required property) was received
	if !nameReceived {
		return errors.New("\"name\" is required but was not present")
	}
	return nil
}

func (strct *NormalJob) UnmarshalYAML(node *yaml.Node) error {
	var jsonMap map[string]json.RawMessage
	if err := node.Decode(&jsonMap); err != nil {
		return err
	}
	return strct.unmarshal(jsonMap)
}

func (strct *NormalJob) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}

	return strct.unmarshal(jsonMap)
}

func (strct *NormalJob) unmarshal(jsonMap map[string]json.RawMessage) error {
	var runsOnReceived bool
	for k, v := range jsonMap {
		switch k {
		case "concurrency":
			if err := json.Unmarshal([]byte(v), &strct.Concurrency); err != nil {
				return err
			}
		case "container":
			if err := json.Unmarshal([]byte(v), &strct.Container); err != nil {
				return err
			}
		case "continue-on-error":
			if err := json.Unmarshal([]byte(v), &strct.ContinueOnError); err != nil {
				return err
			}
		case "defaults":
			if err := json.Unmarshal([]byte(v), &strct.Defaults); err != nil {
				return err
			}
		case "env":
			if err := json.Unmarshal([]byte(v), &strct.Env); err != nil {
				return err
			}
		case "environment":
			if err := json.Unmarshal([]byte(v), &strct.Environment); err != nil {
				return err
			}
		case "if":
			if err := json.Unmarshal([]byte(v), &strct.If); err != nil {
				return err
			}
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
		case "needs":
			if err := json.Unmarshal([]byte(v), &strct.Needs); err != nil {
				return err
			}
		case "outputs":
			if err := json.Unmarshal([]byte(v), &strct.Outputs); err != nil {
				return err
			}
		case "permissions":
			if err := json.Unmarshal([]byte(v), &strct.Permissions); err != nil {
				return err
			}
		case "runs-on":
			if err := json.Unmarshal([]byte(v), &strct.RunsOn); err != nil {
				return err
			}
			runsOnReceived = true
		case "services":
			if err := json.Unmarshal([]byte(v), &strct.Services); err != nil {
				return err
			}
		case "steps":
			if err := json.Unmarshal([]byte(v), &strct.Steps); err != nil {
				return err
			}
		case "strategy":
			if err := json.Unmarshal([]byte(v), &strct.Strategy); err != nil {
				return err
			}
		case "timeout-minutes":
			if err := json.Unmarshal([]byte(v), &strct.TimeoutMinutes); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if runs-on (a required property) was received
	if !runsOnReceived {
		return errors.New("\"runs-on\" is required but was not present")
	}
	return nil
}

func (strct *PermissionsEvent) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "actions":
			if err := json.Unmarshal([]byte(v), &strct.Actions); err != nil {
				return err
			}
		case "checks":
			if err := json.Unmarshal([]byte(v), &strct.Checks); err != nil {
				return err
			}
		case "contents":
			if err := json.Unmarshal([]byte(v), &strct.Contents); err != nil {
				return err
			}
		case "deployments":
			if err := json.Unmarshal([]byte(v), &strct.Deployments); err != nil {
				return err
			}
		case "discussions":
			if err := json.Unmarshal([]byte(v), &strct.Discussions); err != nil {
				return err
			}
		case "id-token":
			if err := json.Unmarshal([]byte(v), &strct.IdToken); err != nil {
				return err
			}
		case "issues":
			if err := json.Unmarshal([]byte(v), &strct.Issues); err != nil {
				return err
			}
		case "packages":
			if err := json.Unmarshal([]byte(v), &strct.Packages); err != nil {
				return err
			}
		case "pages":
			if err := json.Unmarshal([]byte(v), &strct.Pages); err != nil {
				return err
			}
		case "pull-requests":
			if err := json.Unmarshal([]byte(v), &strct.PullRequests); err != nil {
				return err
			}
		case "repository-projects":
			if err := json.Unmarshal([]byte(v), &strct.RepositoryProjects); err != nil {
				return err
			}
		case "security-events":
			if err := json.Unmarshal([]byte(v), &strct.SecurityEvents); err != nil {
				return err
			}
		case "statuses":
			if err := json.Unmarshal([]byte(v), &strct.Statuses); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

func (strct *ReusableWorkflowCallJob) UnmarshalJSON(b []byte) error {
	usesReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "if":
			if err := json.Unmarshal([]byte(v), &strct.If); err != nil {
				return err
			}
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
		case "needs":
			if err := json.Unmarshal([]byte(v), &strct.Needs); err != nil {
				return err
			}
		case "permissions":
			if err := json.Unmarshal([]byte(v), &strct.Permissions); err != nil {
				return err
			}
		case "secrets":
			if err := json.Unmarshal([]byte(v), &strct.Secrets); err != nil {
				return err
			}
		case "uses":
			if err := json.Unmarshal([]byte(v), &strct.Uses); err != nil {
				return err
			}
			usesReceived = true
		case "with":
			if err := json.Unmarshal([]byte(v), &strct.With); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if uses (a required property) was received
	if !usesReceived {
		return errors.New("\"uses\" is required but was not present")
	}
	return nil
}

func (strct *Workflow) UnmarshalJSON(b []byte) error {
	jobsReceived := false
	onReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "concurrency":
			if err := json.Unmarshal([]byte(v), &strct.Concurrency); err != nil {
				return err
			}
		case "defaults":
			if err := json.Unmarshal([]byte(v), &strct.Defaults); err != nil {
				return err
			}
		case "env":
			if err := json.Unmarshal([]byte(v), &strct.Env); err != nil {
				return err
			}
		case "jobs":
			if err := json.Unmarshal([]byte(v), &strct.Jobs); err != nil {
				return err
			}
			jobsReceived = true
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
		case "on":
			if err := json.Unmarshal([]byte(v), &strct.On); err != nil {
				return err
			}
			onReceived = true
		case "permissions":
			if err := json.Unmarshal([]byte(v), &strct.Permissions); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if jobs (a required property) was received
	if !jobsReceived {
		return errors.New("\"jobs\" is required but was not present")
	}
	// check if on (a required property) was received
	if !onReceived {
		return errors.New("\"on\" is required but was not present")
	}
	return nil
}

func (strct *Run) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "shell":
			if err := json.Unmarshal([]byte(v), &strct.Shell); err != nil {
				return err
			}
		case "working-directory":
			if err := json.Unmarshal([]byte(v), &strct.WorkingDirectory); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

func (strct *Step) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "continue-on-error":
			if err := json.Unmarshal([]byte(v), &strct.ContinueOnError); err != nil {
				return err
			}
		case "env":
			if err := json.Unmarshal([]byte(v), &strct.Env); err != nil {
				return err
			}
		case "id":
			if err := json.Unmarshal([]byte(v), &strct.Id); err != nil {
				return err
			}
		case "if":
			if err := json.Unmarshal([]byte(v), &strct.If); err != nil {
				return err
			}
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
		case "run":
			if err := json.Unmarshal([]byte(v), &strct.Run); err != nil {
				return err
			}
		case "shell":
			if err := json.Unmarshal([]byte(v), &strct.Shell); err != nil {
				return err
			}
		case "timeout-minutes":
			if err := json.Unmarshal([]byte(v), &strct.TimeoutMinutes); err != nil {
				return err
			}
		case "uses":
			if err := json.Unmarshal([]byte(v), &strct.Uses); err != nil {
				return err
			}
		case "with":
			if err := json.Unmarshal([]byte(v), &strct.With); err != nil {
				return err
			}
		case "working-directory":
			if err := json.Unmarshal([]byte(v), &strct.WorkingDirectory); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

func (strct *Strategy) UnmarshalJSON(b []byte) error {
	matrixReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "fail-fast":
			if err := json.Unmarshal([]byte(v), &strct.FailFast); err != nil {
				return err
			}
		case "matrix":
			if err := json.Unmarshal([]byte(v), &strct.Matrix); err != nil {
				return err
			}
			matrixReceived = true
		case "max-parallel":
			if err := json.Unmarshal([]byte(v), &strct.MaxParallel); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if matrix (a required property) was received
	if !matrixReceived {
		return errors.New("\"matrix\" is required but was not present")
	}
	return nil
}
