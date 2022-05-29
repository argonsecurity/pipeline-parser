package models

const (
	DockerRunnerType RunnerType = "docker"
	VmRunnerType     RunnerType = "vm"
	ServerRunnerType RunnerType = "server"

	WindowsOS OS = "windows"
	LinuxOS   OS = "linux"
	MacOS     OS = "macos"
)

type OS string
type RunnerType string

type DockerMetadata struct {
	Image                 *string `json:"image,omitempty"`
	Label                 *string `json:"label,omitempty"`
	RegistryURL           *string `json:"registry_url,omitempty"`
	RegistryCredentialsID *string `json:"registry_credentials_id,omitempty"`
}

type Runner struct {
	Type           *string         `json:"type,omitempty"`
	Labels         *[]string       `json:"labels,omitempty"`
	OS             *string         `json:"os,omitempty"`
	Arch           *string         `json:"arch,omitempty"`
	SelfHosted     *bool           `json:"self_hosted,omitempty"`
	DockerMetadata *DockerMetadata `json:"docker_metadata,omitempty"`
	FileReference  *FileReference  `json:"file_reference,omitempty"`
}
