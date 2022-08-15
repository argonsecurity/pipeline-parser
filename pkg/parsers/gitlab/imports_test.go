package gitlab

import (
	"testing"

	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
)

func TestParseImports(t *testing.T) {
	testCases := []struct {
		name            string
		include         *gitlabModels.Include
		expectedImports []string
	}{
		{
			name:            "Include is nil",
			include:         nil,
			expectedImports: nil,
		},
		{
			name:            "Include is empty",
			include:         &gitlabModels.Include{},
			expectedImports: []string{},
		},
		{
			name: "One include with all data",
			include: &gitlabModels.Include{
				{
					Template: "template1",
					Local:    "local1",
					Remote:   "remote1",
				},
			},
			expectedImports: []string{"template1", "local1", "remote1"},
		},
		{
			name: "Some includes with all data",
			include: &gitlabModels.Include{
				{
					Template: "template1",
					Local:    "local1",
					Remote:   "remote1",
				},
				{
					Template: "template2",
					Local:    "local2",
					Remote:   "remote2",
				},
			},
			expectedImports: []string{"template1", "local1", "remote1", "template2", "local2", "remote2"},
		},
		{
			name:            "Include is empty",
			include:         &gitlabModels.Include{},
			expectedImports: []string{},
		},
		{
			name: "One include with partial data",
			include: &gitlabModels.Include{
				{
					Template: "template1",
					Local:    "local1",
				},
			},
			expectedImports: []string{"template1", "local1"},
		},
		{
			name: "Some includes with partial data",
			include: &gitlabModels.Include{
				{
					Template: "template1",
					Local:    "local1",
				},
				{
					Template: "template2",
					Remote:   "remote2",
				},
				{
					Local:  "local3",
					Remote: "remote3",
				},
			},
			expectedImports: []string{"template1", "local1", "template2", "remote2", "local3", "remote3"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseImports(testCase.include)

			testutils.DeepCompare(t, testCase.expectedImports, got)
		})
	}
}
