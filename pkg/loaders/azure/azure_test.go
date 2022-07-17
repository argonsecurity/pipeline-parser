package azure

import (
	"strings"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/r3labs/diff/v3"
)

func TestLoad(t *testing.T) {
	testCases := []struct {
		name             string
		filename         string
		expectedPipeline *models.Pipeline
		expectedError    error
	}{
		{
			name:     "all-triggers",
			filename: "../../../test/fixtures/azure/all-triggers.yaml",
			expectedPipeline: &models.Pipeline{
				Name: "all-triggers",
				Trigger: &models.TriggerRef{
					Trigger: &models.Trigger{
						Batch: true,
						Branches: &models.Filter{
							Include: []string{"master", "main"},
							Exclude: []string{"test/*"},
						},
						Paths: &models.Filter{
							Include: []string{"path/to/file", "another/path/to/file"},
							Exclude: []string{"all/*"},
						},
						Tags: &models.Filter{
							Include: []string{"v1.0.*"},
							Exclude: []string{"v2.0.*"},
						},
					},
					FileReference: testutils.CreateFileReference(2, 1, 20, 15),
				},
				PR: &models.PRRef{
					PR: &models.PR{
						AutoCancel: true,
						Branches: &models.Filter{
							Include: []string{"features/*"},
							Exclude: []string{"features/experimental/*"},
						},
						Paths: &models.Filter{
							Include: []string{"path/to/file"},
							Exclude: []string{"README.md"},
						},
						Drafts: true,
					},
					FileReference: testutils.CreateFileReference(21, 1, 33, 15),
				},
				Schedules: &models.Schedules{
					Crons: &[]models.Cron{
						{
							Cron:        "0 0 * * *",
							DisplayName: "Daily midnight build",
							Branches: &models.Filter{
								Include: []string{"main", "releases/*"},
								Exclude: []string{"releases/ancient/*"},
							},
							FileReference: testutils.CreateFileReference(35, 3, 42, 25),
						},
						{
							Cron:        "0 12 * * 0",
							DisplayName: "Weekly Sunday build",
							Branches: &models.Filter{
								Include: []string{"releases/*"},
							},
							Always:        true,
							FileReference: testutils.CreateFileReference(43, 3, 48, 15),
						},
					},
					FileReference: testutils.CreateFileReference(34, -1, 48, 15),
				},
			},
		},
		{
			name:     "no-trigger",
			filename: "../../../test/fixtures/azure/no-trigger.yaml",
			expectedPipeline: &models.Pipeline{
				Name:    "no-trigger",
				Trigger: &models.TriggerRef{FileReference: testutils.CreateFileReference(2, 10, 2, 14)},
				PR:      &models.PRRef{FileReference: testutils.CreateFileReference(3, 5, 3, 9)},
			},
		},
		{
			name:     "branch-list-trigger",
			filename: "../../../test/fixtures/azure/branch-list-trigger.yaml",
			expectedPipeline: &models.Pipeline{
				Name: "branch-list-trigger",
				Trigger: &models.TriggerRef{
					Trigger: &models.Trigger{
						Branches: &models.Filter{
							Include: []string{"main", "development"},
						},
					},
					FileReference: testutils.CreateFileReference(2, 1, 4, 16),
				},
				PR: &models.PRRef{
					PR: &models.PR{
						Branches: &models.Filter{
							Include: []string{"main", "develop"},
						},
					},
					FileReference: testutils.CreateFileReference(5, -1, 7, 10),
				},
			},
		},
		{
			name:     "variables",
			filename: "../../../test/fixtures/azure/variables.yaml",
			expectedPipeline: &models.Pipeline{
				Name: "variables",
				Variables: &models.Variables{
					{
						Name:          "var1",
						Value:         "value1",
						FileReference: testutils.CreateFileReference(3, 3, 4, 16),
					},
					{
						Name:          "var2",
						Value:         "value2",
						Readonly:      true,
						FileReference: testutils.CreateFileReference(5, 3, 7, 17),
					},
					{
						Group:         "my-group",
						FileReference: testutils.CreateFileReference(8, 3, 8, 18),
					},
					{
						Template: models.Template{
							Template: "variables/var.yml",
							Parameters: map[string]any{
								"param": "value",
							},
						},
						FileReference: testutils.CreateFileReference(9, 3, 11, 17),
					},
				},
				Stages: &models.Stages{
					Stages: &[]models.Stage{
						{
							Stage: "Build",
							Variables: &models.Variables{
								{
									Name:          "STAGE_VAR",
									Value:         "that happened",
									FileReference: testutils.CreateFileReference(16, 5, 16, 18),
								},
							},
							Jobs: &models.Jobs{
								CIJobs: &[]models.CIJob{
									{
										Job: "FirstJob",
										BaseJob: models.BaseJob{
											Variables: &models.Variables{
												{
													Name:          "JOB_VAR",
													Value:         "a job var",
													FileReference: testutils.CreateFileReference(21, 7, 21, 16),
												},
											},
											Steps: &models.Steps{
												{
													Script:        "echo $(MY_VAR) $(STAGE_VAR) $(JOB_VAR)",
													FileReference: testutils.CreateFileReference(23, 7, 23, 53),
												},
											},
										},
										FileReference: testutils.CreateFileReference(19, 5, 23, 53),
									},
								},
								DeploymentJobs: &[]models.DeploymentJob{},
								TemplateJobs:   &[]models.TemplateJob{},
								FileReference:  testutils.CreateFileReference(18, 1, 23, 53),
							},
							FileReference: testutils.CreateFileReference(14, 3, 23, 53),
						},
					},
					TemplateStages: &[]models.TemplateStage{},
					FileReference:  testutils.CreateFileReference(13, -1, 23, 53),
				},
			},
		},
		{
			name:     "pool",
			filename: "../../../test/fixtures/azure/pool.yaml",
			expectedPipeline: &models.Pipeline{
				Name: "pool",
				Pool: &models.Pool{
					Name:          "MyPool",
					Demands:       []string{"demand1", "demand2"},
					VmImage:       "ubuntu-latest",
					FileReference: testutils.CreateFileReference(2, 3, 5, 25),
				},
				Jobs: &models.Jobs{
					CIJobs: &[]models.CIJob{
						{
							Job: "jobId",
							BaseJob: models.BaseJob{
								Pool: &models.Pool{
									Name:          "jobPool",
									Demands:       []string{"demand"},
									VmImage:       "ubuntu-latest",
									FileReference: testutils.CreateFileReference(10, 5, 13, 27),
								},
							},
							FileReference: testutils.CreateFileReference(9, 3, 13, 18),
						},
					},
					DeploymentJobs: &[]models.DeploymentJob{},
					TemplateJobs:   &[]models.TemplateJob{},
					FileReference:  testutils.CreateFileReference(8, -1, 13, 18),
				},
			},
		},
		{
			name:     "parameters",
			filename: "../../../test/fixtures/azure/parameters.yaml",
			expectedPipeline: &models.Pipeline{
				Name: "parameters",
				Parameters: &models.Parameters{
					{
						Name:          "myString",
						Type:          "string",
						Default:       "a string",
						FileReference: testutils.CreateFileReference(3, 3, 5, 20),
					},
					{
						Name:          "myMultiString",
						Type:          "string",
						Default:       "default",
						Values:        []string{"default", "ubuntu"},
						FileReference: testutils.CreateFileReference(6, 3, 11, 11),
					},
					{
						Name:          "myNumber",
						Type:          "number",
						Default:       "2",
						Values:        []string{"1", "2", "4", "8", "16"},
						FileReference: testutils.CreateFileReference(12, 3, 20, 7),
					},
					{
						Name:          "myBoolean",
						Type:          "boolean",
						Default:       "true",
						FileReference: testutils.CreateFileReference(21, 3, 23, 16),
					},
					{
						Name: "myObject",
						Type: "object",
						Default: map[string]any{
							"foo":    "FOO",
							"bar":    "BAR",
							"things": []string{"one", "two", "three"},
							"nested": map[string]any{
								"one":   "apple",
								"two":   "pear",
								"count": 3,
							},
						},
						FileReference: testutils.CreateFileReference(24, 3, 36, 15),
					},
					{
						Name: "myStep",
						Type: "step",
						Default: map[string]any{
							"script": "echo my step",
						},
						FileReference: testutils.CreateFileReference(37, 3, 40, 25),
					},
					{
						Name: "myStepList",
						Type: "stepList",
						Default: []map[string]any{
							{
								"script": "echo step one",
							},
							{
								"script": "echo step two",
							},
						},
						FileReference: testutils.CreateFileReference(41, 3, 45, 28),
					},
				},
			},
		},
		{
			name:     "stages",
			filename: "../../../test/fixtures/azure/stages.yaml",
			expectedPipeline: &models.Pipeline{
				Name: "stages",
				Stages: &models.Stages{
					Stages: &[]models.Stage{
						{
							Stage:         "BuildWin",
							DisplayName:   "Build for Windows",
							FileReference: testutils.CreateFileReference(4, 3, 5, 33),
						},
						{
							Stage:         "BuildMac",
							DisplayName:   "Build for Mac",
							DependsOn:     []models.DependsOn{},
							FileReference: testutils.CreateFileReference(6, 3, 8, 14),
						},
					},
					TemplateStages: &[]models.TemplateStage{
						{
							Template: models.Template{
								Template: "stages/build.yml",
								Parameters: map[string]any{
									"param": "value",
								},
							},
							FileReference: testutils.CreateFileReference(10, 3, 12, 17),
						},
						{
							Template: models.Template{
								Template: "stages/test.yml",
								Parameters: map[string]any{
									"name":     "Full",
									"testFile": "tests/fullSuite.js",
								},
							},
							FileReference: testutils.CreateFileReference(14, 3, 17, 33),
						},
					},
					FileReference: testutils.CreateFileReference(3, -1, 17, 33),
				},
			},
		},
		{
			name:     "jobs",
			filename: "../../../test/fixtures/azure/jobs.yaml",
			expectedPipeline: &models.Pipeline{
				Name: "Jobs",
				Jobs: &models.Jobs{
					CIJobs: &[]models.CIJob{
						{
							Job: "MyJob",
							BaseJob: models.BaseJob{
								DisplayName:     "My First Job",
								ContinueOnError: true,
								Workspace: &models.Workspace{
									Clean: "outputs",
								},
								Steps: &models.Steps{
									{
										Script:        "echo My first job",
										FileReference: testutils.CreateFileReference(10, 5, 10, 30),
									},
								},
							},
							FileReference: testutils.CreateFileReference(4, 3, 10, 30),
						},
					},
					DeploymentJobs: &[]models.DeploymentJob{
						{
							Deployment: "DeployWeb",
							BaseJob: models.BaseJob{
								DisplayName: "deploy Web App",
								Pool: &models.Pool{
									VmImage:       "ubuntu-latest",
									FileReference: testutils.CreateFileReference(13, 5, 14, 27),
								},
							},
							Environment: &models.DeploymentEnvironment{
								Name:          "smarthotel-dev",
								FileReference: testutils.CreateFileReference(16, 16, 16, 30),
							},
							Strategy: &models.DeploymentStrategy{
								RunOnce: &models.BaseDeploymentStrategy{
									Deploy: &models.DeploymentHook{
										Steps: &models.Steps{
											{
												Script:        "echo my first deployment",
												FileReference: testutils.CreateFileReference(22, 11, 22, 43),
											},
										},
									},
								},
							},
							FileReference: testutils.CreateFileReference(11, 3, 22, 43),
						},
					},
					TemplateJobs: &[]models.TemplateJob{
						{
							Template: models.Template{
								Template: "jobs/build.yml",
								Parameters: map[string]any{
									"name": "macOS",
									"pool": map[string]any{
										"vmImage": "macOS-latest",
									},
								},
							},
							FileReference: testutils.CreateFileReference(23, 3, 27, 28),
						},
						{
							Template: models.Template{
								Template: "jobs/build.yml",
								Parameters: map[string]any{
									"name": "Linux",
									"pool": map[string]any{
										"vmImage": "ubuntu-latest",
									},
								},
							},
							FileReference: testutils.CreateFileReference(29, 3, 33, 29),
						},
						{
							Template: models.Template{
								Template: "jobs/build.yml",
								Parameters: map[string]any{
									"name": "Windows",
									"pool": map[string]any{
										"vmImage": "windows-latest",
									},
									"sign": true,
								},
							},
							FileReference: testutils.CreateFileReference(34, 3, 39, 15),
						},
					},
					FileReference: testutils.CreateFileReference(3, -1, 39, 15),
				},
			},
		},
		{
			name:     "steps",
			filename: "../../../test/fixtures/azure/steps.yaml",
			expectedPipeline: &models.Pipeline{
				Name: "steps",
				Steps: &models.Steps{
					{
						Bash:        "which bash\necho Hello $name\n",
						DisplayName: "Multiline Bash script",
						Env: &models.EnvironmentVariablesRef{
							EnvironmentVariables: map[string]any{
								"name": "Microsoft",
							},
							FileReference: testutils.CreateFileReference(7, 5, 8, 20),
						},
						FileReference: testutils.CreateFileReference(3, 3, 8, 20),
					},
					{
						Checkout:           "self",
						Submodules:         "true",
						PersistCredentials: true,
						FileReference:      testutils.CreateFileReference(9, 3, 11, 27),
					},
					{
						Download:      "current",
						Artifact:      "WebApp",
						Patterns:      "**/.js",
						DisplayName:   "Download artifact WebApp",
						FileReference: testutils.CreateFileReference(12, 3, 15, 40),
					},
					{
						DownloadBuild: "current",
						Artifact:      "WebApp",
						Path:          "build",
						Patterns:      "**/.js",
						DisplayName:   "Download artifact WebApp",
						FileReference: testutils.CreateFileReference(16, 3, 20, 40),
					},
					{
						GetPackage:    "packageID",
						Path:          "dist",
						FileReference: testutils.CreateFileReference(21, 3, 22, 13),
					},
					{
						Powershell:       "Write-Host Hello $(name)",
						DisplayName:      "Say hello",
						Name:             "firstStep",
						WorkingDirectory: "$(build.sourcesDirectory)",
						FailOnStderr:     true,
						Env: &models.EnvironmentVariablesRef{
							EnvironmentVariables: map[string]any{
								"name": "Microsoft",
							},
							FileReference: testutils.CreateFileReference(28, 5, 29, 20),
						},
						FileReference: testutils.CreateFileReference(23, 3, 29, 20),
					},
					{
						Publish:       "$(Build.SourcesDirectory)/build",
						Artifact:      "WebApp",
						DisplayName:   "Publish artifact WebApp",
						FileReference: testutils.CreateFileReference(30, 3, 32, 39),
					},
					{
						Pwsh:             "Write-Host Hello $(name)",
						DisplayName:      "Say hello",
						Name:             "firstStep",
						WorkingDirectory: "$(build.sourcesDirectory)",
						FailOnStderr:     true,
						Env: &models.EnvironmentVariablesRef{
							EnvironmentVariables: map[string]any{
								"name": "Microsoft",
							},
							FileReference: testutils.CreateFileReference(38, 5, 39, 20),
						},
						FileReference: testutils.CreateFileReference(33, 3, 39, 20),
					},
					{
						ReviewApp:     "review",
						FileReference: testutils.CreateFileReference(40, 3, 40, 20),
					},
					{
						Script:        "echo This runs in the default shell on any machine",
						FileReference: testutils.CreateFileReference(41, 3, 41, 61),
					},
					{
						Task:             "VSBuild@1",
						DisplayName:      "Build",
						TimeoutInMinutes: 120,
						Inputs: map[string]any{
							"solution": "**\\*.sln",
						},
						FileReference: testutils.CreateFileReference(42, 3, 46, 23),
					},
					{
						Template: models.Template{
							Template: "steps/build.yml",
							Parameters: map[string]any{
								"key": "value",
							},
						},
						FileReference: testutils.CreateFileReference(47, 3, 49, 15),
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			loader := &AzureLoader{}
			pipeline, err := loader.Load(testutils.ReadFile(testCase.filename))

			if err != testCase.expectedError {
				t.Errorf("Expected error: %v, got: %v", testCase.expectedError, err)
			}

			changelog, _ := diff.Diff(pipeline, testCase.expectedPipeline)

			if len(changelog) > 0 {
				t.Errorf("Loader result is not as expected:")
				for _, change := range changelog {
					t.Errorf("field: %s, got: %v, expected: %v", strings.Join(change.Path, "."), change.From, change.To)
				}
			}
		})
	}
}
