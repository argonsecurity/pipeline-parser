package models

import (
	"fmt"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"gopkg.in/yaml.v3"
)

type Jobs struct {
	NormalJobs               map[string]*Job
	ReusableWorkflowCallJobs map[string]*ReusableWorkflowCallJob
}

type Needs []string

func (n *Needs) UnmarshalYAML(node *yaml.Node) error {
	var tags []string
	var err error

	if node.Tag == consts.SequenceTag {
		if tags, err = loadersUtils.ParseYamlStringSequenceToSlice(node); err != nil {
			return err
		}
	} else if node.Tag == consts.StringTag {
		tags = []string{node.Value}
	} else {
		return fmt.Errorf("unexpected tag %s", node.Tag)
	}

	*n = tags
	return nil
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
	normalJobs := map[string]*Job{}
	reusableWorkflowCallJobs := map[string]*ReusableWorkflowCallJob{}

	for i := 0; i < len(node.Content); i += 2 {
		jobIDNode := node.Content[i]
		jobNode := node.Content[i+1]

		jobID := jobIDNode.Value

		if isReusableWorkflowJob(jobNode) {
			job, err := parseReusableWorkflowNode(jobIDNode, jobNode)
			if err != nil {
				return err
			}
			reusableWorkflowCallJobs[jobID] = job
		} else {
			job, err := parseJobNode(jobIDNode, jobNode)
			if err != nil {
				return err
			}
			normalJobs[jobID] = job
		}
	}

	*j = Jobs{
		NormalJobs:               normalJobs,
		ReusableWorkflowCallJobs: reusableWorkflowCallJobs,
	}
	return nil
}

func parseJobNode(jobID, job *yaml.Node) (*Job, error) {
	parsedJob := &Job{ID: utils.GetPtr(jobID.Value)}
	if err := job.Decode(parsedJob); err != nil {
		return nil, err
	}
	parsedJob.FileLocation = loadersUtils.GetMapKeyFileLocation(jobID, job)
	return parsedJob, nil
}

func parseReusableWorkflowNode(jobID, job *yaml.Node) (*ReusableWorkflowCallJob, error) {
	reusableJob := &ReusableWorkflowCallJob{ID: utils.GetPtr(jobID.Value)}
	if err := job.Decode(reusableJob); err != nil {
		return nil, err
	}
	reusableJob.FileLocation = loadersUtils.GetMapKeyFileLocation(jobID, job)
	return reusableJob, nil
}

func isReusableWorkflowJob(job *yaml.Node) bool {
	for _, node := range job.Content {
		if node.Tag == consts.StringTag && node.Value == "uses" {
			return true
		}
	}
	return false
}
