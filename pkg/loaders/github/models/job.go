package models

import (
	"errors"

	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
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
	CancelInProgress *bool   `yaml:"cancel-in-progress,omitempty"`
	Group            *string `yaml:"group"`
}

func (c *Concurrency) UnmarshalYAML(node *yaml.Node) error {
	(*c).Group = &node.Value
	return nil
}

type Job struct {
	ID              *string                  `yaml:"id"`
	Concurrency     *Concurrency             `yaml:"concurrency,omitempty"`
	Container       interface{}              `yaml:"container,omitempty"`
	ContinueOnError bool                     `yaml:"continue-on-error,omitempty"`
	Defaults        *Defaults                `yaml:"defaults,omitempty"`
	Env             *EnvironmentVariablesRef `yaml:"env,omitempty"`
	Environment     interface{}              `yaml:"environment,omitempty"`
	If              string                   `yaml:"if,omitempty"`
	Name            string                   `yaml:"name,omitempty"`
	Needs           *Needs                   `yaml:"needs,omitempty"`
	Outputs         map[string]string        `yaml:"outputs,omitempty"`
	Permissions     *PermissionsEvent        `yaml:"permissions,omitempty"`
	RunsOn          *RunsOn                  `yaml:"runs-on"`
	Services        map[string]*Container    `yaml:"services,omitempty"`
	Steps           *Steps                   `yaml:"steps,omitempty"`
	Strategy        *Strategy                `yaml:"strategy,omitempty"`
	TimeoutMinutes  *float64                 `yaml:"timeout-minutes,omitempty"`
	FileLocation    *models.FileLocation
}

type ReusableWorkflowCallJob struct {
	ID           *string           `yaml:"id"`
	If           string            `yaml:"if,omitempty"`
	Name         string            `yaml:"name,omitempty"`
	Needs        *Needs            `yaml:"needs,omitempty"`
	Permissions  *PermissionsEvent `yaml:"permissions,omitempty"`
	Secrets      interface{}       `yaml:"secrets,omitempty"`
	Uses         string            `yaml:"uses"`
	With         map[string]any
	FileLocation *models.FileLocation
}

func (j *Jobs) UnmarshalYAML(node *yaml.Node) error {
	var jobIdsToJobNodes map[string]yaml.Node
	if err := node.Decode(&jobIdsToJobNodes); err != nil {
		return err
	}

	normalJobs := make(map[string]*Job, 0)
	reusableWorkflowCallJobs := make(map[string]*ReusableWorkflowCallJob, 0)

	for i := 0; i < len(node.Content); i += 2 {
		jobIDNode := node.Content[i]
		jobNode := node.Content[i+1]

		jobID := jobIDNode.Value

		if isJobReusableWorkflowJob(jobNode) {
			reusableJob := &ReusableWorkflowCallJob{ID: utils.GetPtr(jobID)}
			if err := jobNode.Decode(reusableJob); err != nil {
				return err
			}
			reusableJob.FileLocation = loadersUtils.GetMapKeyFileLocation(jobIDNode, jobNode)
			reusableWorkflowCallJobs[jobID] = reusableJob
		} else {
			job := &Job{ID: utils.GetPtr(jobID)}
			if err := jobNode.Decode(job); err != nil {
				return err
			}
			job.FileLocation = loadersUtils.GetMapKeyFileLocation(jobIDNode, jobNode)
			normalJobs[jobID] = job
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
