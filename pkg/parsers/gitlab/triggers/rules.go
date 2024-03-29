package triggers

import (
	"fmt"

	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

const (
	never  = "never"
	always = "always"
	manual = "manual"

	eventVariable                    = "$CI_PIPELINE_SOURCE"
	branchVariable                   = "$CI_COMMIT_REF_NAME"
	mergeRequestSourceBranchVariable = "$CI_MERGE_REQUEST_SOURCE_BRANCH_NAME"
)

var (
	eventMapping = map[models.EventType]string{
		models.PullRequestEvent: "merge_request_event",
		models.PushEvent:        "push",
		models.ScheduledEvent:   "schedule",
	}

	branchVariables = []string{
		branchVariable,
		mergeRequestSourceBranchVariable,
	}
)

func ParseRules(rules *common.Rules) (*models.Triggers, []*models.Condition) {
	triggers := []*models.Trigger{}
	conditions := []*models.Condition{}
	for _, rule := range rules.RulesList {
		if ruleTriggers := parseTriggerRules(rule); len(ruleTriggers) > 0 {
			triggers = append(triggers, ruleTriggers...)
		} else {
			conditions = append(conditions, parseConditionRule(rule))
		}
	}

	return &models.Triggers{
		Triggers:      triggers,
		FileReference: rules.FileReference,
	}, conditions
}

func ParseConditionRules(rules *common.Rules) []*models.Condition {
	if rules == nil {
		return nil
	}
	return utils.Map(rules.RulesList, parseConditionRule)
}

func parseConditionRule(rule *common.Rule) *models.Condition {
	return &models.Condition{
		Statement: rule.If,
		Allow:     utils.GetPtr(rule.When != never),
		Paths:     generateRuleFileFilter(rule),
		Exists:    generateRuleExistsFilter(rule),
		Variables: generateRuleVariables(rule),
	}
}

func parseTriggerRules(rule *common.Rule) []*models.Trigger {
	comparisons := getComparisons(rule.If)
	triggers := []*models.Trigger{}
	for _, comparison := range comparisons {
		for event, eventValue := range eventMapping {
			if (comparison.IsPositive() == (rule.When != never)) &&
				comparison.Variable == eventVariable &&
				comparison.Value == fmt.Sprintf(`"%s"`, eventValue) {
				triggers = append(triggers, generateTriggerFromRule(rule, event))
			}
		}
	}
	return triggers
}

func generateTriggerFromRule(rule *common.Rule, event models.EventType) *models.Trigger {
	return &models.Trigger{
		Event:         event,
		FileReference: rule.FileReference,
		Paths:         generateRuleFileFilter(rule),
		Branches:      generateRuleBranchFilter(rule),
	}
}

func generateRuleBranchFilter(rule *common.Rule) *models.Filter {
	denyList := []string{}
	allowList := []string{}
	for _, comparison := range getComparisons(rule.If) {
		if utils.SliceContains(branchVariables, comparison.Variable) {
			if comparison.IsPositive() == (rule.When != never) {
				allowList = append(allowList, comparison.Value)
				continue
			}
			denyList = append(denyList, comparison.Value)
		}
	}

	if len(denyList) > 0 || len(allowList) > 0 {
		return &models.Filter{
			DenyList:  denyList,
			AllowList: allowList,
		}
	}

	return nil
}

func generateRuleFileFilter(rule *common.Rule) *models.Filter {
	if rule.Changes == nil {
		return nil
	}

	if rule.When == never {
		return &models.Filter{
			DenyList: rule.Changes,
		}
	}

	return &models.Filter{
		AllowList: rule.Changes,
	}
}

func generateRuleExistsFilter(rule *common.Rule) *models.Filter {
	if rule.Exists == nil {
		return nil
	}

	if rule.When == never {
		return &models.Filter{
			DenyList: rule.Exists,
		}
	}

	return &models.Filter{
		AllowList: rule.Exists,
	}
}

func generateRuleVariables(rule *common.Rule) map[string]string {
	if rule.Variables == nil {
		return nil
	}

	variables := make(map[string]string)
	for key, value := range *rule.Variables.Variables {
		variables[key] = fmt.Sprint(value)
	}
	return variables
}
