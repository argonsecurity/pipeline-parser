package github

import (
	"github.com/argonsecurity/pipeline-parser/pkg/models"

	githubModels "github.com/argonsecurity/pipeline-parser/pkg/parsers/github/models"
)

func parseRunsOnToRunner(runsOn *githubModels.RunsOn) *models.Runner {
	if runsOn == nil {
		return nil
	}

	runner := &models.Runner{
		OS:     runsOn.OS,
		Arch:   runsOn.Arch,
		Labels: &runsOn.Tags,
	}
	return runner
}
