package bitbucket

import (
	"strings"
	"testing"

	bbModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/r3labs/diff/v3"
)

func Ptr[T any](v T) *T {
	return &v
}

func TestLoad(t *testing.T) {
	testCases := []struct {
		name             string
		filename         string
		expectedPipeline *bbModels.Pipeline
		expectedError    error
	}{
		{
			name:     "parallel steps",
			filename: "../../../test/fixtures/bitbucket/parallel-steps.yml",
			expectedPipeline: &bbModels.Pipeline{
				Image: &bbModels.Image{
					Name: "node:16",
				},
				Pipelines: &bbModels.BuildPipelines{
					Default: []*bbModels.Step{
						{
							Parallel: []*bbModels.ParallelSteps{
								{
									Step: &bbModels.ExecutionUnitRef{
										ExecutionUnit: &bbModels.ExecutionUnit{
											Name: "Build and Test",
											Caches: []string{
												"node",
											},
											Script: []bbModels.Script{
												{
													String:        "npm install",
													FileReference: testutils.CreateFileReference(11, 17, 11, 28),
												},
												{
													String:        "npm test",
													FileReference: testutils.CreateFileReference(12, 17, 12, 25),
												},
											},
										},
										FileReference: testutils.CreateFileReference(7, 13, 12, 25),
									},
								},
								{
									Step: &bbModels.ExecutionUnitRef{
										ExecutionUnit: &bbModels.ExecutionUnit{
											Name: "Code linting",
											Caches: []string{
												"node",
											},
											Script: []bbModels.Script{
												{
													String:        "npm install eslint",
													FileReference: testutils.CreateFileReference(16, 17, 16, 35),
												},
												{
													String:        "npx eslint .",
													FileReference: testutils.CreateFileReference(17, 17, 17, 29),
												},
											},
										},
										FileReference: testutils.CreateFileReference(14, 13, 19, 21),
									},
								},
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:     "sync steps",
			filename: "../../../test/fixtures/bitbucket/sync-steps.yml",
			expectedPipeline: &bbModels.Pipeline{
				Image: &bbModels.Image{
					Name: "node:16",
				},
				Pipelines: &bbModels.BuildPipelines{
					PullRequests: &bbModels.StepMap{
						"**": {
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: "Build and Test",
										Caches: []string{
											"node",
										},
										Script: []bbModels.Script{
											{
												String:        "npm install",
												FileReference: testutils.CreateFileReference(11, 15, 11, 26),
											},
											{
												String:        "npm test",
												FileReference: testutils.CreateFileReference(12, 15, 12, 23),
											},
										},
									},
									FileReference: testutils.CreateFileReference(6, 9, 12, 23),
								},
							},
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: "Code linting",
										Caches: []string{
											"node",
										},
										Script: []bbModels.Script{
											{
												String:        "npm install eslint",
												FileReference: testutils.CreateFileReference(16, 15, 16, 33),
											},
											{
												String:        "npx eslint .",
												FileReference: testutils.CreateFileReference(17, 15, 17, 27),
											},
										},
									},
									FileReference: testutils.CreateFileReference(13, 9, 19, 19),
								},
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:     "multiple pipelines definitions",
			filename: "../../../test/fixtures/bitbucket/multiple-pipelines-types.yml",
			expectedPipeline: &bbModels.Pipeline{
				Image: &bbModels.Image{
					Name: "node:16",
				},
				Pipelines: &bbModels.BuildPipelines{
					Custom: &bbModels.StepMap{
						"notify": {
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: "Notify Teams",
										Caches: []string{
											"node",
										},
										Script: []bbModels.Script{
											{
												String:        "npx notify -s \"deployment\"",
												FileReference: testutils.CreateFileReference(11, 15, 11, 41),
											},
										},
									},
									FileReference: testutils.CreateFileReference(6, 9, 11, 41),
								},
							},
						},
					},
					Branches: &bbModels.StepMap{
						"master": {
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: "step 1",
									},
									FileReference: testutils.CreateFileReference(14, 9, 15, 23),
								},
							},
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: "step 2",
									},
									FileReference: testutils.CreateFileReference(16, 9, 17, 23),
								},
							},
							{
								Parallel: []*bbModels.ParallelSteps{
									{
										Step: &bbModels.ExecutionUnitRef{
											ExecutionUnit: &bbModels.ExecutionUnit{
												Name: "step 3",
											},
											FileReference: testutils.CreateFileReference(20, 15, 20, 27),
										},
									},
									{
										Step: &bbModels.ExecutionUnitRef{
											ExecutionUnit: &bbModels.ExecutionUnit{
												Name: "step 4",
											},
											FileReference: testutils.CreateFileReference(22, 15, 22, 27),
										},
									},
								},
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:     "step with image",
			filename: "../../../test/fixtures/bitbucket/image.yml",
			expectedPipeline: &bbModels.Pipeline{
				Image: &bbModels.Image{
					Name: "node:16",
				},
				Pipelines: &bbModels.BuildPipelines{
					Default: []*bbModels.Step{
						{
							Step: &bbModels.ExecutionUnitRef{
								ExecutionUnit: &bbModels.ExecutionUnit{
									Image: &bbModels.Image{
										ImageWithCustomUser: &bbModels.ImageWithCustomUser{
											Name:     "node:16",
											Email:    Ptr("test@test.com"),
											Username: Ptr("test"),
											Password: Ptr("test"),
											Aws: &bbModels.Aws{
												AccessKey: "123456",
												SecretKey: "7891011",
											},
										},
									},
								},
								FileReference: testutils.CreateFileReference(5, 7, 13, 32),
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:     "additional config: clone, options, services, variables",
			filename: "../../../test/fixtures/bitbucket/config-options.yml",
			expectedPipeline: &bbModels.Pipeline{
				Clone: &bbModels.Clone{
					Depth:   1,
					LFS:     Ptr(true),
					Enabled: Ptr(true),
				},
				Options: &bbModels.GlobalSettings{
					MaxTime: Ptr(int64(30)),
					Docker:  Ptr(true),
					Size:    Ptr(bbModels.X1),
				},
				Definitions: &bbModels.Definitions{
					Caches: &bbModels.Caches{
						"custom-npm": "node_modules",
					},
					Services: map[string]*bbModels.Service{
						"service": {
							Memory: Ptr(int64(128)),
							Image: &bbModels.Image{
								Name: "node:16",
							},
							Variables: &bbModels.EnvironmentVariablesRef{
								EnvironmentVariables: models.EnvironmentVariables{
									"TEST":  "development",
									"TEST2": "production",
								},
								FileReference: testutils.CreateFileReference(16, 9, 18, 26),
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:     "definitions",
			filename: "../../../test/fixtures/bitbucket/definitions.yml",
			expectedPipeline: &bbModels.Pipeline{
				Definitions: &bbModels.Definitions{
					Steps: []*bbModels.Step{
						{
							Step: &bbModels.ExecutionUnitRef{
								ExecutionUnit: &bbModels.ExecutionUnit{
									Name: "scripts step",
									Script: []bbModels.Script{
										{
											String:        "echo \"hello world\"",
											FileReference: testutils.CreateFileReference(6, 13, 6, 31),
										},
									},
									AfterScript: []bbModels.Script{
										{
											String:        "echo \"goodbye world\"",
											FileReference: testutils.CreateFileReference(8, 13, 8, 33),
										},
										{
											PipeToExecute: &bbModels.PipeToExecute{
												Pipe: "notify",
												Variables: bbModels.EnvironmentVariablesRef{
													EnvironmentVariables: models.EnvironmentVariables{
														"FOO": "bar",
													},
													FileReference: testutils.CreateFileReference(10, 15, 11, 23),
												},
											},
										},
									},
								},
								FileReference: testutils.CreateFileReference(3, 7, 11, 23),
							},
						},
						{
							Step: &bbModels.ExecutionUnitRef{
								ExecutionUnit: &bbModels.ExecutionUnit{
									Name: "artifacts step",
									Caches: []string{
										"package.json",
									},
								},
								FileReference: testutils.CreateFileReference(12, 7, 15, 25),
							},
						},
						{
							Step: &bbModels.ExecutionUnitRef{
								ExecutionUnit: &bbModels.ExecutionUnit{
									Name: "shared artifact step",
									Artifacts: &bbModels.Artifacts{
										SharedStepFiles: &bbModels.SharedStepFiles{
											Download: Ptr(false),
											Paths: []string{
												"dist/*",
												"package-lock.json",
											},
										},
									},
								},
								FileReference: testutils.CreateFileReference(16, 7, 22, 32),
							},
						},
						{
							Parallel: []*bbModels.ParallelSteps{
								{
									Step: &bbModels.ExecutionUnitRef{
										ExecutionUnit: &bbModels.ExecutionUnit{
											Name:    "parallel step 1",
											Trigger: Ptr(bbModels.MANUAL),
										},
										FileReference: testutils.CreateFileReference(25, 13, 26, 28),
									},
								},
								{
									Step: &bbModels.ExecutionUnitRef{
										ExecutionUnit: &bbModels.ExecutionUnit{
											Name:    "parallel step 2",
											Trigger: Ptr(bbModels.AUTOMATIC),
										},
										FileReference: testutils.CreateFileReference(28, 13, 29, 31),
									},
								},
							},
						},
					},
				},
				Pipelines: &bbModels.BuildPipelines{
					Custom: &bbModels.StepMap{
						"test": {
							{
								Variables: []*bbModels.CustomStepVariable{
									{
										Name:          Ptr("Username"),
										FileReference: testutils.CreateFileReference(34, 13, 34, 27),
									},
									{
										Name:          Ptr("Role"),
										Default:       Ptr("admin"),
										FileReference: testutils.CreateFileReference(35, 13, 36, 27),
									},
									{
										Name:    Ptr("Region"),
										Default: Ptr("ap-southeast-2"),
										AllowedValues: []string{
											"ap-southeast-2",
											"us-east-1",
											"us-west-2",
										},
										FileReference: testutils.CreateFileReference(37, 13, 42, 26),
									},
								},
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:     "alias nodes",
			filename: "../../../test/fixtures/bitbucket/alias-nodes.yml",
			expectedPipeline: &bbModels.Pipeline{
				Definitions: &bbModels.Definitions{
					Caches: &bbModels.Caches{
						"cypress": "/root/.cache/Cypress",
					},
					Services: map[string]*bbModels.Service{
						"docker": {
							Memory: Ptr(int64(2048)),
						},
					},
					Steps: []*bbModels.Step{
						{
							Step: &bbModels.ExecutionUnitRef{
								ExecutionUnit: &bbModels.ExecutionUnit{
									Name: "Install and build",
									Script: []bbModels.Script{
										{
											String:        "yarn build",
											FileReference: testutils.CreateFileReference(11, 13, 11, 23),
										},
									},
									Artifacts: &bbModels.Artifacts{
										Paths: []string{
											"dist/**",
										},
									},
								},
								FileReference: testutils.CreateFileReference(8, 7, 13, 20),
							},
						},
					},
				},
				Pipelines: &bbModels.BuildPipelines{
					PullRequests: &bbModels.StepMap{
						"*": {
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: "Install and build",
										Script: []bbModels.Script{
											{
												String:        "yarn build",
												FileReference: testutils.CreateFileReference(11, 13, 11, 23),
											},
										},
										Artifacts: &bbModels.Artifacts{
											Paths: []string{
												"dist/**",
											},
										},
									},
									FileReference: testutils.CreateFileReference(8, 7, 13, 20),
								},
							},
						},
						"**": {
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: "Install and build",
										Script: []bbModels.Script{
											{
												String:        "yarn build",
												FileReference: testutils.CreateFileReference(11, 13, 11, 23),
											},
										},
										Artifacts: &bbModels.Artifacts{
											Paths: []string{
												"dist/**",
											},
										},
									},
									FileReference: testutils.CreateFileReference(20, 9, 21, 28),
								},
							},
						},
					},
				},
			},
			expectedError: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			loader := &BitbucketLoader{}
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
