package bitbucket

import (
	"testing"

	bbModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name              string
		bitbucketPipeline *bbModels.Pipeline
		expectedPipeline  *models.Pipeline
	}{
		{
			name:              "Pipeline is nil",
			bitbucketPipeline: nil,
			expectedPipeline:  nil,
		},
		{
			name: "Pipeline with default job",
			bitbucketPipeline: &bbModels.Pipeline{
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
										FileReference: testutils.CreateFileReference(7, 13, 12, 25),
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
										FileReference: testutils.CreateFileReference(14, 13, 19, 21),
									},
								},
							},
						},
					},
				},
			},
			expectedPipeline: &models.Pipeline{
				Defaults: &models.Defaults{},
				Jobs: []*models.Job{
					{
						FileReference: testutils.CreateFileReference(7, 13, 19, 21),
						ID:            utils.GetPtr("job-default"),
						Name:          utils.GetPtr("default"),
						Steps: []*models.Step{
							{
								Type: "shell",
								Name: utils.GetPtr("Build and Test"),
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("npm install\nnpm test"),
									FileReference: testutils.CreateFileReference(11, 17, 12, 25),
								},
								FileReference: testutils.CreateFileReference(7, 13, 12, 25),
							},
							{
								Type: "shell",
								Name: utils.GetPtr("Code linting"),
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("npm install eslint\nnpx eslint ."),
									FileReference: testutils.CreateFileReference(16, 17, 17, 29),
								},
								FileReference: testutils.CreateFileReference(14, 13, 19, 21),
							},
						},
					},
				},
			},
		},
		{
			name: "Pipeline with pull-request jobs",
			bitbucketPipeline: &bbModels.Pipeline{
				Pipelines: &bbModels.BuildPipelines{
					PullRequests: &bbModels.StepMap{
						"**": []*bbModels.Step{
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: utils.GetPtr("Build and Test"),
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
									FileReference: testutils.CreateFileReference(7, 13, 12, 25),
								},
							},
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: utils.GetPtr("Deploy"),
										Script: []*bbModels.Script{
											{
												String:        utils.GetPtr("deploy.sh"),
												FileReference: testutils.CreateFileReference(11, 17, 11, 28),
											},
										},
									},
									FileReference: testutils.CreateFileReference(7, 13, 12, 25),
								},
							},
						},
						"master": []*bbModels.Step{
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: utils.GetPtr("Test"),
										Script: []*bbModels.Script{
											{
												String:        utils.GetPtr("npm test"),
												FileReference: testutils.CreateFileReference(12, 17, 12, 25),
											},
										},
									},
									FileReference: testutils.CreateFileReference(12, 17, 12, 25),
								},
							},
						},
					},
				},
			},
			expectedPipeline: &models.Pipeline{
				Defaults: &models.Defaults{},
				Jobs: []*models.Job{
					{
						FileReference: testutils.CreateFileReference(7, 13, 12, 25),
						ID:            utils.GetPtr("job-**"),
						Name:          utils.GetPtr("**"),
						Steps: []*models.Step{
							{
								Type: "shell",
								Name: utils.GetPtr("Build and Test"),
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("npm install\nnpm test"),
									FileReference: testutils.CreateFileReference(11, 17, 12, 25),
								},
								FileReference: testutils.CreateFileReference(7, 13, 12, 25),
							},
							{
								Type: "shell",
								Name: utils.GetPtr("Deploy"),
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("deploy.sh"),
									FileReference: testutils.CreateFileReference(11, 17, 11, 28),
								},
								FileReference: testutils.CreateFileReference(7, 13, 12, 25),
							},
						},
					},
					{
						FileReference: testutils.CreateFileReference(12, 17, 12, 25),
						ID:            utils.GetPtr("job-master"),
						Name:          utils.GetPtr("master"),
						Steps: []*models.Step{
							{
								Type: "shell",
								Name: utils.GetPtr("Test"),
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("npm test"),
									FileReference: testutils.CreateFileReference(12, 17, 12, 25),
								},
								FileReference: testutils.CreateFileReference(12, 17, 12, 25),
							},
						},
					},
				},
			},
		},
		{
			name: "Pipeline with options and definitions",
			bitbucketPipeline: &bbModels.Pipeline{
				Image: &bbModels.Image{
					ImageData: &bbModels.ImageData{
						Name: utils.GetPtr("node:8"),
					},
				},
				Options: &bbModels.GlobalSettings{
					Docker:  utils.GetPtr(true),
					MaxTime: utils.GetPtr(int64(60)),
					Size:    utils.GetPtr(bbModels.X2),
				},
				Definitions: &bbModels.Definitions{
					Steps: []*bbModels.Step{
						{
							Step: &bbModels.ExecutionUnitRef{
								ExecutionUnit: &bbModels.ExecutionUnit{
									Name: utils.GetPtr("Build and Test"),
									Script: []*bbModels.Script{
										{
											PipeToExecute: &bbModels.PipeToExecute{
												Pipe: &bbModels.Pipe{
													String:        utils.GetPtr("npm install"),
													FileReference: testutils.CreateFileReference(7, 13, 11, 23),
												},
												Variables: &bbModels.EnvironmentVariablesRef{
													EnvironmentVariables: models.EnvironmentVariables{
														"NPM_TOKEN": "secret",
													},
													FileReference: testutils.CreateFileReference(11, 17, 11, 28),
												},
											},
											FileReference: testutils.CreateFileReference(7, 13, 12, 25),
										},
									},
								},
								FileReference: testutils.CreateFileReference(7, 13, 12, 25),
							},
						},
					},
				},
				Pipelines: &bbModels.BuildPipelines{
					Custom: &bbModels.StepMap{
						"install": []*bbModels.Step{
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: utils.GetPtr("Build and Test"),
										Script: []*bbModels.Script{
											{
												PipeToExecute: &bbModels.PipeToExecute{
													Pipe: &bbModels.Pipe{
														String:        utils.GetPtr("npm install"),
														FileReference: testutils.CreateFileReference(10, 17, 10, 28),
													},
													Variables: &bbModels.EnvironmentVariablesRef{
														EnvironmentVariables: models.EnvironmentVariables{
															"NPM_TOKEN": "secret",
														},
														FileReference: testutils.CreateFileReference(11, 17, 11, 28),
													},
												},
												FileReference: testutils.CreateFileReference(7, 13, 12, 25),
											},
										},
									},
									FileReference: testutils.CreateAliasFileReference(8, 13, 12, 25, true),
								},
							},
						},
					},
				},
			},
			expectedPipeline: &models.Pipeline{
				Defaults: &models.Defaults{
					Runner: &models.Runner{
						DockerMetadata: &models.DockerMetadata{
							Image: utils.GetPtr("node:8"),
						},
					},
					Settings: &map[string]any{
						"docker":   utils.GetPtr(true),
						"max-time": utils.GetPtr(int64(60)),
						"size":     utils.GetPtr(bbModels.X2),
					},
				},
				Jobs: []*models.Job{
					{
						FileReference: testutils.CreateFileReference(8, 13, 12, 25),
						ID:            utils.GetPtr("job-install"),
						Name:          utils.GetPtr("install"),
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("Build and Test"),
								Type: "task",
								Task: &models.Task{
									Name:        utils.GetPtr("npm install"),
									VersionType: "none",
								},
								EnvironmentVariables: &models.EnvironmentVariablesRef{
									EnvironmentVariables: models.EnvironmentVariables{
										"NPM_TOKEN": "secret",
									},
									FileReference: testutils.CreateFileReference(11, 17, 11, 28),
								},
								FileReference: testutils.CreateAliasFileReference(8, 13, 12, 25, true),
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			parser := BitbucketParser{}

			pipeline, err := parser.Parse(testCase.bitbucketPipeline)
			assert.NoError(t, err)
			testutils.DeepCompare(t, testCase.expectedPipeline, pipeline)
		})
	}
}
