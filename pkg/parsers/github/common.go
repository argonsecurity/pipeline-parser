package github

import (
	"regexp"

	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

var (
	sha1Regex   = regexp.MustCompile(`[0-9a-fA-F]{40}`)
	semverRegex = regexp.MustCompile(`v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`)
)

func parseEnvironmentVariablesRef(envRef *githubModels.EnvironmentVariablesRef) *models.EnvironmentVariablesRef {
	if envRef == nil {
		return nil
	}

	return &models.EnvironmentVariablesRef{
		EnvironmentVariables: envRef.EnvironmentVariables,
		FileReference:        envRef.FileReference,
	}
}

func detectVersionType(version string) models.VersionType {
	versionType := models.None

	if version != "" {
		if sha1Regex.MatchString(version) {
			versionType = models.CommitSHA
		} else if semverRegex.MatchString(version) {
			versionType = models.TagVersion
		} else {
			versionType = models.BranchVersion
		}
	}

	return versionType
}
