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
	command := script.Commands[commandIndex]
	
	firstLine := script.FileReference.StartRef.Line + countLines(script, commandIndex)
	endLine := firstLine

	// handle multiline command
	if strings.Contains(command, "\n") {
		splitValue := strings.Split(command, "\n")
		command = splitValue[len(splitValue)-1]
		endLine = firstLine + len(splitValue) - 1
	}

	return &models.FileReference{
		StartRef: &models.FileLocation{
			Line:   firstLine,                               // +1 for the script header
			Column: script.FileReference.StartRef.Column + 2, // +2 for the "- " section
		},
		EndRef: &models.FileLocation{
			Line:   endLine,
			Column: script.FileReference.EndRef.Column + len(command), // start column + the length of the last command
		},
	}
}

func countLines(script *common.Script, commandIndex int) int {
	lines := 0
	for i := 0; i < commandIndex; i++ {
		splitCommand := strings.Split(script.Commands[i], "\n")
		lines += len(splitCommand)
	}

	return lines
}
