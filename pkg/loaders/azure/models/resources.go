package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

type Build struct {
	Build      string `yaml:"build,omitempty"`
	Type       string `yaml:"type,omitempty"`
	Connection string `yaml:"connection,omitempty"`
	Source     string `yaml:"source,omitempty"`
	Version    string `yaml:"version,omitempty"`
	Branch     string `yaml:"branch,omitempty"`
	Trigger    string `yaml:"trigger,omitempty"`
}

type BuildRef struct {
	Build         *Build
	FileReference *models.FileReference
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
}

type ResourceContainerRef struct {
	ResourceContainer *ResourceContainer
	FileReference     *models.FileReference
}

type ResourcePipeline struct {
	Pipeline string      `yaml:"pipeline,omitempty"`
	Project  string      `yaml:"project,omitempty"`
	Source   string      `yaml:"source,omitempty"`
	Version  string      `yaml:"version,omitempty"`
	Branch   string      `yaml:"branch,omitempty"`
	Tags     []string    `yaml:"tags,omitempty"`
	Trigger  *TriggerRef `yaml:"trigger,omitempty"`
}

type ResourcePipelineRef struct {
	ResourcePipeline *ResourcePipeline
	FileReference    *models.FileReference
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

type RepositoryRef struct {
	Repository    *Repository
	FileReference *models.FileReference
}

type Path struct {
	Path  string `yaml:"path,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type Webhook struct {
	Webhook    string `yaml:"webhook,omitempty"`
	Connection string `yaml:"connection,omitempty"`
	Type       string `yaml:"type,omitempty"`
	Filters    []Path
}

type WebhookRef struct {
	Webhook       *Webhook
	FileReference *models.FileReference
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

type PackageRef struct {
	Package       *Package
	FileReference *models.FileReference
}

type Resources struct {
	Builds        []*BuildRef             `yaml:"builds,omitempty"`
	Containers    []*ResourceContainerRef `yaml:"containers,omitempty"`
	Pipelines     []*ResourcePipelineRef  `yaml:"pipelines,omitempty"`
	Repositories  []*RepositoryRef        `yaml:"repositories,omitempty"`
	Webhooks      []*WebhookRef           `yaml:"webhooks,omitempty"`
	Packages      []*PackageRef           `yaml:"packages,omitempty"`
	FileReference *models.FileReference
}

func (br *BuildRef) UnmarshalYAML(node *yaml.Node) error {
	br.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&br.Build)
}

func (rcf *ResourceContainerRef) UnmarshalYAML(node *yaml.Node) error {
	rcf.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&rcf.ResourceContainer)
}

func (rpr *ResourcePipelineRef) UnmarshalYAML(node *yaml.Node) error {
	rpr.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&rpr.ResourcePipeline)
}

func (rr *RepositoryRef) UnmarshalYAML(node *yaml.Node) error {
	rr.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&rr.Repository)
}

func (wr *WebhookRef) UnmarshalYAML(node *yaml.Node) error {
	wr.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&wr.Webhook)
}

func (pr *PackageRef) UnmarshalYAML(node *yaml.Node) error {
	pr.FileReference = loadersUtils.GetFileReference(node)
	return node.Decode(&pr.Package)
}

func (r *Resources) UnmarshalYAML(node *yaml.Node) error {
	r.FileReference = loadersUtils.GetFileReference(node)
	r.FileReference.StartRef.Line--      // The "resources" node is not accessible, this is a patch
	r.FileReference.StartRef.Column -= 2 // The "resources" node is not accessible, this is a patch
	return loadersUtils.IterateOnMap(node, func(key string, value *yaml.Node) error {
		switch key {
		case "builds":
			var builds []*BuildRef
			if err := value.Decode(&builds); err != nil {
				return err
			}
			r.Builds = builds
		case "containers":
			var containers []*ResourceContainerRef
			if err := value.Decode(&containers); err != nil {
				return err
			}
			r.Containers = containers
		case "pipelines":
			var pipelines []*ResourcePipelineRef
			if err := value.Decode(&pipelines); err != nil {
				return err
			}
			r.Pipelines = pipelines
		case "repositories":
			var repositories []*RepositoryRef
			if err := value.Decode(&repositories); err != nil {
				return err
			}
			r.Repositories = repositories
		case "webhooks":
			var webhooks []*WebhookRef
			if err := value.Decode(&webhooks); err != nil {
				return err
			}
			r.Webhooks = webhooks
		case "packages":
			var packages []*PackageRef
			if err := value.Decode(&packages); err != nil {
				return err
			}
			r.Packages = packages
		}
		return nil
	}, "Resources")
}
