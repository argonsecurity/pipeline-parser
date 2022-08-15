package gitlab

import (
	"testing"

	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
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

			testutils.DeepCompare(t, testCase.expectedDefaults, got)
		})
	}

}
