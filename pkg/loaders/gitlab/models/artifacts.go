package models

type Artifacts struct {
	Exclude   []string `yaml:"exclude,omitempty"`
	ExpireIn  string   `yaml:"expire_in,omitempty"`
	ExposeAs  string   `yaml:"expose_as,omitempty"`
	Name      string   `yaml:"name,omitempty"`
	Paths     []string `yaml:"paths,omitempty"`
	Reports   *Reports `yaml:"reports,omitempty"`
	Untracked bool     `yaml:"untracked,omitempty"`
	When      string   `yaml:"when,omitempty"`
}

type Reports struct {
	CoverageReport *CoverageReport `yaml:"coverage_report,omitempty"`

	Codequality        any `yaml:"codequality,omitempty"`
	ContainerScanning  any `yaml:"container_scanning,omitempty"`
	Dast               any `yaml:"dast,omitempty"`
	DependencyScanning any `yaml:"dependency_scanning,omitempty"`
	Dotenv             any `yaml:"dotenv,omitempty"`
	Junit              any `yaml:"junit,omitempty"`
	LicenseManagement  any `yaml:"license_management,omitempty"`
	LicenseScanning    any `yaml:"license_scanning,omitempty"`
	Lsif               any `yaml:"lsif,omitempty"`
	Metrics            any `yaml:"metrics,omitempty"`
	Performance        any `yaml:"performance,omitempty"`
	Requirements       any `yaml:"requirements,omitempty"`
	Sast               any `yaml:"sast,omitempty"`
	SecretDetection    any `yaml:"secret_detection,omitempty"`
	Terraform          any `yaml:"terraform,omitempty"`
}

type CoverageReport struct {
	CoverageFormat any    `yaml:"coverage_format,omitempty"`
	Path           string `yaml:"path,omitempty"`
}
