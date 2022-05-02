package models

import (
	"errors"

	loaderUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
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

func (n *Needs) UnmarshalYAML(node *yaml.Node) error {
	var needs []string
	if err := node.Decode(&needs); err == nil {
		*n = needs
		return nil
	}

	var needsString string
	if err := node.Decode(&needsString); err == nil {
		*n = []string{needsString}
		return nil
	}
	return errors.New("unable to decode needs")
}

type Concurrency struct {
	CancelInProgress *bool   `mapstructure:"cancel-in-progress,omitempty" yaml:"cancel-in-progress,omitempty"`
	Group            *string `mapstructure:"group" yaml:"group"`
}

func (c *Concurrency) UnmarshalYAML(node *yaml.Node) error {
	(*c).Group = &node.Value
	return nil
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
	FileLocation    *models.FileLocation
}

type ReusableWorkflowCallJob struct {
	ID           *string              `mapstructure:"id" yaml:"id"`
	If           string               `mapstructure:"if,omitempty" yaml:"if,omitempty"`
	Name         string               `mapstructure:"name,omitempty" yaml:"name,omitempty"`
	Needs        *Needs               `mapstructure:"needs,omitempty" yaml:"needs,omitempty"`
	Permissions  *PermissionsEvent    `mapstructure:"permissions,omitempty" yaml:"permissions,omitempty"`
	Secrets      interface{}          `mapstructure:"secrets,omitempty" yaml:"secrets,omitempty"`
	Uses         string               `mapstructure:"uses" yaml:"uses"`
	With         map[string]any       `mapstructure:"with,omitempty"`
	FileLocation *models.FileLocation `mapstructure:""`
}

func (j *Jobs) UnmarshalYAML(node *yaml.Node) error {
	var nodeAsMap map[string]yaml.Node
	if err := node.Decode(&nodeAsMap); err != nil {
		return err
	}

	normalJobs := make(map[string]*Job, 0)
	reusableWorkflowCallJobs := make(map[string]*ReusableWorkflowCallJob, 0)

	for jobId, jobObject := range nodeAsMap {
		if isJobReusableWorkflowJob(jobObject) {
			reusableJob := &ReusableWorkflowCallJob{ID: utils.GetPtr(jobId)}
			if err := jobObject.Decode(reusableJob); err != nil {
				return err
			}
			reusableJob.FileLocation = loaderUtils.GetFileLocation(node)
			reusableWorkflowCallJobs[jobId] = reusableJob
		} else {
			job := &Job{ID: utils.GetPtr(jobId)}
			if err := jobObject.Decode(job); err != nil {
				return err
			}
			job.FileLocation = loaderUtils.GetFileLocation(node)
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
