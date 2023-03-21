package general

import (
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers/general/config"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func enhanceStep(step *models.Step, config *config.EnhancementConfiguration) *models.Step {
	if step.Type == models.ShellStepType {
		step = enhanceShellStep(step, config)
	}

	if step.Type == models.TaskStepType {
		step = enhanceTaskStep(step, config)
	}

	if utils.AnyMatch(config.Build.Names, step.Name) {
		step.Metadata.Build = true
	}

	if utils.AnyMatch(config.Test.Names, step.Name) {
		step.Metadata.Test = true
	}

	if utils.AnyMatch(config.Deploy.Names, step.Name) {
		step.Metadata.Deploy = true
	}

	return step
}

func enhanceShellStep(step *models.Step, config *config.EnhancementConfiguration) *models.Step {
	if utils.AnyMatch(config.Build.ShellRegexes, step.Shell.Script) {
		step.Metadata.Build = true
	}
	if utils.AnyMatch(config.Test.ShellRegexes, step.Shell.Script) {
		step.Metadata.Test = true
	}

	if utils.AnyMatch(config.Deploy.ShellRegexes, step.Shell.Script) {
		step.Metadata.Deploy = true
	}

	return step
}

func enhanceTaskStep(step *models.Step, config *config.EnhancementConfiguration) *models.Step {
	if utils.SliceContains(config.Build.Tasks, *step.Task.Name) {
		step.Metadata.Build = true
	}

	if utils.SliceContains(config.Test.Tasks, *step.Task.Name) {
		step.Metadata.Test = true
	}

	if utils.SliceContains(config.Deploy.Tasks, *step.Task.Name) {
		step.Metadata.Deploy = true
	}

	return step
}
