package models

// import "gopkg.in/yaml.v3"

type MountReadOnly struct {
	Work     bool `yaml:"work,omitempty"`
	External bool `yaml:"external,omitempty"`
	Tools    bool `yaml:"tools,omitempty"`
	Tasks    bool `yaml:"tasks,omitempty"`
}

type ContainerTrigger struct {
	Trigger *Trigger `yaml:"trigger,omitempty"`
}

// type Build struct {
// 	Build      string `yaml:"build,omitempty"`
// 	Type       string `yaml:"type,omitempty"`
// 	Connection string `yaml:"connection,omitempty"`
// 	Source     string `yaml:"source,omitempty"`
// 	Version    string `yaml:"version,omitempty"`
// 	Branch     string `yaml:"branch,omitempty"`
// 	Trigger    string `yaml:"trigger,omitempty"`
// }

type Container struct {
	Container         string                   `yaml:"container,omitempty"`
	Image             string                   `yaml:"image,omitempty"`
	Type              string                   `yaml:"type,omitempty"`
	Trigger           *ContainerTrigger        `yaml:"trigger,omitempty"`
	Endpoint          string                   `yaml:"endpoint,omitempty"`
	Env               *EnvironmentVariablesRef `yaml:"env,omitempty"`
	MapDockerSocket   bool                     `yaml:"mapDockerSocket,omitempty"`
	Options           string                   `yaml:"options,omitempty"`
	Ports             []string                 `yaml:"ports,omitempty"`
	Volumes           []string                 `yaml:"volumes,omitempty"`
	MountReadOnly     *MountReadOnly           `yaml:"mountReadOnly,omitempty"`
	AzureSubscription string                   `yaml:"azureSubscription,omitempty"`
	ResourceGroup     string                   `yaml:"resourceGroup,omitempty"`
	Registry          string                   `yaml:"registry,omitempty"`
	Repository        string                   `yaml:"repository,omitempty"`
}

// type Package struct {
// 	Package    string `yaml:"package,omitempty"`
// 	Type       string `yaml:"type,omitempty"`
// 	Connection string `yaml:"connection,omitempty"`
// 	Name       string `yaml:"name,omitempty"`
// 	Version    string `yaml:"version,omitempty"`
// 	Tag        string `yaml:"tag,omitempty"`
// 	Trigger    string `yaml:"trigger,omitempty"`
// }

// type ResourcePipeline struct {
// 	Pipeline string   `yaml:"pipeline,omitempty"`
// 	Project  string   `yaml:"project,omitempty"`
// 	Source   string   `yaml:"source,omitempty"`
// 	Version  string   `yaml:"version,omitempty"`
// 	Branch   string   `yaml:"branch,omitempty"`
// 	Tags     []string `yaml:"tags,omitempty"`
// }

// type Resources struct {
// 	Builds     *[]Build     `yaml:"builds,omitempty"`
// 	Containers *[]Container `yaml:"containers,omitempty"`
// 	Packages   *[]Package   `yaml:"packages,omitempty"`
// }

// func (r *Resources) UnmarshalYAML(node *yaml.Node) error {
// 	return node.Decode(r)
// }

// func (r *ContainerTrigger) UnmarshalYAML(node *yaml.Node) error {
// 	return node.Decode(r)
// }
