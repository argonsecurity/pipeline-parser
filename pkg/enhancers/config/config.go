package config

import "regexp"

type ObjectiveConfiguration struct {
	Tasks        []string
	Names        []*regexp.Regexp
	ShellRegexes []*regexp.Regexp
}

type EnhancementConfiguration struct {
	Build  ObjectiveConfiguration
	Test   ObjectiveConfiguration
	Deploy ObjectiveConfiguration
}
