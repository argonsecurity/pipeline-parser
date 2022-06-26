package gitlab

import (
	"testing"

	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/r3labs/diff/v3"
	"github.com/stretchr/testify/assert"
)

func TestParseDefaults(t *testing.T) {
	testCases := []struct {
		name                  string
		gitlabCIConfiguration *gitlabModels.GitlabCIConfiguration
		expectedDefaults      *models.Defaults
	}{
		{
			name:                  "Gitlab CI config is empty",
			gitlabCIConfiguration: &gitlabModels.GitlabCIConfiguration{},
			expectedDefaults:      &models.Defaults{},
		},
		// {
		// 	name:                  "Gitlab CI config with defaults data",
		// 	gitlabCIConfiguration: &gitlabModels.GitlabCIConfiguration{
		// 		Variables: &common.EnvironmentVariablesRef{

		// 		},
		// 	},
		// 	expectedDefaults:      &models.Defaults{},
		// },
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseDefaults(testCase.gitlabCIConfiguration)

			changelog, err := diff.Diff(testCase.expectedDefaults, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0)
		})
	}

}
