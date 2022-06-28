package github

import (
	"fmt"
	"regexp"
	"strings"

	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

var (
	sha1Regex   = regexp.MustCompile(`[0-9a-fA-F]{40}`)
	semverRegex = regexp.MustCompile(`v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`)

	githubActionNameRegex = regexp.MustCompile(`(.+?)(?:@(.+)|$)`)

	regexToType = map[*regexp.Regexp]models.VersionType{
		sha1Regex:   models.CommitSHA,
		semverRegex: models.TagVersion,
	}
)

func parseJobSteps(steps *githubModels.Steps) []*models.Step {
	if steps == nil {
		return nil
	}
	parsedSteps := utils.Map(*steps, parseJobStep)
	return parsedSteps
}

func parseJobStep(step githubModels.Step) *models.Step {
	parsedStep := &models.Step{
		Name:                 &step.Name,
		EnvironmentVariables: parseEnvironmentVariablesRef(step.Env),
		FileReference:        step.FileReference,
	}

	if step.ContinueOnError != nil {
		parsedStep.FailsPipeline = utils.GetPtr(!*step.ContinueOnError)
	}

	if step.If != "" {
		parsedStep.Conditions = &[]models.Condition{{Statement: step.If}}
	}

	if step.TimeoutMinutes != 0 {
		parsedStep.Timeout = utils.GetPtr(step.TimeoutMinutes * 60 * 1000)
	}

	if step.WorkingDirectory != "" {
		parsedStep.WorkingDirectory = &step.WorkingDirectory
	}

	if step.Id != "" {
		parsedStep.ID = &step.Id
	}

	if step.Run != nil {
		parsedStep.Shell = &models.Shell{
			Script:        &step.Run.Script,
			FileReference: step.Run.FileReference,
		}

		if step.Shell != "" {
			parsedStep.Shell.Type = &step.Shell
		}
		parsedStep.Type = models.ShellStepType
	} else if step.Uses != "" {
		actionName, version, versionType := parseActionHeader(step.Uses)
		parsedStep.Task = &models.Task{
			Name:        &actionName,
			Version:     &version,
			VersionType: versionType,
			Inputs:      parseActionInput(step.With),
		}
		parsedStep.Type = models.TaskStepType
	}

	return parsedStep
}

func parseActionHeader(header string) (string, string, models.VersionType) {
	result := githubActionNameRegex.FindStringSubmatch(header)
	actionName := result[1]
	versionType := models.None
	version := ""
	if len(result) == 3 {
		version = result[2]
	}

	if version != "" {
		if sha1Regex.MatchString(version) {
			versionType = models.CommitSHA
		} else if semverRegex.MatchString(version) {
			versionType = models.TagVersion
		} else {
			versionType = models.BranchVersion
		}
	}

	return actionName, version, versionType
}

func parseActionInput(with *githubModels.With) *[]models.Parameter {
	if with == nil {
		return nil
	}

	parameters := make([]models.Parameter, 0)
	currentLine := -1
	startColumn := -1

	if with.FileReference != nil {
		currentLine = with.FileReference.StartRef.Line + 1
		startColumn = with.FileReference.StartRef.Column + 2
	}

	for key, value := range with.Inputs {
		name := key
		parameter := models.Parameter{
			Name:          &name,
			Value:         value,
			FileReference: calcParameterFileReference(currentLine, startColumn, key, value),
		}
		currentLine = parameter.FileReference.EndRef.Line + 1
		parameters = append(parameters, parameter)
	}

	return &parameters
}

func calcParameterFileReference(startLine int, startColumn int, key string, val any) *models.FileReference {
	if startLine == -1 || startColumn == -1 {
		return nil
	}

	splitValue := strings.Split(fmt.Sprint(val), "\n")
	valuesLengths := utils.Map(splitValue, func(value string) int { return len(value) })
	longestValue := utils.GetSliceMaxValue(valuesLengths, func(a, b int) bool { return a < b })

	return &models.FileReference{
		StartRef: &models.FileLocation{
			Line:   startLine,
			Column: startColumn, // for the tab after the inputs
		},
		EndRef: &models.FileLocation{
			Line:   startLine + strings.Count(fmt.Sprint(val), "\n"),
			Column: startColumn + len(key) + 2 + longestValue, // for the key: val. len(key) for the key, 2 fot the ": " + longestValue for the value
		},
	}
}
