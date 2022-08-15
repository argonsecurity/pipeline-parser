package triggers

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestParseConditionRules(t *testing.T) {
	testCases := []struct {
		name               string
		rules              *common.Rules
		expectedConditions []*models.Condition
	}{
		{
			name:               "Rules list is null",
			rules:              nil,
			expectedConditions: nil,
		},
		{
			name:               "Rules list is empty",
			rules:              &common.Rules{},
			expectedConditions: []*models.Condition{},
		},
		{
			name: "Rules list with one rule",
			rules: &common.Rules{
				RulesList: []*common.Rule{
					{
						If:      "condition",
						When:    never,
						Changes: []string{"a", "b", "c"},
						Exists:  []string{"d", "e", "f"},
						Variables: &common.EnvironmentVariablesRef{
							Variables: &common.Variables{
								"string": "string",
								"number": 1,
								"bool":   true,
							},
						},
					},
				},
			},
			expectedConditions: []*models.Condition{
				{
					Statement: "condition",
					Allow:     utils.GetPtr(false),
					Paths: &models.Filter{
						DenyList: []string{"a", "b", "c"},
					},
					Exists: &models.Filter{
						DenyList: []string{"d", "e", "f"},
					},
					Variables: map[string]string{
						"string": "string",
						"number": "1",
						"bool":   "true",
					},
				},
			},
		},
		{
			name: "Rules list with some rules",
			rules: &common.Rules{
				RulesList: []*common.Rule{
					{
						If:      "condition",
						When:    never,
						Changes: []string{"a", "b", "c"},
						Exists:  []string{"d", "e", "f"},
						Variables: &common.EnvironmentVariablesRef{
							Variables: &common.Variables{
								"string": "string",
								"number": 1,
								"bool":   true,
							},
						},
					},
					{
						If:      "condition",
						When:    always,
						Changes: []string{"a", "b", "c"},
						Exists:  []string{"d", "e", "f"},
						Variables: &common.EnvironmentVariablesRef{
							Variables: &common.Variables{
								"string": "string",
								"number": 1,
								"bool":   true,
							},
						},
					},
				},
			},
			expectedConditions: []*models.Condition{
				{
					Statement: "condition",
					Allow:     utils.GetPtr(false),
					Paths: &models.Filter{
						DenyList: []string{"a", "b", "c"},
					},
					Exists: &models.Filter{
						DenyList: []string{"d", "e", "f"},
					},
					Variables: map[string]string{
						"string": "string",
						"number": "1",
						"bool":   "true",
					},
				},
				{
					Statement: "condition",
					Allow:     utils.GetPtr(true),
					Paths: &models.Filter{
						AllowList: []string{"a", "b", "c"},
					},
					Exists: &models.Filter{
						AllowList: []string{"d", "e", "f"},
					},
					Variables: map[string]string{
						"string": "string",
						"number": "1",
						"bool":   "true",
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := ParseConditionRules(testCase.rules)

			testutils.DeepCompare(t, testCase.expectedConditions, got)

		})
	}
}

func TestParseConditionRule(t *testing.T) {
	testCases := []struct {
		name              string
		rule              *common.Rule
		expectedCondition *models.Condition
	}{
		{
			name: "Rule is empty",
			rule: &common.Rule{},
			expectedCondition: &models.Condition{
				Statement: "",
				Allow:     utils.GetPtr(true),
				Paths:     nil,
				Exists:    nil,
				Variables: nil,
			},
		},
		{
			name: "Rule with data and never",
			rule: &common.Rule{
				If:      "condition",
				When:    never,
				Changes: []string{"a", "b", "c"},
				Exists:  []string{"d", "e", "f"},
				Variables: &common.EnvironmentVariablesRef{
					Variables: &common.Variables{
						"string": "string",
						"number": 1,
						"bool":   true,
					},
				},
			},
			expectedCondition: &models.Condition{
				Statement: "condition",
				Allow:     utils.GetPtr(false),
				Paths: &models.Filter{
					DenyList: []string{"a", "b", "c"},
				},
				Exists: &models.Filter{
					DenyList: []string{"d", "e", "f"},
				},
				Variables: map[string]string{
					"string": "string",
					"number": "1",
					"bool":   "true",
				},
			},
		},
		{
			name: "Rule with data and never",
			rule: &common.Rule{
				If:      "condition",
				When:    always,
				Changes: []string{"a", "b", "c"},
				Exists:  []string{"d", "e", "f"},
				Variables: &common.EnvironmentVariablesRef{
					Variables: &common.Variables{
						"string": "string",
						"number": 1,
						"bool":   true,
					},
				},
			},
			expectedCondition: &models.Condition{
				Statement: "condition",
				Allow:     utils.GetPtr(true),
				Paths: &models.Filter{
					AllowList: []string{"a", "b", "c"},
				},
				Exists: &models.Filter{
					AllowList: []string{"d", "e", "f"},
				},
				Variables: map[string]string{
					"string": "string",
					"number": "1",
					"bool":   "true",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseConditionRule(testCase.rule)

			testutils.DeepCompare(t, testCase.expectedCondition, got)

		})
	}
}

// func TestGenerateRuleBranchFilters(t *testing.T) {
// 	testCases := []struct {
// 		name            string
// 		rule            *common.Rule
// 		expectedFilters []*models.Filter
// 	}{
// 		{
// 			name:            "Rule is empty",
// 			rule:            &common.Rule{},
// 			expectedFilters: nil,
// 		},
// 		// {
// 		// 	name: "Rule with branch filters",
// 		// 	rule: &common.Rule{},
// 		// 	expectedFilters: []*models.Filter{
// 		// 		{
// 		// 			AllowList: []string{"a", "b", "c"},
// 		// 		},
// 		// 		{
// 		// 			DenyList: []string{"d", "e", "f"},
// 		// 		},
// 		// 	},
// 		// },
// 	}

// 	for _, testCase := range testCases {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			got := generateRuleBranchFilters(testCase.rule)

// 			assert.Equal(t, testCase.expectedFilters, got, testCase.name)
// 		})
// 	}
// }

func TestGenerateRuleFileFilter(t *testing.T) {
	testCases := []struct {
		name           string
		rule           *common.Rule
		expectedFilter *models.Filter
	}{
		{
			name:           "Rule is empty",
			rule:           &common.Rule{},
			expectedFilter: nil,
		},
		{
			name: "Rule when is empty",
			rule: &common.Rule{
				Changes: []string{"a", "b", "c"},
			},
			expectedFilter: &models.Filter{
				AllowList: []string{"a", "b", "c"},
			},
		},
		{
			name: "Rule when is never",
			rule: &common.Rule{
				When:    never,
				Changes: []string{"a", "b", "c"},
			},
			expectedFilter: &models.Filter{
				DenyList: []string{"a", "b", "c"},
			},
		},
		{
			name: "Rule when is not never",
			rule: &common.Rule{
				When:    always,
				Changes: []string{"a", "b", "c"},
			},
			expectedFilter: &models.Filter{
				AllowList: []string{"a", "b", "c"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := generateRuleFileFilter(testCase.rule)

			testutils.DeepCompare(t, testCase.expectedFilter, got)
		})
	}
}

func TestGenerateRuleExistsFilter(t *testing.T) {
	testCases := []struct {
		name           string
		rule           *common.Rule
		expectedFilter *models.Filter
	}{
		{
			name:           "Rule is empty",
			rule:           &common.Rule{},
			expectedFilter: nil,
		},
		{
			name: "Rule when is empty",
			rule: &common.Rule{
				Exists: []string{"a", "b", "c"},
			},
			expectedFilter: &models.Filter{
				AllowList: []string{"a", "b", "c"},
			},
		},
		{
			name: "Rule is never",
			rule: &common.Rule{
				When:   never,
				Exists: []string{"a", "b", "c"},
			},
			expectedFilter: &models.Filter{
				DenyList: []string{"a", "b", "c"},
			},
		},
		{
			name: "Rule when is not never",
			rule: &common.Rule{
				When:   always,
				Exists: []string{"a", "b", "c"},
			},
			expectedFilter: &models.Filter{
				AllowList: []string{"a", "b", "c"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := generateRuleExistsFilter(testCase.rule)

			testutils.DeepCompare(t, testCase.expectedFilter, got)
		})
	}
}

func TestGenerateRuleVariables(t *testing.T) {
	testCases := []struct {
		name              string
		rule              *common.Rule
		expectedVariables map[string]string
	}{
		{
			name:              "Rule is empty",
			rule:              &common.Rule{},
			expectedVariables: nil,
		},
		{
			name: "Rule with variables",
			rule: &common.Rule{
				Variables: &common.EnvironmentVariablesRef{
					Variables: &common.Variables{
						"string": "string",
						"number": 1,
						"bool":   true,
					},
				},
			},
			expectedVariables: map[string]string{
				"string": "string",
				"number": "1",
				"bool":   "true",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := generateRuleVariables(testCase.rule)

			testutils.DeepCompare(t, testCase.expectedVariables, got)
		})
	}
}
