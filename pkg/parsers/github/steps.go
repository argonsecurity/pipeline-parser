package github

import (
	"fmt"
	"regexp"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/parsers/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

var (
	sha1Regex        = regexp.MustCompile(`[0-9a-fA-F]{40}`)
	semverRegex      = regexp.MustCompile(`^v?((0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))|\d+)?`)
	semverMajorRegex = regexp.MustCompile(`v?(0|[1-9]\d*)`)

	githubActionNameRegex = regexp.MustCompile(`(.+?)(?:@(.+)|$)`)

	regexToType = map[*regexp.Regexp]models.VersionType{
		sha1Regex:        models.CommitSHA,
		semverRegex:      models.TagVersion,
		semverMajorRegex: models.TagVersion,
	}
)

func parseJobSteps(steps *[]githubModels.Step) *[]models.Step {
	if steps == nil {
		return nil
	}

	parsedSteps := utils.Map(*steps, parseJobStep)
	return utils.GetPtr(parsedSteps)
}

func parseJobStep(step githubModels.Step) models.Step {
	var failsPipeline *bool
	if step.ContinueOnError != nil {
		failsPipeline = utils.GetPtr(!*step.ContinueOnError)
	}

	parsedStep := models.Step{
		Name:                 &step.Name,
		ID:                   &step.Id,
		FailsPipeline:        failsPipeline,
		EnvironmentVariables: parseEnvironmentVariables(step.Env),
		WorkingDirectory:     &step.WorkingDirectory,
		Conditions:           &[]models.Condition{models.Condition(step.If)},
		Timeout:              utils.GetPtr(step.TimeoutMinutes * 60 * 1000),
	}

	if step.Run != "" {
		parsedStep.Shell = &step.Run
		parsedStep.Type = models.ShellType
	} else if step.Uses != "" {
		actionName, version, versionType := parseActionHeader(step.Uses)
		parsedStep.Task = &models.Task{
			Name:        &actionName,
			ID:          &step.Id,
			Version:     &version,
			VersionType: versionType,
			Inputs:      parseActionInput(step.With),
		}
		parsedStep.Type = models.TaskType
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
		} else if semverMajorRegex.MatchString(version) {
			versionType = models.TagVersion
		} else {
			versionType = models.BranchVersion
		}
	}

	return actionName, version, versionType
}

func parseActionInput(with map[string]any) *[]models.Parameter {
	if with == nil {
		return nil
	}

	parsedInputs := utils.MapToSlice(with, parseActionInputItem)
	return &parsedInputs
}

func parseActionInputItem(k string, val any) models.Parameter {
	return models.Parameter{
		Name:  &k,
		Value: utils.GetPtr(fmt.Sprint(val)),
	}
}

// 	Shell            any
// 	With             any
