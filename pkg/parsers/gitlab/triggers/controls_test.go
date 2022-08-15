package triggers

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
)

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
			name:              "Expressions is empty",
			expressions:       []string{},
			expectedVariables: map[string]string{},
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
