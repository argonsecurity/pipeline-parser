package models

import (
	"errors"
	"reflect"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type Jobs struct {
	NormalJobs               map[string]*Job
	ReusableWorkflowCallJobs map[string]*ReusableWorkflowCallJob
}

type Needs []string

type Concurrency struct {
	CancelInProgress *bool   `mapstructure:"cancel-in-progress,omitempty" yaml:"cancel-in-progress,omitempty"`
	Group            *string `mapstructure:"group" yaml:"group"`
}

type Job struct {
	ID              *string                      `mapstructure:"id" yaml:"id"`
	Concurrency     *Concurrency                 `mapstructure:"concurrency,omitempty" yaml:"concurrency,omitempty"`
	Container       interface{}                  `mapstructure:"container,omitempty" yaml:"container,omitempty"`
	ContinueOnError bool                         `mapstructure:"continue-on-error,omitempty" yaml:"continue-on-error,omitempty"`
	Defaults        *Defaults                    `mapstructure:"defaults,omitempty" yaml:"defaults,omitempty"`
	Env             *models.EnvironmentVariables `mapstructure:"env,omitempty" yaml:"env,omitempty"`
	Environment     interface{}                  `mapstructure:"environment,omitempty" yaml:"environment,omitempty"`
	If              string                       `mapstructure:"if,omitempty" yaml:"if,omitempty"`
	Name            string                       `mapstructure:"name,omitempty" yaml:"name,omitempty"`
	Needs           *Needs                       `mapstructure:"needs,omitempty" yaml:"needs,omitempty"`
	Outputs         map[string]string            `mapstructure:"outputs,omitempty" yaml:"outputs,omitempty"`
	Permissions     *PermissionsEvent            `mapstructure:"permissions,omitempty" yaml:"permissions,omitempty"`
	RunsOn          *RunsOn                      `mapstructure:"runs-on" yaml:"runs-on"`
	Services        map[string]*Container        `mapstructure:"services,omitempty" yaml:"services,omitempty"`
	Steps           *[]Step                      `mapstructure:"steps,omitempty" yaml:"steps,omitempty"`
	Strategy        *Strategy                    `mapstructure:"strategy,omitempty" yaml:"strategy,omitempty"`
	TimeoutMinutes  *float64                     `mapstructure:"timeout-minutes,omitempty" yaml:"timeout-minutes,omitempty"`
}

type ReusableWorkflowCallJob struct {
	ID          *string           `mapstructure:"id" yaml:"id"`
	If          string            `mapstructure:"if,omitempty" yaml:"if,omitempty"`
	Name        string            `mapstructure:"name,omitempty" yaml:"name,omitempty"`
	Needs       *Needs            `mapstructure:"needs,omitempty" yaml:"needs,omitempty"`
	Permissions *PermissionsEvent `mapstructure:"permissions,omitempty" yaml:"permissions,omitempty"`
	Secrets     interface{}       `mapstructure:"secrets,omitempty" yaml:"secrets,omitempty"`
	Uses        string            `mapstructure:"uses" yaml:"uses"`
	With        map[string]any    `mapstructure:"with,omitempty"`
}

func (j *Jobs) UnmarshalYAML(node *yaml.Node) error {
	var nodeAsMap map[string]any
	if err := node.Decode(&nodeAsMap); err != nil {
		return err
	}

	normalJobs := make(map[string]*Job, 0)
	reusableWorkflowCallJobs := make(map[string]*ReusableWorkflowCallJob, 0)

	for jobId, jobObject := range nodeAsMap {
		if isJobReusableWorkflowJob(jobObject) {
			reusableJob := &ReusableWorkflowCallJob{ID: utils.GetPtr(jobId)}
			if err := decodeWithHooks(jobObject, reusableJob); err != nil {
				return err
			}
			reusableWorkflowCallJobs[jobId] = reusableJob
		} else {
			job := &Job{ID: utils.GetPtr(jobId)}
			if err := decodeWithHooks(jobObject, job); err != nil {
				return err
			}
			normalJobs[jobId] = job
		}
	}
	*j = Jobs{
		NormalJobs:               normalJobs,
		ReusableWorkflowCallJobs: reusableWorkflowCallJobs,
	}
	return nil
}

func isJobReusableWorkflowJob(job any) bool {
	var jobAsMap map[string]any
	if err := mapstructure.Decode(job, &jobAsMap); err != nil {
		return false
	}
	_, ok := jobAsMap["uses"]
	return ok
}

func decodeWithHooks[T any](data any, target T) error {
	dc := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.TextUnmarshallerHookFunc(),
			DecodeRunsOnHookFunc(),
			DecodeNeedsHookFunc(),
			DecodeTokenPermissionsHookFunc(),
		),
		Result: &target,
	}
	decoder, err := mapstructure.NewDecoder(dc)
	if err != nil {
		return err
	}
	return decoder.Decode(data)
}

func DecodeNeedsHookFunc() mapstructure.DecodeHookFunc {
	return func(f, t reflect.Type, data any) (any, error) {
		if t != reflect.TypeOf(Needs{}) {
			return data, nil
		}

		var needs []string
		if err := mapstructure.Decode(data, &needs); err == nil {
			return needs, nil
		}

		var needsString string
		if err := mapstructure.Decode(data, &needsString); err == nil {
			return []string{needsString}, nil
		}
		return nil, errors.New("unable to decode needs")
	}
}

func (c *Concurrency) UnmarshalText(text []byte) error {
	c.Group = utils.GetPtr(string(text))
	return nil
}
