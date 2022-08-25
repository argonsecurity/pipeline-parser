package triggers

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/job"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestParseControls(t *testing.T) {
	testCases := []struct {
		name              string
		controls          *job.Controls
		isDeny            bool
		expectedCondition *models.Condition
	}{
		{
			name:              "Controls is nil",
			controls:          nil,
			expectedCondition: nil,
		},
		{
			name:     "Controls is empty",
			controls: &job.Controls{},
			expectedCondition: &models.Condition{
				Allow: utils.GetPtr(true),
				Branches: &models.Filter{
					AllowList: []string{},
				},
				Events: []models.EventType{},
			},
		},
		{
			name: "Controls with data + isDeny = false",
			controls: &job.Controls{
				Refs: []string{
					"master",
					"main",
					"test",
					"pushes",
					"merge_requests",
					"api",
					"branches",
				},
				Changes: []string{
					"/path/to/file",
					"/path/to/another/file",
				},
				Variables: []string{
					`$VAR1 == "VALUE1"`,
					`$VAR2 == /VALUE2/`,
				},
			},
			isDeny: false,
			expectedCondition: &models.Condition{
				Allow: utils.GetPtr(true),
				Branches: &models.Filter{
					AllowList: []string{
						"master",
						"main",
						"test",
					},
				},
				Paths: &models.Filter{
					AllowList: []string{
						"/path/to/file",
						"/path/to/another/file",
					},
				},
				Events: []models.EventType{
					models.PushEvent,
					models.PullRequestEvent,
					models.EventType("api"),
					models.EventType("branches"),
				},
				Variables: map[string]string{
					"$VAR1": `"VALUE1"`,
					"$VAR2": "/VALUE2/",
				},
			},
		},
		{
			name: "Controls with data + isDeny = true",
			controls: &job.Controls{
				Refs: []string{
					"master",
					"main",
					"test",
					"pushes",
					"merge_requests",
					"api",
					"branches",
				},
				Changes: []string{
					"/path/to/file",
					"/path/to/another/file",
				},
				Variables: []string{
					`$VAR1 == "VALUE1"`,
					`$VAR2 == /VALUE2/`,
				},
			},
			isDeny: true,
			expectedCondition: &models.Condition{
				Allow: utils.GetPtr(false),
				Branches: &models.Filter{
					DenyList: []string{
						"master",
						"main",
						"test",
					},
				},
				Paths: &models.Filter{
					DenyList: []string{
						"/path/to/file",
						"/path/to/another/file",
					},
				},
				Events: []models.EventType{
					models.PushEvent,
					models.PullRequestEvent,
					models.EventType("api"),
					models.EventType("branches"),
				},
				Variables: map[string]string{
					"$VAR1": `"VALUE1"`,
					"$VAR2": "/VALUE2/",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := ParseControls(testCase.controls, testCase.isDeny)

			testutils.DeepCompare(t, testCase.expectedCondition, got)

		})
	}
}

func TestParseVariables(t *testing.T) {
	testCases := []struct {
		name              string
		expressions       []string
		expectedVariables map[string]string
	}{
		{
			name:              "Expressions is nil",
			expressions:       nil,
			expectedVariables: nil,
		},
		{
			name:              "Expressions is empty",
			expressions:       []string{},
			expectedVariables: map[string]string{},
		},
		{
			name: "Expressions with data",
			expressions: []string{
				`$VAR1 == "VALUE1"`,
				`$VAR2 == /VALUE2/`,
			},
			expectedVariables: map[string]string{
				"$VAR1": `"VALUE1"`,
				"$VAR2": "/VALUE2/",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseVariables(testCase.expressions)

			testutils.DeepCompare(t, testCase.expectedVariables, got)

		})
	}
}

func TestGetBranchesAndEvents(t *testing.T) {
	testCases := []struct {
		name             string
		refs             []string
		expectedBranches []string
		expectedEvents   []models.EventType
	}{
		{
			name:             "Refs is nil",
			refs:             nil,
			expectedBranches: []string{},
			expectedEvents:   []models.EventType{},
		},
		{
			name:             "Refs is empty",
			refs:             []string{},
			expectedBranches: []string{},
			expectedEvents:   []models.EventType{},
		},
		{
			name: "Refs contains only ref keywords",
			refs: []string{
				"api",
				"branches",
			},
			expectedBranches: []string{},
			expectedEvents: []models.EventType{
				models.EventType("api"),
				models.EventType("branches"),
			},
		},
		{
			name: "Refs contains only ref events",
			refs: []string{
				"pushes",
				"merge_requests",
			},
			expectedBranches: []string{},
			expectedEvents: []models.EventType{
				models.PushEvent,
				models.PullRequestEvent,
			},
		},
		{
			name: "Refs contains only branches",
			refs: []string{
				"master",
				"main",
				"test",
			},
			expectedBranches: []string{
				"master",
				"main",
				"test",
			},
			expectedEvents: []models.EventType{},
		},
		{
			name: "Refs contains all ref types",
			refs: []string{
				"master",
				"main",
				"test",
				"pushes",
				"merge_requests",
				"api",
				"branches",
			},
			expectedBranches: []string{
				"master",
				"main",
				"test",
			},
			expectedEvents: []models.EventType{
				models.PushEvent,
				models.PullRequestEvent,
				models.EventType("api"),
				models.EventType("branches"),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			gotBranches, gotEvents := getBranchesAndEvents(testCase.refs)

			testutils.DeepCompare(t, testCase.expectedBranches, gotBranches)
			testutils.DeepCompare(t, testCase.expectedEvents, gotEvents)

		})
	}
}

func TestGenerateFilter(t *testing.T) {
	testCases := []struct {
		name           string
		list           []string
		isDeny         bool
		expectedFilter *models.Filter
	}{
		{
			name:           "List is nil",
			list:           nil,
			expectedFilter: nil,
		},
		{
			name:   "IsDeny is true",
			list:   []string{"a", "b", "c"},
			isDeny: true,
			expectedFilter: &models.Filter{
				DenyList: []string{"a", "b", "c"},
			},
		},
		{
			name:   "IsDeny is false",
			list:   []string{"a", "b", "c"},
			isDeny: false,
			expectedFilter: &models.Filter{
				AllowList: []string{"a", "b", "c"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := generateFilter(testCase.list, testCase.isDeny)

			testutils.DeepCompare(t, testCase.expectedFilter, got)
		})
	}
}
