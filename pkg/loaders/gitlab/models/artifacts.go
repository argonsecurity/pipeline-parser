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

	Codequality        interface{} `yaml:"codequality,omitempty"`
	ContainerScanning  interface{} `yaml:"container_scanning,omitempty"`
	Dast               interface{} `yaml:"dast,omitempty"`
	DependencyScanning interface{} `yaml:"dependency_scanning,omitempty"`
	Dotenv             interface{} `yaml:"dotenv,omitempty"`
	Junit              interface{} `yaml:"junit,omitempty"`
	LicenseManagement  interface{} `yaml:"license_management,omitempty"`
	LicenseScanning    interface{} `yaml:"license_scanning,omitempty"`
	Lsif               interface{} `yaml:"lsif,omitempty"`
	Metrics            interface{} `yaml:"metrics,omitempty"`
	Performance        interface{} `yaml:"performance,omitempty"`
	Requirements       interface{} `yaml:"requirements,omitempty"`
	Sast               interface{} `yaml:"sast,omitempty"`
	SecretDetection    interface{} `yaml:"secret_detection,omitempty"`
	Terraform          interface{} `yaml:"terraform,omitempty"`
}

type CoverageReport struct {
	CoverageFormat interface{} `yaml:"coverage_format,omitempty"`
	Path           string      `yaml:"path,omitempty"`
}
