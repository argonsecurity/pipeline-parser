package config

import (
	"regexp"
)

var (
	commonBuildShellRegexes = []*regexp.Regexp{
		regexp.MustCompile(`docker build`),
		regexp.MustCompile(`docker-compose build`),
		regexp.MustCompile(`npm build`),
		regexp.MustCompile(`yarn build`),
		regexp.MustCompile(`go build`),
	}

	commonTestShellRegexes = []*regexp.Regexp{
		regexp.MustCompile(`go test`),
		regexp.MustCompile(`npm test`),
		regexp.MustCompile(`yarn test`),
	}

	commonBuildNameRegexes = []*regexp.Regexp{
		regexp.MustCompile(`build`),
		regexp.MustCompile(`build-*`),
	}

	commonTestNameRegexes = []*regexp.Regexp{
		regexp.MustCompile(`tests?`),
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
