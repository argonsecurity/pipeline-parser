package models

type Clone struct {
	Depth   any   `yaml:"depth"`             // Depth of Git clones for all pipelines (supported only for Git repositories)
	Enabled *bool `yaml:"enabled,omitempty"` // Enables cloning of the repository
	LFS     *bool `yaml:"lfs,omitempty"`     // Enables the download of LFS files in the clone (supported only for Git repositories)
}
