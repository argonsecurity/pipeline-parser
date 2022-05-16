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

	eventVariable  = "$CI_PIPELINE_SOURCE"
	branchVariable = "$CI_COMMIT_REF_NAME"
)

var (
	eventMapping = map[models.EventType]string{
		models.PullRequestEvent: "merge_request_event",
		models.PushEvent:        "push",
		models.ScheduledEvent:   "schedule",
	}
)

func ParseRules(rules *common.Rules) (*models.Triggers, []*models.Condition) {
	triggers := []*models.Trigger{}
	conditions := []*models.Condition{}
	for _, rule := range rules.RulesList {
		if ruleTriggers := parseRuleTriggers(rule); ruleTriggers != nil {
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

func ParseRulesConditions(rules *common.Rules) []*models.Condition {
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

func parseRuleTriggers(rule *common.Rule) []*models.Trigger {
	comparisons := getComparisons(rule.If)
	triggers := []*models.Trigger{}
	for event, eventValue := range eventMapping {
		for _, comparison := range comparisons {
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
		Branches:      generateRuleBranchFilters(rule),
	}
}

func generateRuleBranchFilters(rule *common.Rule) *models.Filter {
	denyList := []string{}
	allowList := []string{}
	for _, comparison := range getComparisons(rule.If) {
		if comparison.Variable == branchVariable {
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
	for key, value := range rule.Variables.Variables {
		variables[key] = fmt.Sprint(value)
	}
	return variables
}
