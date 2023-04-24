package models

import (
	"fmt"

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
	DependsOn              *DependsOn        `yaml:"dependsOn,omitempty"`
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

type DeploymentEnvironmentRef struct {
	DeploymentEnvironment *DeploymentEnvironment `yaml:"deploymentEnvironment,omitempty"`
	FileReference         *models.FileReference
}

type DeploymentEnvironment struct {
	Name         string `yaml:"name,omitempty"`
	ResourceName string `yaml:"resourceName,omitempty"`
	ResourceId   string `yaml:"resourceId,omitempty"`
	ResourceType string `yaml:"resourceType,omitempty"`
	Tags         string `yaml:"tags,omitempty"`
}

type DeploymentJob struct {
	Deployment    string                    `yaml:"deployment,omitempty"`
	Strategy      *DeploymentStrategy       `yaml:"strategy,omitempty"`
	Environment   *DeploymentEnvironmentRef `yaml:"environment,omitempty"`
	BaseJob       `yaml:",inline"`
	FileReference *models.FileReference
}

type TemplateJob struct {
	Template      `yaml:",inline"`
	FileReference *models.FileReference
}

type Jobs struct {
	CIJobs         []*CIJob
	DeploymentJobs []*DeploymentJob
	TemplateJobs   []*TemplateJob
	FileReference  *models.FileReference
}

func (der *DeploymentEnvironmentRef) UnmarshalYAML(node *yaml.Node) error {
	der.FileReference = loadersUtils.GetFileReference(node)
	if node.Tag == consts.StringTag {
		der.DeploymentEnvironment = &DeploymentEnvironment{
			Name: node.Value,
		}
		return nil
	}

	return node.Decode(&der.DeploymentEnvironment)
}

func (j *Jobs) UnmarshalYAML(node *yaml.Node) error {
	var ciJobs []*CIJob
	var deploymentJobs []*DeploymentJob
	var templateJobs []*TemplateJob

	for _, jobNode := range node.Content {
		if len(jobNode.Content) == 0 {
			consts.NewErrInvalidYaml("job is empty")
		}

		if jobNode.Tag == consts.StringTag {
			var job TemplateJob
			job.Template = Template{
				Template: jobNode.Value,
			}
			job.FileReference = loadersUtils.GetFileReference(jobNode)
			templateJobs = append(templateJobs, &job)
			continue
		}

		switch JobType(jobNode.Content[0].Value) {
		case CIJobType:
			var job CIJob
			if err := jobNode.Decode(&job); err != nil {
				return err
			}
			job.FileReference = loadersUtils.GetFileReference(jobNode)
			ciJobs = append(ciJobs, &job)
		case DeploymentJobType:
			var job DeploymentJob
			if err := jobNode.Decode(&job); err != nil {
				return err
			}
			job.FileReference = loadersUtils.GetFileReference(jobNode)
			deploymentJobs = append(deploymentJobs, &job)
		case TemplateJobType:
			var job TemplateJob
			if err := jobNode.Decode(&job); err != nil {
				return err
			}
			job.FileReference = loadersUtils.GetFileReference(jobNode)
			templateJobs = append(templateJobs, &job)
		default:
			consts.NewErrInvalidYaml(fmt.Sprintf("job must start with one of %s, %s or %s", CIJobType, DeploymentJobType, TemplateJobType))
		}
	}

	*j = Jobs{
		CIJobs:         ciJobs,
		DeploymentJobs: deploymentJobs,
		TemplateJobs:   templateJobs,
		FileReference:  loadersUtils.GetFileReference(node),
	}
	return nil
}
