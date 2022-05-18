package common

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func ParseScript(script *common.Script) []*models.Step {
	if script == nil {
		return nil
	}

	if len(script.Commands) == 1 {
		return []*models.Step{
			{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: &script.Commands[0],
				},
				FileReference: script.FileReference,
			},
		}
	}

	return utils.MapWithIndex(script.Commands, func(command string, index int) *models.Step {
		return &models.Step{
			Type: models.ShellStepType,
			Shell: &models.Shell{
				Script: &command,
			},
			FileReference: parseCommandFileReference(script, index),
		}
	})
}

func parseCommandFileReference(script *common.Script, commandIndex int) *models.FileReference {
	return &models.FileReference{
		StartRef: &models.FileLocation{
			Line:   script.FileReference.StartRef.Line + commandIndex + 1, // +1 for the script header
			Column: script.FileReference.StartRef.Column + 2,              // +2 for the "- " section
		},
		EndRef: &models.FileLocation{
			Line:   script.FileReference.StartRef.Line + commandIndex + 1,
			Column: script.FileReference.EndRef.Column + len(script.Commands[commandIndex]),
		},
	}

}
