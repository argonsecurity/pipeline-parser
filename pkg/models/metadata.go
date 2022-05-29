package models

type Metadata struct {
	Build  bool `json:"build,omitempty"`
	Test   bool `json:"test,omitempty"`
	Deploy bool `json:"deploy,omitempty"`
}
