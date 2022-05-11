package gitlab

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/r3labs/diff/v3"
)

type TestCase struct {
	Name                   string
	Filename               string
	ExpectedError          error
	ExpectedGitlabCIConfig *models.GitlabCIConfiguration
}

func readFile(filename string) []byte {
	data, err := ioutil.ReadFile(filepath.Join("../../../test/fixtures/gitlab", filename))
	if err != nil {
		panic(err)
	}

	return data
}

func TestGitlabLoader(t *testing.T) {
	testCases := []TestCase{
		{
			Name:     "Build Job",
			Filename: "build-job.yaml",
			ExpectedGitlabCIConfig: &models.GitlabCIConfiguration{
				Stages: []string{"build"},
				BeforeScript: []*common.Script{
					{
						Commands:      []string{`echo "before_script"`},
						FileReference: testutils.CreateFileReference(11, 5, 11, 5),
					},
				},
				Jobs: map[string]*models.Job{
					"python-build": {
						FileReference: testutils.CreateFileReference(4, 0, 8, 7),
						Stage:         "build",
						Script: &common.Script{
							Commands: []string{"cd requests",
								"python3 setup.py sdist",
							},
							FileReference: testutils.CreateFileReference(6, 3, 8, 7),
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			gitlabLoader := &GitLabLoader{}
			gitlabCIConfig, err := gitlabLoader.Load(readFile(tc.Filename))

			if err != tc.ExpectedError {
				t.Errorf("Expected error: %v, got: %v", tc.ExpectedError, err)
			}

			changelog, _ := diff.Diff(gitlabCIConfig, tc.ExpectedGitlabCIConfig)

			if len(changelog) > 0 {
				t.Errorf("Loader result is not as expected:")
				for _, change := range changelog {
					t.Errorf("field: %s, from: %v, to: %v", strings.Join(change.Path, "."), change.From, change.To)
				}
			}
		})
	}
}
