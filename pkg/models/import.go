package models

type SourceType string

const (
	SourceTypeLocal  SourceType = "local"
	SourceTypeRemote SourceType = "remote"
)

type ImportSource struct {
	SCM             Platform   `json:"scm,omitempty"`
	Organization    *string    `json:"organization,omitempty"`
	Repository      *string    `json:"repository,omitempty"`
	Path            *string    `json:"path,omitempty"`
	Type            SourceType `json:"type,omitempty"`
	RepositoryAlias *string    `json:"alias,omitempty"`
	Reference       *string    `json:"reference,omitempty"`
}

type Import struct {
	Source        *ImportSource  `json:"source,omitempty"`
	Version       *string        `json:"version,omitempty"`
	VersionType   VersionType    `json:"version_type,omitempty"`
	Pipeline      *Pipeline      `json:"pipeline,omitempty"`
	Parameters    map[string]any `json:"parameters,omitempty"`
	Secrets       *SecretsRef    `json:"secrets,omitempty"`
	FileReference *FileReference `json:"file_reference,omitempty"`
}
