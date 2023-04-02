package gitlab

import (
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

func TestLoad(t *testing.T) {
	testCases := []TestCase{
		{
			Name:     "Build Job",
			Filename: "../../../test/fixtures/gitlab/build-job.yaml",
			ExpectedGitlabCIConfig: &models.GitlabCIConfiguration{
				Workflow: &models.Workflow{
					Rules: &common.Rules{
						RulesList: []*common.Rule{
							{
								If:            `$CI_PIPELINE_SOURCE == "merge_request_event"`,
								FileReference: testutils.CreateFileReference(23, 7, 23, 55),
							},
							{
								If:            "$CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH",
								When:          "never",
								FileReference: testutils.CreateFileReference(24, 7, 25, 18),
							},
						},
						FileReference: testutils.CreateFileReference(22, 3, 25, 18),
					},
				},
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
				Default: &models.Default{
					Artifacts: &models.Artifacts{
						Reports: &models.Reports{
							SecretDetection:    "secrets.json",
							Sast:               "sast.json",
							Terraform:          "terraform.json",
							DependencyScanning: "dependency.json",
							LicenseScanning:    "license.json",
						},
					},
				},
			},
		},
		{
			Name:     "Gradle",
			Filename: "../../../test/fixtures/gitlab/gradle.yaml",
			ExpectedGitlabCIConfig: &models.GitlabCIConfiguration{
				Image: &common.Image{
					Name:          "gradle:alpine",
					FileReference: testutils.CreateFileReference(10, 8, 10, 28),
				},
				Variables: &common.EnvironmentVariablesRef{
					Variables: &common.Variables{
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
			Filename: "../../../test/fixtures/gitlab/terraform.yaml",
			ExpectedGitlabCIConfig: &models.GitlabCIConfiguration{
				Stages: []string{
					"validate",
					"test",
					"build",
					"deploy",
				},
				Jobs: map[string]*models.Job{
					"validate": {
						Extends:       ".terraform:validate",
						FileReference: testutils.CreateFileReference(16, 1, 18, 10),
						Needs:         &job.Needs{},
					},
					"fmt": {
						Extends:       ".terraform:fmt",
						FileReference: testutils.CreateFileReference(12, 1, 14, 10),
						Needs:         &job.Needs{},
					},
					"build": {
						Extends:       ".terraform:build",
						FileReference: testutils.CreateFileReference(20, 1, 21, 28),
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
		{
			Name:     "Baserow",
			Filename: "../../../test/fixtures/gitlab/baserow.yaml",
			ExpectedGitlabCIConfig: &models.GitlabCIConfiguration{
				Stages: []string{
					"build",
					"test",
					"build-final",
					"publish",
				},
				Jobs: map[string]*models.Job{
					"build-ci-util-image": {
						Image: &common.Image{
							Name:          "docker:20.10.12",
							FileReference: testutils.CreateFileReference(202, 3, 202, 25),
						},
						Stage:    "build",
						Services: []any{"docker:20.10.12-dind"},
						Variables: &common.EnvironmentVariablesRef{
							Variables: &common.Variables{
								"DOCKER_BUILDKIT": "1",
								"DOCKER_HOST":     "tcp://docker:2376",
							},
							FileReference: testutils.CreateFileReference(205, 1, 208, 37),
						},
						BeforeScript: &common.Script{
							Commands: []string{
								"echo \"$CI_REGISTRY_PASSWORD\" | \\\n  docker login -u \"$CI_REGISTRY_USER\" \"$CI_REGISTRY\" --password-stdin\n",
							},
							FileReference: testutils.CreateFileReference(209, 3, 212, 7),
						},
						Script: &common.Script{
							Commands: []string{
								"cd .gitlab/ci_util_image",
								"docker build -t $CI_UTIL_IMAGE .",
								"docker push $CI_UTIL_IMAGE",
							},
							FileReference: testutils.CreateFileReference(213, 3, 216, 33),
						},
						When: "manual",
						Only: &job.Controls{
							Changes:       []string{".gitlab/ci_util_image/*"},
							FileReference: testutils.CreateFileReference(218, 3, 220, 32),
						},
						Except: &job.Controls{
							Refs: []string{
								"pipelines",
								"tags",
							},
							FileReference: testutils.CreateFileReference(221, 3, 224, 13),
						},
						Parallel: &job.Parallel{
							Matrix: &job.Matrix{
								job.MatrixItem{
									"key1": []string{"value1", "value2"},
								},
								job.MatrixItem{
									"key2": []string{"value"},
								},
							},
						},
						FileReference: testutils.CreateFileReference(201, 1, 228, 20),
					},
				},
				Variables: &common.EnvironmentVariablesRef{
					Variables: &common.Variables{
						"TRIGGER_FULL_IMAGE_REBUILD": "no",
						"ENABLE_JOB_SKIPPING":        "false",
						"ENABLE_COVERAGE":            "true",
						"ENABLE_RELEASES":            "false",
						"TESTED_IMAGE_PREFIX":        "ci-tested-",
						"BACKEND_IMAGE_NAME":         "backend",
					},
					FileReference: testutils.CreateFileReference(184, 1, 198, 32),
				},
			},
		},
		{
			Name:     "Include Local",
			Filename: "../../../test/fixtures/gitlab/include-local.yaml",
			ExpectedGitlabCIConfig: &models.GitlabCIConfiguration{
				Include: &models.Include{
					{
						Local:         "/../../test/fixtures/gitlab/gradle.yaml",
						FileReference: testutils.CreateFileReference(1, 10, 1, 49),
					},
				},
			},
		},
		{
			Name:     "Include Remote",
			Filename: "../../../test/fixtures/gitlab/include-remote.yaml",
			ExpectedGitlabCIConfig: &models.GitlabCIConfiguration{
				Include: &models.Include{
					{
						Remote:        "https://gitlab.com/gitlab-org/gitlab/-/raw/master/imported.yaml",
						FileReference: testutils.CreateFileReference(1, 10, 1, 73),
					},
				},
			},
		},
		{
			Name:     "Include Multiple",
			Filename: "../../../test/fixtures/gitlab/include-multiple.yaml",
			ExpectedGitlabCIConfig: &models.GitlabCIConfiguration{
				Include: &models.Include{
					{
						File:          "/imported.yaml",
						Ref:           "master",
						Project:       "gitlab-org/gitlab",
						FileReference: testutils.CreateFileReference(2, 5, 4, 16),
					},
					{
						Remote:        "https://gitlab.com/gitlab-org/gitlab/-/raw/master/imported.yaml",
						FileReference: testutils.CreateFileReference(5, 5, 5, 68),
					},
					{
						Local:         "/../../test/fixtures/gitlab/gradle.yaml",
						FileReference: testutils.CreateFileReference(6, 5, 6, 44),
					},
					{
						Template:      "Android.gitlab-ci.yml",
						FileReference: testutils.CreateFileReference(7, 5, 7, 36),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			gitlabLoader := &GitLabLoader{}
			gitlabCIConfig, err := gitlabLoader.Load(testutils.ReadFile(tc.Filename))

			if err != tc.ExpectedError {
				t.Errorf("Expected error: %v, got: %v", tc.ExpectedError, err)
			}

			changelog, _ := diff.Diff(gitlabCIConfig, tc.ExpectedGitlabCIConfig)

			if len(changelog) > 0 {
				t.Errorf("Loader result is not as expected:")
				for _, change := range changelog {
					t.Errorf("field: %s, got: %v, expected: %v", strings.Join(change.Path, "."), change.From, change.To)
				}
			}
		})
	}
}
