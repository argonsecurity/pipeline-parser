package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type JobType string

const (
	CIJobType         JobType = "job"
	DeploymentJobType JobType = "deployment"
	TemplateJobType   JobType = "template"
)

type Workspace struct {
	Clean string `yaml:"clean,omitempty"`
}

type Uses struct {
	Repositories []string `yaml:"repositories,omitempty"`
	Pools        []string `yaml:"pools,omitempty"`
}

type BaseJob struct {
	DisplayName            string            `yaml:"displayName,omitempty"`
	DependsOn              []string          `yaml:"dependsOn,omitempty"`
	Condition              string            `yaml:"condition,omitempty"`
	ContinueOnError        bool              `yaml:"continueOnError,omitempty"`
	TimeoutInMinutes       int               `yaml:"timeoutInMinutes,omitempty"`
	CancelTimeoutInMinutes int               `yaml:"cancelTimeoutInMinutes,omitempty"`
	Variables              *Variables        `yaml:"variables,omitempty"`
	Pool                   *Pool             `yaml:"pool,omitempty"`
	Container              *JobContainer     `yaml:"container,omitempty"`
	Services               map[string]string `yaml:"services,omitempty"`
	Workspace              *Workspace        `yaml:"workspace,omitempty"`
	Uses                   *Uses             `yaml:"uses,omitempty"`
	Steps                  *Steps            `yaml:"steps,omitempty"`
	TemplateContext        map[string]any    `yaml:"templateContext,omitempty"`
}

type CIJob struct {
	Job           string       `yaml:"job,omitempty"`
	Strategy      *JobStrategy `yaml:"strategy,omitempty"`
	BaseJob       `yaml:",inline"`
	FileReference *models.FileReference
}

type DeploymentEnvironment struct {
	Name          string `yaml:"name,omitempty"`
	ResourceName  string `yaml:"resourceName,omitempty"`
	ResourceId    string `yaml:"resourceId,omitempty"`
	ResourceType  string `yaml:"resourceType,omitempty"`
	Tags          string `yaml:"tags,omitempty"`
	FileReference *models.FileReference
}

func (de *DeploymentEnvironment) UnmarshalYAML(node *yaml.Node) error {
	de.FileReference = loadersUtils.GetFileReference(node)
	if node.Tag == consts.StringTag {
		de.Name = node.Value
		return nil
	}

	return node.Decode(&de)
}

type DeploymentJob struct {
	Deployment    string                 `yaml:"deployment,omitempty"`
	Strategy      *DeploymentStrategy    `yaml:"strategy,omitempty"`
	Environment   *DeploymentEnvironment `yaml:"environment,omitempty"`
	BaseJob       `yaml:",inline"`
	FileReference *models.FileReference
}

type TemplateJob struct {
	Template      `yaml:",inline"`
	FileReference *models.FileReference
}

type Jobs struct {
	CIJobs         *[]CIJob
	DeploymentJobs *[]DeploymentJob
	TemplateJobs   *[]TemplateJob
	FileReference  *models.FileReference
}

func (j *Jobs) UnmarshalYAML(node *yaml.Node) error {
	var ciJobs []CIJob
	var deploymentJobs []DeploymentJob
	var templateJobs []TemplateJob

	for _, jobNode := range node.Content {
		switch getJobType(jobNode) {
		case CIJobType:
			job, err := parseCIJob(jobNode)
			if err != nil {
				return err
			}
			ciJobs = append(ciJobs, job)
		case DeploymentJobType:
			job, err := parseDeploymentJob(jobNode)
			if err != nil {
				return err
			}
			deploymentJobs = append(deploymentJobs, job)
		case TemplateJobType:
			job, err := parseTemplateJob(jobNode)
			if err != nil {
				return err
			}
			templateJobs = append(templateJobs, job)
		}
	}

	*j = Jobs{
		CIJobs:         &ciJobs,
		DeploymentJobs: &deploymentJobs,
		TemplateJobs:   &templateJobs,
		FileReference:  loadersUtils.GetFileReference(node),
	}
	return nil
}

func parseCIJob(node *yaml.Node) (CIJob, error) {
	var job CIJob
	if err := node.Decode(&job); err != nil {
		return job, err
	}
	job.FileReference = loadersUtils.GetFileReference(node)
	return job, nil
}

func parseDeploymentJob(node *yaml.Node) (DeploymentJob, error) {
	var job DeploymentJob
	if err := node.Decode(&job); err != nil {
		return job, err
	}
	job.FileReference = loadersUtils.GetFileReference(node)
	return job, nil
}

func parseTemplateJob(node *yaml.Node) (TemplateJob, error) {
	var job TemplateJob
	if err := node.Decode(&job); err != nil {
		return job, err
	}
	job.FileReference = loadersUtils.GetFileReference(node)
	return job, nil
}

func getJobType(job *yaml.Node) JobType {
	for _, node := range job.Content {
		if node.Tag == consts.StringTag {
			switch JobType(node.Value) {
			case CIJobType:
				return CIJobType
			case DeploymentJobType:
				return DeploymentJobType
			case TemplateJobType:
				return TemplateJobType
			}
		}
	}
	return CIJobType
}
