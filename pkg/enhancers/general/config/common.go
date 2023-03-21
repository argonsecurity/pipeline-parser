package config

import (
	"regexp"
)

var (
	commonBuildShellRegexes = []*regexp.Regexp{
		regexp.MustCompile(`docker build`),
		regexp.MustCompile(`docker-compose build`),
		regexp.MustCompile(`npm run build`),
		regexp.MustCompile(`yarn( run)? build`),
		regexp.MustCompile(`go build`),
	}

	commonTestShellRegexes = []*regexp.Regexp{
		regexp.MustCompile(`go test`),
		regexp.MustCompile(`npm run test`),
		regexp.MustCompile(`yarn( run)? test`),
	}

	commonBuildNameRegexes = []*regexp.Regexp{
		regexp.MustCompile(`(?i)build.*`),
	}

	commonTestNameRegexes = []*regexp.Regexp{
		regexp.MustCompile(`(?i)tests?`),
	}

	CommonConfiguration = &EnhancementConfiguration{
		Build: ObjectiveConfiguration{
			Names:        commonBuildNameRegexes,
			ShellRegexes: commonBuildShellRegexes,
		},
		Test: ObjectiveConfiguration{
			Names:        commonTestNameRegexes,
			ShellRegexes: commonTestShellRegexes,
		},
	}
)
