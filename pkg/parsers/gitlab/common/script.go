package common

import (
	"strings"

	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func ParseScript(script *common.Script) []*models.Step {
	if script == nil {
		return nil
	}

	if len(script.Commands) == 1 { // format: script: command
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

	// format
	// script:
	//   - command1
	//   - command2
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
	scriptLine := script.FileReference.StartRef.Line + commandIndex + 1
	lastCommand := script.Commands[commandIndex]

	if strings.Contains(lastCommand, "\n") {
		splitValue := strings.Split(lastCommand, "\n")
		lastCommand = splitValue[len(splitValue)-1]
	}

	return &models.FileReference{
		StartRef: &models.FileLocation{
			Line:   scriptLine,                               // +1 for the script header
			Column: script.FileReference.StartRef.Column + 2, // +2 for the "- " section
		},
		EndRef: &models.FileLocation{
			Line:   scriptLine,
			Column: script.FileReference.EndRef.Column + len(lastCommand), // start column + the length of the last command
		},
	}

}
