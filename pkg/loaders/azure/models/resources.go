package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Build struct {
	Build         string `yaml:"build,omitempty"`
	Type          string `yaml:"type,omitempty"`
	Connection    string `yaml:"connection,omitempty"`
	Source        string `yaml:"source,omitempty"`
	Version       string `yaml:"version,omitempty"`
	Branch        string `yaml:"branch,omitempty"`
	Trigger       string `yaml:"trigger,omitempty"`
	FileReference *models.FileReference
}

func (b *Build) UnmarshalYAML(node *yaml.Node) error {
	b.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&b)
}

type MountReadOnly struct {
	Work     bool `yaml:"work,omitempty"`
	External bool `yaml:"external,omitempty"`
	Tools    bool `yaml:"tools,omitempty"`
	Tasks    bool `yaml:"tasks,omitempty"`
}

type JobContainer struct {
	Image           string                   `yaml:"image,omitempty"`
	Endpoint        string                   `yaml:"endpoint,omitempty"`
	Env             *EnvironmentVariablesRef `yaml:"env,omitempty"`
	MapDockerSocket bool                     `yaml:"mapDockerSocket,omitempty"`
	Options         string                   `yaml:"options,omitempty"`
	Ports           []string                 `yaml:"ports,omitempty"`
	Volumes         []string                 `yaml:"volumes,omitempty"`
	MountReadOnly   *MountReadOnly           `yaml:"mountReadOnly,omitempty"`
}

type ResourceContainer struct {
	Container         string      `yaml:"container,omitempty"`
	Type              string      `yaml:"type,omitempty"`
	Trigger           *TriggerRef `yaml:"trigger,omitempty"`
	AzureSubscription string      `yaml:"azureSubscription,omitempty"`
	ResourceGroup     string      `yaml:"resourceGroup,omitempty"`
	Registry          string      `yaml:"registry,omitempty"`
	Repository        string      `yaml:"repository,omitempty"`
	JobContainer      `yaml:",inline"`
	FileReference     *models.FileReference
}

func (c *ResourceContainer) UnmarshalYAML(node *yaml.Node) error {
	c.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&c)
}

type ResourcePipeline struct {
	Pipeline string      `yaml:"pipeline,omitempty"`
	Project  string      `yaml:"project,omitempty"`
	Source   string      `yaml:"source,omitempty"`
	Version  string      `yaml:"version,omitempty"`
	Branch   string      `yaml:"branch,omitempty"`
	Tags     []string    `yaml:"tags,omitempty"`
	Trigger  *TriggerRef `yaml:"trigger,omitempty"`

	FileReference *models.FileReference
}

func (p *ResourcePipeline) UnmarshalYAML(node *yaml.Node) error {
	p.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&p)
}

type Repository struct {
	Repository    string      `yaml:"repository,omitempty"`
	Endpoint      string      `yaml:"endpoint,omitempty"`
	Trigger       *TriggerRef `yaml:"trigger,omitempty"`
	Name          string      `yaml:"name,omitempty"`
	Type          string      `yaml:"type,omitempty"`
	Ref           string      `yaml:"ref,omitempty"`
	FileReference *models.FileReference
}

func (r *Repository) UnmarshalYAML(node *yaml.Node) error {
	r.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&r)
}

type Path struct {
	Path  string `yaml:"path,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type Webhook struct {
	Webhook       string `yaml:"webhook,omitempty"`
	Connection    string `yaml:"connection,omitempty"`
	Type          string `yaml:"type,omitempty"`
	Filters       []Path
	FileReference *models.FileReference
}

func (w *Webhook) UnmarshalYAML(node *yaml.Node) error {
	w.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&w)
}

type Package struct {
	Package       string `yaml:"package,omitempty"`
	Type          string `yaml:"type,omitempty"`
	Connection    string `yaml:"connection,omitempty"`
	Name          string `yaml:"name,omitempty"`
	Version       string `yaml:"version,omitempty"`
	Tag           string `yaml:"tag,omitempty"`
	Trigger       string `yaml:"trigger,omitempty"`
	FileReference *models.FileReference
}

func (p *Package) UnmarshalYAML(node *yaml.Node) error {
	p.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&p)
}

type Resources struct {
	Builds        *[]Build             `yaml:"builds,omitempty"`
	Containers    *[]ResourceContainer `yaml:"containers,omitempty"`
	Pipelines     *[]ResourcePipeline  `yaml:"pipelines,omitempty"`
	Repositories  *[]Repository        `yaml:"repositories,omitempty"`
	Webhooks      *[]Webhook           `yaml:"webhooks,omitempty"`
	Packages      *[]Package           `yaml:"packages,omitempty"`
	FileReference *models.FileReference
}

func (r *Resources) UnmarshalYAML(node *yaml.Node) error {
	r.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&r)
}
