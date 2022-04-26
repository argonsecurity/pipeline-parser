package enhancers

import (
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers/config"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func enhanceStep(step models.Step, config config.EnhancementConfiguration) models.Step {
	if step.Type == models.ShellStepType {
		if utils.AnyMatch(config.Build.ShellRegexes, *step.Shell.Script) {
			step.Metadata.Build = true
		}
		if utils.AnyMatch(config.Test.ShellRegexes, *step.Shell.Script) {
			step.Metadata.Test = true
		}

		if utils.AnyMatch(config.Deploy.ShellRegexes, *step.Shell.Script) {
			step.Metadata.Deploy = true
		}
	}

	if step.Type == models.TaskStepType {
		if utils.SliceContains(config.Build.Tasks, *step.Task.Name) {
			step.Metadata.Build = true
		}

		if utils.AnyMatch(config.Test.Names, *step.Name) {
			step.Metadata.Test = true
		}

		if utils.AnyMatch(config.Deploy.Names, *step.Name) {
			step.Metadata.Deploy = true
		}
	}

	if utils.AnyMatch(config.Build.Names, *step.Name) {
		step.Metadata.Build = true
	}

	if utils.AnyMatch(config.Test.Names, *step.Name) {
		step.Metadata.Test = true
	}

	if utils.AnyMatch(config.Deploy.Names, *step.Name) {
		step.Metadata.Deploy = true
	}

	return step
}
