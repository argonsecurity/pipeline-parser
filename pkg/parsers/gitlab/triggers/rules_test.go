package triggers

import (
	"fmt"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestParseRules(t *testing.T) {
	testCases := []struct {
		name               string
		rules              *common.Rules
		expectedTriggers   *models.Triggers
		expectedConditions []*models.Condition
	}{
		{
			name:  "Rules are empty",
			rules: &common.Rules{},
			expectedTriggers: &models.Triggers{
				Triggers: []*models.Trigger{},
			},
			expectedConditions: []*models.Condition{},
		},
		{
			name:  "Rules are empty",
			rules: &common.Rules{},
			expectedTriggers: &models.Triggers{
				Triggers: []*models.Trigger{},
			},
			expectedConditions: []*models.Condition{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			gotTriggers, gotConditions := ParseRules(testCase.rules)

			testutils.DeepCompare(t, testCase.expectedTriggers, gotTriggers)
			testutils.DeepCompare(t, testCase.expectedConditions, gotConditions)
		})
	}
}

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

func TestParseTriggerRules(t *testing.T) {
	testCases := []struct {
		name             string
		rule             *common.Rule
		expectedTriggers []*models.Trigger
	}{
		{
			name:             "Rule is empty",
			rule:             &common.Rule{},
			expectedTriggers: []*models.Trigger{},
		},
		{
			name: "Rule with data, condition == and always",
			rule: &common.Rule{
				If:            `$CI_PIPELINE_SOURCE == "merge_request_event"`,
				When:          always,
				Changes:       []string{"a", "b", "c"},
				Exists:        []string{"d", "e", "f"},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTriggers: []*models.Trigger{
				{
					Event: models.PullRequestEvent,
					Paths: &models.Filter{
						AllowList: []string{"a", "b", "c"},
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
			},
		},
		{
			name: "Rule with data, condition == and never",
			rule: &common.Rule{
				If:            `$CI_PIPELINE_SOURCE == "merge_request_event"`,
				When:          never,
				Changes:       []string{"a", "b", "c"},
				Exists:        []string{"d", "e", "f"},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTriggers: []*models.Trigger{},
		},
		{
			name: "Rule with data, condition != and never",
			rule: &common.Rule{
				If:            `$CI_PIPELINE_SOURCE != "merge_request_event"`,
				When:          never,
				Changes:       []string{"a", "b", "c"},
				Exists:        []string{"d", "e", "f"},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTriggers: []*models.Trigger{
				{
					Event: models.PullRequestEvent,
					Paths: &models.Filter{
						DenyList: []string{"a", "b", "c"},
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
			},
		},
		{
			name: "Rule with data, condition != and always",
			rule: &common.Rule{
				If:            `$CI_PIPELINE_SOURCE != "merge_request_event"`,
				When:          always,
				Changes:       []string{"a", "b", "c"},
				Exists:        []string{"d", "e", "f"},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTriggers: []*models.Trigger{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseTriggerRules(testCase.rule)

			testutils.DeepCompare(t, testCase.expectedTriggers, got)
		})
	}
}

func TestGenerateTriggerFromRule(t *testing.T) {
	testCases := []struct {
		name            string
		rule            *common.Rule
		event           models.EventType
		expectedTrigger *models.Trigger
	}{
		{
			name:            "Rule is empty",
			rule:            &common.Rule{},
			expectedTrigger: &models.Trigger{},
		},
		{
			name:  "Rule with data",
			event: models.PullRequestEvent,
			rule: &common.Rule{
				If:            fmt.Sprintf(`%s == "a"`, branchVariable),
				When:          always,
				Changes:       []string{"a", "b", "c"},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event: models.PullRequestEvent,
				Branches: &models.Filter{
					AllowList: []string{`"a"`},
					DenyList:  []string{},
				},
				Paths: &models.Filter{
					AllowList: []string{"a", "b", "c"},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := generateTriggerFromRule(testCase.rule, testCase.event)

			testutils.DeepCompare(t, testCase.expectedTrigger, got)

		})
	}
}

func TestGenerateRuleBranchFilters(t *testing.T) {
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
			name: "Rule with if with no branch variable",
			rule: &common.Rule{
				If: `$var == "value"`,
			},
			expectedFilter: nil,
		},
		{
			name: "Rule with positive if, branch variable and when is always",
			rule: &common.Rule{
				When: always,
				If:   fmt.Sprintf(`%s == "a"`, branchVariable),
			},
			expectedFilter: &models.Filter{
				AllowList: []string{`"a"`},
				DenyList:  []string{},
			},
		},
		{
			name: "Rule with negative if, branch variable and when is always",
			rule: &common.Rule{
				When: always,
				If:   fmt.Sprintf(`%s != "a"`, branchVariable),
			},
			expectedFilter: &models.Filter{
				AllowList: []string{},
				DenyList:  []string{`"a"`},
			},
		},
		{
			name: "Rule with positive if, branch variable and when is never",
			rule: &common.Rule{
				When: never,
				If:   fmt.Sprintf(`%s == "a"`, branchVariable),
			},
			expectedFilter: &models.Filter{
				AllowList: []string{},
				DenyList:  []string{`"a"`},
			},
		},
		{
			name: "Rule with negative if, branch variable and when is never",
			rule: &common.Rule{
				When: never,
				If:   fmt.Sprintf(`%s != "a"`, branchVariable),
			},
			expectedFilter: &models.Filter{
				AllowList: []string{`"a"`},
				DenyList:  []string{},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := generateRuleBranchFilter(testCase.rule)

			testutils.DeepCompare(t, testCase.expectedFilter, got)
		})
	}
}

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
