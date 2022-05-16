package blackbox

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestGitLab(t *testing.T) {
	testCases := []TestCase{
		{
			Filename: "gradle.yaml",
			Expected: &models.Pipeline{
				Jobs: []*models.Job{
					{},
				},
				Defaults: &models.Defaults{
					Runner: &models.Runner{
						DockerMetadata: &models.DockerMetadata{
							Image: utils.GetPtr("gradle"),
							Label: utils.GetPtr("alpine"),
						},
						FileReference: testutils.CreateFileReference(10, 1, 10, 8),
					},
					EnvironmentVariables: &models.EnvironmentVariablesRef{
						EnvironmentVariables: models.EnvironmentVariables{
							"GRADLE_OPTS": "-Dorg.gradle.daemon=false",
						},
						FileReference: testutils.CreateFileReference(16, 0, 17, 16),
					},
				},
			},
		},
	}

	executeTestCases(t, testCases, "gitlab", consts.GitLabPlatform)
}
