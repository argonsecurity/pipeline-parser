package gitlab

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/job"
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
				BeforeScript: &common.Script{
					Commands:      []string{`echo "before_script"`},
					FileReference: testutils.CreateFileReference(18, 1, 19, 25),
				},
				Jobs: map[string]*models.Job{
					"python-build": {
						Rules: &common.Rules{
							RulesList: []*common.Rule{
								{
									If:            "$CI_MERGE_REQUEST_SOURCE_BRANCH_NAME =~ /^feature/",
									FileReference: testutils.CreateFileReference(13, 7, 13, 61),
								},
							},
							FileReference: testutils.CreateFileReference(12, 3, 13, 61),
						},
						FileReference: testutils.CreateFileReference(4, 1, 16, 29),
						Stage:         "build",
						Script: &common.Script{
							Commands: []string{
								"cd requests",
								"python3 setup.py sdist",
							},
							FileReference: testutils.CreateFileReference(14, 3, 16, 29),
						},
						Only: &job.Controls{
							Refs: []string{
								"merge_requests",
								"/^feature-.*/",
								"main",
								"api",
							},
							FileReference: testutils.CreateFileReference(5, 3, 10, 12),
						},
					},
				},
			},
		},
		{
			Name:     "Gradle",
			Filename: "gradle.yaml",
			ExpectedGitlabCIConfig: &models.GitlabCIConfiguration{
				Image: &common.Image{
					Name:          "gradle:alpine",
					FileReference: testutils.CreateFileReference(10, 1, 10, 21),
				},
				Variables: &common.EnvironmentVariablesRef{
					Variables: map[string]any{
						"GRADLE_OPTS": "-Dorg.gradle.daemon=false",
					},
					FileReference: testutils.CreateFileReference(16, 1, 17, 43),
				},
				BeforeScript: &common.Script{
					Commands: []string{
						`GRADLE_USER_HOME="$(pwd)/.gradle"`,
						`export GRADLE_USER_HOME`,
					},
					FileReference: testutils.CreateFileReference(19, 1, 21, 28),
				},
				Jobs: map[string]*models.Job{
					"build": {
						Stage: "build",
						Script: &common.Script{
							Commands: []string{
								`gradle --build-cache assemble`,
							},
							FileReference: testutils.CreateFileReference(25, 3, 25, 40),
						},
						Cache: &common.Cache{
							Key:    "$CI_COMMIT_REF_NAME",
							Policy: "push",
							Paths: []string{
								"build",
								".gradle",
							},
						},
						FileReference: testutils.CreateFileReference(23, 1, 31, 16),
					},
					"test": {
						Stage: "test",
						Script: &common.Script{
							Commands:      []string{`gradle check`},
							FileReference: testutils.CreateFileReference(35, 3, 35, 23),
						},
						FileReference: testutils.CreateFileReference(33, 1, 35, 23),
					},
				},
			},
		},
		{
			Name:     "Terraform",
			Filename: "terraform.yaml",
			ExpectedGitlabCIConfig: &models.GitlabCIConfiguration{
				Include: &models.Include{
					{
						Template:      "Terraform/Base.latest.gitlab-ci.yml",
						FileReference: testutils.CreateFileReference(7, 5, 7, 50),
					},
					{
						Template:      "Jobs/SAST-IaC.latest.gitlab-ci.yml",
						FileReference: testutils.CreateFileReference(8, 5, 8, 49),
					},
				},
				Stages: []string{
					"validate",
					"test",
					"build",
					"deploy",
				},
				Jobs: map[string]*models.Job{
					"validate": {
						Extends:       ".terraform:validate",
						FileReference: testutils.CreateFileReference(20, 1, 22, 10),
						Needs:         &job.Needs{},
					},
					"fmt": {
						Extends:       ".terraform:fmt",
						FileReference: testutils.CreateFileReference(16, 1, 18, 10),
						Needs:         &job.Needs{},
					},
					"build": {
						Extends:       ".terraform:build",
						FileReference: testutils.CreateFileReference(24, 1, 25, 28),
					},
					"deploy": {
						Extends: ".terraform:deploy",
						Dependencies: []string{
							"build",
						},
						Environment: map[string]string{
							"name": "$TF_STATE_NAME",
						},
						FileReference: testutils.CreateFileReference(27, 1, 32, 10),
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
