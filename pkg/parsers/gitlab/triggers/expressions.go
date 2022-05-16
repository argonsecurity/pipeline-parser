package triggers

import (
	"regexp"

	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

// Variable Expressions are GitLab's way to filter according to variable values.
// https://docs.gitlab.com/ee/ci/jobs/job_control.html#cicd-variable-expressions
// TODO: support exists and parentheses

var (
	// Operators
	equals    Operator = "=="
	notEquals Operator = "!="
	match     Operator = "=~"
	notMatch  Operator = "!~"

	// Regexes
	comparisonRegex = regexp.MustCompile(`(\$\w+)\s*(==|!=|=~|!~)\s*(["/].*?["/]|\$\w+)`)
)

type Operator string

type Comparison struct {
	Variable string
	Value    string
	Operator Operator
}

func (c *Comparison) IsPositive() bool {
	return c.Operator == equals || c.Operator == match
}

func getComparisons(expression string) []*Comparison {
	matches := comparisonRegex.FindAllStringSubmatch(expression, -1)
	return utils.Map(matches, func(match []string) *Comparison {
		return &Comparison{
			Variable: match[1],
			Value:    match[3],
			Operator: Operator(match[2]),
		}
	})
}
