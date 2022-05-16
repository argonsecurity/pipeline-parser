package utils

import "regexp"

func AnyMatch(regexes []*regexp.Regexp, s *string) bool {
	if s == nil {
		return false
	}

	for _, regexp := range regexes {
		if regexp.MatchString(*s) {
			return true
		}
	}
	return false
}
