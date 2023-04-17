package azure

import (
	"strings"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	loadersCommonModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/common/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	parserUtils "github.com/argonsecurity/pipeline-parser/pkg/parsers/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"golang.org/x/exp/slices"
)

func parseSteps(steps *azureModels.Steps) []*models.Step {
	if steps == nil {
		return nil
	}

	return utils.Map(*steps, parseStep)
}

func parseStep(step azureModels.Step) *models.Step {
	parsedStep := &models.Step{
		Name:                 &step.DisplayName,
		EnvironmentVariables: parseEnvironmentVariablesRef(step.Env),
		FileReference:        step.FileReference,
	}

	if step.ContinueOnError != nil {
		parsedStep.FailsPipeline = utils.GetPtr(!*step.ContinueOnError)
	}

	if step.Condition != "" {
		parsedStep.Conditions = &[]models.Condition{{Statement: step.Condition}}
	}

	if step.TimeoutInMinutes != 0 {
		parsedStep.Timeout = utils.GetPtr(step.TimeoutInMinutes * 60 * 1000)
	}

	if step.WorkingDirectory != "" {
		parsedStep.WorkingDirectory = &step.WorkingDirectory
	}

	if step.Name != "" {
		parsedStep.ID = &step.Name
	}

	if step.Task != "" {
		actionName, version, versionType := parseTaskHeader(step.Task)
		parsedStep.Task = &models.Task{
			Name:        &actionName,
			Version:     &version,
			VersionType: versionType,
		}

		if step.Inputs != nil {
			parsedStep.Task.Inputs = parserUtils.ParseMapToParameters(loadersCommonModels.Map(*step.Inputs))
		}

		parsedStep.Type = models.TaskStepType
	}

	if shell := parseStepScript(step); shell != nil {
		parsedStep.Shell = shell
		parsedStep.Type = models.ShellStepType
	}

	if step.Enabled != nil {
		parsedStep.Disabled = utils.GetPtr(!slices.Contains(consts.TrueValues, *step.Enabled))
	}

	if step.Template.Template != "" {
		if parsedStep.ID == nil {
			parsedStep.ID = &step.Template.Template
		}
		parsedStep.Imports = &models.Import{
			Source: &models.ImportSource{
				Path: &step.Template.Template,
			},
			Parameters:    step.Template.Parameters,
			FileReference: step.FileReference,
		}
	}

	return parsedStep
}

func parseStepScript(step azureModels.Step) *models.Shell {
	var shellType string
	var script string

	if step.Script != "" {
		shellType = ""
		script = step.Script
	} else if step.Bash != "" {
		shellType = "bash"
		script = step.Bash
	} else if step.Powershell != "" {
		shellType = "powershell"
		script = step.Powershell
	} else if step.Pwsh != "" {
		shellType = "powershell core"
		script = step.Pwsh
	} else {
		return nil
	}

	return &models.Shell{
		Type:   &shellType,
		Script: &script,
	}
}

func parseTaskHeader(header string) (string, string, models.VersionType) {
	result := strings.Split(header, "@")
	if len(result) == 1 {
		return result[0], "", models.None
	}

	return result[0], result[1], models.TagVersion
}
