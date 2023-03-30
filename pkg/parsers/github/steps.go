package github

import (
	"regexp"

	loadersCommonModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/common/models"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	parserUtils "github.com/argonsecurity/pipeline-parser/pkg/parsers/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

var (
	githubActionNameRegex = regexp.MustCompile(`(.+?)(?:@(.+)|$)`)
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
		isContinueOnError := *step.ContinueOnError == "true"
		parsedStep.FailsPipeline = utils.GetPtr(!isContinueOnError)
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
		}

		if step.With != nil {
			parsedStep.Task.Inputs = parserUtils.ParseMapToParameters(loadersCommonModels.Map(*step.With))
		}

		parsedStep.Type = models.TaskStepType
	}

	return parsedStep
}

func parseActionHeader(header string) (string, string, models.VersionType) {
	result := githubActionNameRegex.FindStringSubmatch(header)
	actionName := result[1]
	version := ""
	if len(result) == 3 {
		version = result[2]
	}

	versionType := parserUtils.DetectVersionType(version)

	return actionName, version, versionType
}
