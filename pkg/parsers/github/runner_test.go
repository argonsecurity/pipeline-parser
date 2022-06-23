package github

import (
	"testing"

	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseRunsOnToRunner(t *testing.T) {
	testCases := []struct {
		name           string
		runsOn         *githubModels.RunsOn
		expectedRunner *models.Runner
	}{
		{
			name:           "runsOn nil",
			runsOn:         nil,
			expectedRunner: nil,
		},
		{
			name: "Full runsOn",
			runsOn: &githubModels.RunsOn{
				OS:         utils.GetPtr("linux"),
				Arch:       utils.GetPtr("amd64"),
				SelfHosted: true,
				Tags:       []string{"tag1", "tag2"},
				FileReference: &models.FileReference{
					StartRef: &models.FileLocation{
						Line:   1,
						Column: 1,
					},
					EndRef: &models.FileLocation{
						Line:   2,
						Column: 2,
					},
				},
			},
			expectedRunner: &models.Runner{
				OS:         utils.GetPtr("linux"),
				Arch:       utils.GetPtr("amd64"),
				SelfHosted: utils.GetPtr(true),
				Labels:     &[]string{"tag1", "tag2"},
				FileReference: &models.FileReference{
					StartRef: &models.FileLocation{
						Line:   1,
						Column: 1,
					},
					EndRef: &models.FileLocation{
						Line:   2,
						Column: 2,
					},
				},
			},
		},
		{
			name:   "Empty runsOn",
			runsOn: &githubModels.RunsOn{},
			expectedRunner: &models.Runner{
				SelfHosted: utils.GetPtr(false),
				Labels:     utils.GetPtr[[]string](nil),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			runner := parseRunsOnToRunner(testCase.runsOn)
			assert.Equal(t, testCase.expectedRunner, runner, testCase.name)
		})
	}
}
