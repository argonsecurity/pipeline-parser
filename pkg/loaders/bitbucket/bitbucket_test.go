package bitbucket

import (
	"strings"
	"testing"

	bbModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/r3labs/diff/v3"
)

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
					ImageData: &bbModels.ImageData{
						Name: utils.GetPtr("node:16"),
					},
				},
				Pipelines: &bbModels.BuildPipelines{
					Default: []*bbModels.Step{
						{
							Parallel: []*bbModels.ParallelSteps{
								{
									Step: &bbModels.ExecutionUnitRef{
										ExecutionUnit: &bbModels.ExecutionUnit{
											Name: utils.GetPtr("Build and Test"),
											Caches: []*string{
												utils.GetPtr("node"),
											},
											Script: []*bbModels.Script{
												{
													String:        utils.GetPtr("npm install"),
													FileReference: testutils.CreateFileReference(11, 17, 11, 28),
												},
												{
													String:        utils.GetPtr("npm test"),
													FileReference: testutils.CreateFileReference(12, 17, 12, 25),
												},
											},
										},
										FileReference: testutils.CreateFileReference(6, 11, 12, 25),
									},
								},
								{
									Step: &bbModels.ExecutionUnitRef{
										ExecutionUnit: &bbModels.ExecutionUnit{
											Name: utils.GetPtr("Code linting"),
											Caches: []*string{
												utils.GetPtr("node"),
											},
											Script: []*bbModels.Script{
												{
													String:        utils.GetPtr("npm install eslint"),
													FileReference: testutils.CreateFileReference(16, 17, 16, 35),
												},
												{
													String:        utils.GetPtr("npx eslint ."),
													FileReference: testutils.CreateFileReference(17, 17, 17, 29),
												},
											},
										},
										FileReference: testutils.CreateFileReference(13, 11, 19, 21),
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
					ImageData: &bbModels.ImageData{
						Name: utils.GetPtr("node:16"),
					},
				},
				Pipelines: &bbModels.BuildPipelines{
					PullRequests: &bbModels.StepMap{
						"**": {
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: utils.GetPtr("Build and Test"),
										Caches: []*string{
											utils.GetPtr("node"),
										},
										Script: []*bbModels.Script{
											{
												String:        utils.GetPtr("npm install"),
												FileReference: testutils.CreateFileReference(11, 15, 11, 26),
											},
											{
												String:        utils.GetPtr("npm test"),
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
										Name: utils.GetPtr("Code linting"),
										Caches: []*string{
											utils.GetPtr("node"),
										},
										Script: []*bbModels.Script{
											{
												String:        utils.GetPtr("npm install eslint"),
												FileReference: testutils.CreateFileReference(16, 15, 16, 33),
											},
											{
												String:        utils.GetPtr("npx eslint ."),
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
					ImageData: &bbModels.ImageData{
						Name: utils.GetPtr("node:16"),
					},
				},
				Pipelines: &bbModels.BuildPipelines{
					Custom: &bbModels.StepMap{
						"notify": {
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: utils.GetPtr("Notify Teams"),
										Caches: []*string{
											utils.GetPtr("node"),
										},
										Script: []*bbModels.Script{
											{
												String:        utils.GetPtr("npx notify -s \"deployment\""),
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
										Name: utils.GetPtr("step 1"),
									},
									FileReference: testutils.CreateFileReference(14, 9, 15, 23),
								},
							},
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: utils.GetPtr("step 2"),
									},
									FileReference: testutils.CreateFileReference(16, 9, 17, 23),
								},
							},
							{
								Parallel: []*bbModels.ParallelSteps{
									{
										Step: &bbModels.ExecutionUnitRef{
											ExecutionUnit: &bbModels.ExecutionUnit{
												Name: utils.GetPtr("step 3"),
											},
											FileReference: testutils.CreateFileReference(19, 13, 20, 27),
										},
									},
									{
										Step: &bbModels.ExecutionUnitRef{
											ExecutionUnit: &bbModels.ExecutionUnit{
												Name: utils.GetPtr("step 4"),
											},
											FileReference: testutils.CreateFileReference(21, 13, 22, 27),
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
					ImageData: &bbModels.ImageData{
						Name: utils.GetPtr("node:16"),
					},
				},
				Pipelines: &bbModels.BuildPipelines{
					Default: []*bbModels.Step{
						{
							Step: &bbModels.ExecutionUnitRef{
								ExecutionUnit: &bbModels.ExecutionUnit{
									Image: &bbModels.Image{
										ImageData: &bbModels.ImageData{
											Name:     utils.GetPtr("node:16"),
											Email:    utils.GetPtr("test@test.com"),
											Username: utils.GetPtr("test"),
											Password: utils.GetPtr("test"),
											Aws: &bbModels.Aws{
												AccessKey: utils.GetPtr("123456"),
												SecretKey: utils.GetPtr("7891011"),
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
					LFS:     utils.GetPtr(true),
					Enabled: utils.GetPtr(true),
				},
				Options: &bbModels.GlobalSettings{
					MaxTime: utils.GetPtr(int64(30)),
					Docker:  utils.GetPtr(true),
					Size:    utils.GetPtr(bbModels.X1),
				},
				Definitions: &bbModels.Definitions{
					Caches: &bbModels.Caches{
						"custom-npm": utils.GetPtr("node_modules"),
					},
					Services: map[string]*bbModels.Service{
						"service": {
							Memory: utils.GetPtr(int64(128)),
							Image: &bbModels.Image{
								ImageData: &bbModels.ImageData{
									Name: utils.GetPtr("node:16"),
								},
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
									Name: utils.GetPtr("scripts step"),
									Script: []*bbModels.Script{
										{
											String:        utils.GetPtr("echo \"hello world\""),
											FileReference: testutils.CreateFileReference(6, 13, 6, 31),
										},
									},
									AfterScript: []*bbModels.Script{
										{
											String:        utils.GetPtr("echo \"goodbye world\""),
											FileReference: testutils.CreateFileReference(8, 13, 8, 33),
										},
										{
											PipeToExecute: &bbModels.PipeToExecute{
												Pipe: &bbModels.Pipe{
													String:        utils.GetPtr("notify"),
													FileReference: testutils.CreateFileReference(9, 19, 9, 25),
												},
												Variables: &bbModels.EnvironmentVariablesRef{
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
									Name: utils.GetPtr("artifacts step"),
									Caches: []*string{
										utils.GetPtr("package.json"),
									},
								},
								FileReference: testutils.CreateFileReference(12, 7, 15, 25),
							},
						},
						{
							Step: &bbModels.ExecutionUnitRef{
								ExecutionUnit: &bbModels.ExecutionUnit{
									Name: utils.GetPtr("shared artifact step"),
									Artifacts: &bbModels.Artifacts{
										SharedStepFiles: &bbModels.SharedStepFiles{
											Download: utils.GetPtr(false),
											Paths: []*string{
												utils.GetPtr("dist/*"),
												utils.GetPtr("package-lock.json"),
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
											Name:    utils.GetPtr("parallel step 1"),
											Trigger: utils.GetPtr(bbModels.MANUAL),
										},
										FileReference: testutils.CreateFileReference(24, 11, 26, 28),
									},
								},
								{
									Step: &bbModels.ExecutionUnitRef{
										ExecutionUnit: &bbModels.ExecutionUnit{
											Name:    utils.GetPtr("parallel step 2"),
											Trigger: utils.GetPtr(bbModels.AUTOMATIC),
										},
										FileReference: testutils.CreateFileReference(27, 11, 29, 31),
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
										Name:          utils.GetPtr("Username"),
										FileReference: testutils.CreateFileReference(34, 13, 34, 27),
									},
									{
										Name:          utils.GetPtr("Role"),
										Default:       utils.GetPtr("admin"),
										FileReference: testutils.CreateFileReference(35, 13, 36, 27),
									},
									{
										Name:    utils.GetPtr("Region"),
										Default: utils.GetPtr("ap-southeast-2"),
										AllowedValues: []*string{
											utils.GetPtr("ap-southeast-2"),
											utils.GetPtr("us-east-1"),
											utils.GetPtr("us-west-2"),
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
						"cypress": utils.GetPtr("/root/.cache/Cypress"),
					},
					Services: map[string]*bbModels.Service{
						"docker": {
							Memory: utils.GetPtr(int64(2048)),
						},
					},
					Steps: []*bbModels.Step{
						{
							Step: &bbModels.ExecutionUnitRef{
								ExecutionUnit: &bbModels.ExecutionUnit{
									Name: utils.GetPtr("Install and build"),
									Script: []*bbModels.Script{
										{
											String:        utils.GetPtr("yarn build"),
											FileReference: testutils.CreateFileReference(11, 13, 11, 23),
										},
									},
									Artifacts: &bbModels.Artifacts{
										Paths: []*string{
											utils.GetPtr("dist/**"),
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
										Name: utils.GetPtr("Install and build"),
										Script: []*bbModels.Script{
											{
												String:        utils.GetPtr("yarn build"),
												FileReference: testutils.CreateFileReference(11, 13, 11, 23),
											},
										},
										Artifacts: &bbModels.Artifacts{
											Paths: []*string{
												utils.GetPtr("dist/**"),
											},
										},
									},
									FileReference: testutils.CreateAliasFileReference(8, 7, 13, 20, true),
								},
							},
						},
						"**": {
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: utils.GetPtr("merge test"),
										Script: []*bbModels.Script{
											{
												String:        utils.GetPtr("yarn build"),
												FileReference: testutils.CreateFileReference(11, 13, 11, 23),
											},
										},
										Artifacts: &bbModels.Artifacts{
											Paths: []*string{
												utils.GetPtr("dist/**"),
											},
										},
									},
									FileReference: testutils.CreateAliasFileReference(20, 9, 22, 27, true),
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
