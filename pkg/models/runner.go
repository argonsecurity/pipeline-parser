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
	Image                 *string
	Label                 *string
	RegistryURL           *string
	RegistryCredentialsID *string
}

type Runner struct {
	Type           *string
	Labels         *[]string
	OS             *string
	Arch           *string
	SelfHosted     *bool
	DockerMetadata *DockerMetadata
	FileReference  *FileReference
}
