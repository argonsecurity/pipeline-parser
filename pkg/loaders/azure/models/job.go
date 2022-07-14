package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Jobs []Job

type DependsOn []string

func (n *DependsOn) UnmarshalYAML(node *yaml.Node) error {
	var tags []string
	var err error

	if node.Tag == consts.SequenceTag {
		if tags, err = loadersUtils.ParseYamlStringSequenceToSlice(node); err != nil {
			return err
		}
	} else if node.Tag == consts.StringTag {
		tags = []string{node.Value}
	} else {
		return consts.NewErrInvalidYamlTag(node.Tag)
	}

	*n = tags
	return nil
}

// type JobContainer struct {
// 	Image           string                   `yaml:"image,omitempty"`
// 	Endpoint        string                   `yaml:"endpoint,omitempty"`
// 	Env             *EnvironmentVariablesRef `yaml:"env,omitempty"`
// 	MapDockerSocket bool                     `yaml:"mapDockerSocket,omitempty"`
// 	Options         string                 `yaml:"options,omitempty"`
// 	Ports		   []string                 `yaml:"ports,omitempty"`
// 	Volumes		   []string                 `yaml:"volumes,omitempty"`
// 	MountReadOnly   *MountReadOnly           `yaml:"mountReadOnly,omitempty"`

// }

type Workspace struct {
	Clean string `yaml:"clean,omitempty"`
}

type Uses struct {
	Repositories []string `yaml:"repositories,omitempty"`
	Pools        []string `yaml:"pools,omitempty"`
}

type Job struct {
	Job                    string       `yaml:"job,omitempty"`
	DisplayName            string       `yaml:"displayName,omitempty"`
	DependsOn              []string     `yaml:"dependsOn,omitempty"`
	Condition              string       `yaml:"condition,omitempty"`
	ContinueOnError        bool         `yaml:"continueOnError,omitempty"`
	TimeoutInMinutes       int          `yaml:"timeoutInMinutes,omitempty"`
	CancelTimeoutInMinutes int          `yaml:"cancelTimeoutInMinutes,omitempty"`
	Variables              *Variables   `yaml:"variables,omitempty"`
	Strategy               *JobStrategy `yaml:"strategy,omitempty"`
	Pool                   *Pool        `yaml:"pool,omitempty"`
	// Container              *JobContainer     `yaml:"container,omitempty"`
	Services        map[string]string `yaml:"services,omitempty"`
	Workspace       *Workspace        `yaml:"workspace,omitempty"`
	Uses            *Uses             `yaml:"uses,omitempty"`
	Steps           *Steps            `yaml:"steps,omitempty"`
	TemplateContext map[string]any    `yaml:"templateContext,omitempty"`
	FileReference   *models.FileReference
}

func (j *Jobs) UnmarshalYAML(node *yaml.Node) error {
	var jobs []Job

	for _, jobNode := range node.Content {
		job, err := parseJob(jobNode)
		if err != nil {
			return err
		}
		jobs = append(jobs, job)
	}

	*j = jobs
	return nil
}

func parseJob(node *yaml.Node) (Job, error) {
	var job Job
	if err := node.Decode(&job); err != nil {
		return job, err
	}
	job.FileReference = loadersUtils.GetFileReference(node)
	return job, nil
}
