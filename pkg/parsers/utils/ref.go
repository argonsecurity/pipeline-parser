package utils

import (
	"regexp"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

var (
	sha1Regex   = regexp.MustCompile(`[0-9a-fA-F]{40}`)
	semverRegex = regexp.MustCompile(`v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`)
)

func DetectVersionType(version string) models.VersionType {
	versionType := models.None

	if version != "" {
		if sha1Regex.MatchString(version) {
			versionType = models.CommitSHA
		} else if semverRegex.MatchString(version) {
			versionType = models.TagVersion
		} else if version == "latest" {
			versionType = models.Latest
		} else {
			versionType = models.BranchVersion
		}
	}

	return versionType
}
