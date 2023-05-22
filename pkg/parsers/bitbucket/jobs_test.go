package bitbucket

import (
	"testing"

	bitbucketModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestJobsParse(t *testing.T) {
	testCases := []struct {
		name              string
		bitbucketPipeline *bitbucketModels.Pipeline
		expectedJobs      []*models.Job
	}{
		{
			name:              "Pipeline is nil",
			bitbucketPipeline: nil,
			expectedJobs:      nil,
		},
		{
			name: "Pipeline has no jobs",
			bitbucketPipeline: &bitbucketModels.Pipeline{
				Image: &bitbucketModels.Image{
					ImageData: &bitbucketModels.ImageData{
						Name: utils.GetPtr("node:10.15.3"),
					},
				},
			},
			expectedJobs: nil,
		},
		{
			name: "Pipeline has default job",
			bitbucketPipeline: &bitbucketModels.Pipeline{
				Pipelines: &bitbucketModels.BuildPipelines{
					Default: []*bitbucketModels.Step{
						{
							Step: &bitbucketModels.ExecutionUnitRef{
								ExecutionUnit: &bitbucketModels.ExecutionUnit{
									Script: []*bitbucketModels.Script{
										{
											String:        utils.GetPtr("echo 'hello world'"),
											FileReference: testutils.CreateFileReference(1, 2, 3, 4),
										},
									},
								},
								FileReference: testutils.CreateFileReference(5, 6, 7, 8),
							},
						},
					},
				},
			},
			expectedJobs: []*models.Job{
				{
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
					ID:            utils.GetPtr("job-default"),
					Name:          utils.GetPtr("default"),
					Steps: []*models.Step{
						{
							Type: "shell",
							Shell: &models.Shell{
								Type:          utils.GetPtr("shell"),
								Script:        utils.GetPtr("echo 'hello world'"),
								FileReference: testutils.CreateFileReference(1, 2, 3, 4),
							},
							FileReference: testutils.CreateFileReference(5, 6, 7, 8),
						},
					},
				},
			},
		},
		{
			name: "Pipeline has pull request job",
			bitbucketPipeline: &bitbucketModels.Pipeline{
				Pipelines: &bitbucketModels.BuildPipelines{
					PullRequests: &bitbucketModels.StepMap{
						"*": []*bitbucketModels.Step{
							{
								Step: &bitbucketModels.ExecutionUnitRef{
									ExecutionUnit: &bitbucketModels.ExecutionUnit{
										Script: []*bitbucketModels.Script{
											{
												String:        utils.GetPtr("echo 'hello world'"),
												FileReference: testutils.CreateFileReference(1, 2, 3, 4),
											},
											{
												String:        utils.GetPtr("echo 'hello world2'"),
												FileReference: testutils.CreateFileReference(5, 6, 7, 8),
											},
										},
									},
									FileReference: testutils.CreateFileReference(5, 6, 7, 8),
								},
							},
						},
					},
				},
			},
			expectedJobs: []*models.Job{
				{
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
					ID:            utils.GetPtr("job-*"),
					Name:          utils.GetPtr("*"),
					Steps: []*models.Step{
						{
							Type: "shell",
							Shell: &models.Shell{
								Type:          utils.GetPtr("shell"),
								Script:        utils.GetPtr("echo 'hello world'\necho 'hello world2'"),
								FileReference: testutils.CreateFileReference(1, 2, 7, 8),
							},
							FileReference: testutils.CreateFileReference(5, 6, 7, 8),
						},
					},
				},
			},
		},
		{
			name: "Pipeline has branch job",
			bitbucketPipeline: &bitbucketModels.Pipeline{
				Pipelines: &bitbucketModels.BuildPipelines{
					Branches: &bitbucketModels.StepMap{
						"master": []*bitbucketModels.Step{
							{
								Step: &bitbucketModels.ExecutionUnitRef{
									ExecutionUnit: &bitbucketModels.ExecutionUnit{
										Script: []*bitbucketModels.Script{
											{
												PipeToExecute: &bitbucketModels.PipeToExecute{
													Pipe: &bitbucketModels.Pipe{
														String:        utils.GetPtr("echo 'hello world'"),
														FileReference: testutils.CreateFileReference(3, 4, 5, 6),
													},
													Variables: &bitbucketModels.EnvironmentVariablesRef{
														EnvironmentVariables: map[string]any{
															"key": "value",
														},
														FileReference: testutils.CreateFileReference(1, 2, 3, 4),
													},
												},
												FileReference: testutils.CreateFileReference(1, 2, 3, 4),
											},
										},
									},
									FileReference: testutils.CreateFileReference(5, 6, 7, 8),
								},
							},
						},
					},
				},
			},
			expectedJobs: []*models.Job{
				{
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
					ID:            utils.GetPtr("job-master"),
					Name:          utils.GetPtr("master"),
					Steps: []*models.Step{
						{
							Type: "task",
							Task: &models.Task{
								Name:        utils.GetPtr("echo 'hello world'"),
								VersionType: "none",
							},
							EnvironmentVariables: &models.EnvironmentVariablesRef{
								EnvironmentVariables: map[string]any{
									"key": "value",
								},
								FileReference: testutils.CreateFileReference(1, 2, 3, 4),
							},
							FileReference: testutils.CreateFileReference(5, 6, 7, 8),
						},
					},
				},
			},
		},
		{
			name: "Pipeline has tag job",
			bitbucketPipeline: &bitbucketModels.Pipeline{
				Pipelines: &bitbucketModels.BuildPipelines{
					Tags: &bitbucketModels.StepMap{
						"test:1.2.3": []*bitbucketModels.Step{
							{
								Step: &bitbucketModels.ExecutionUnitRef{
									ExecutionUnit: &bitbucketModels.ExecutionUnit{},
									FileReference: testutils.CreateFileReference(5, 6, 7, 8),
								},
							},
						},
					},
				},
			},
			expectedJobs: []*models.Job{
				{
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
					ID:            utils.GetPtr("job-test:1.2.3"),
					Name:          utils.GetPtr("test:1.2.3"),
					Steps: []*models.Step{
						{
							FileReference: testutils.CreateFileReference(5, 6, 7, 8),
						},
					},
				},
			},
		},
		{
			name: "Pipeline has Custom job",
			bitbucketPipeline: &bitbucketModels.Pipeline{
				Pipelines: &bitbucketModels.BuildPipelines{
					Custom: &bitbucketModels.StepMap{
						"on-push": []*bitbucketModels.Step{
							{
								Step: &bitbucketModels.ExecutionUnitRef{
									ExecutionUnit: &bitbucketModels.ExecutionUnit{},
									FileReference: testutils.CreateFileReference(5, 6, 7, 8),
								},
							},
						},
					},
				},
			},
			expectedJobs: []*models.Job{
				{
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
					ID:            utils.GetPtr("job-on-push"),
					Name:          utils.GetPtr("on-push"),
					Steps: []*models.Step{
						{
							FileReference: testutils.CreateFileReference(5, 6, 7, 8),
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			jobs := parseJobs(testCase.bitbucketPipeline)
			testutils.DeepCompare(t, testCase.expectedJobs, jobs)
		})
	}
}

func TestStepParse(t *testing.T) {
	testCases := []struct {
		name          string
		bitbucketStep *bitbucketModels.Step
		expectedStep  []*models.Step
	}{
		{
			name:          "Step is nil",
			bitbucketStep: nil,
			expectedStep:  nil,
		},
		{
			name: "single step",
			bitbucketStep: &bitbucketModels.Step{
				Step: &bitbucketModels.ExecutionUnitRef{
					ExecutionUnit: &bitbucketModels.ExecutionUnit{
						Script: []*bitbucketModels.Script{
							{
								String:        utils.GetPtr("echo 'hello world'"),
								FileReference: testutils.CreateFileReference(1, 2, 3, 4),
							},
						},
					},
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
			},
			expectedStep: []*models.Step{
				{
					Type: "shell",
					Shell: &models.Shell{
						Type:          utils.GetPtr("shell"),
						Script:        utils.GetPtr("echo 'hello world'"),
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
			},
		},
		{
			name: "parallel steps in step",
			bitbucketStep: &bitbucketModels.Step{
				Parallel: []*bitbucketModels.ParallelSteps{
					{
						Step: &bitbucketModels.ExecutionUnitRef{
							ExecutionUnit: &bitbucketModels.ExecutionUnit{
								Script: []*bitbucketModels.Script{
									{
										String:        utils.GetPtr("echo 'hello world'"),
										FileReference: testutils.CreateFileReference(1, 2, 3, 4),
									},
								},
							},
							FileReference: testutils.CreateFileReference(5, 6, 7, 8),
						},
					},
					{
						Step: &bitbucketModels.ExecutionUnitRef{
							ExecutionUnit: &bitbucketModels.ExecutionUnit{
								Script: []*bitbucketModels.Script{
									{
										String:        utils.GetPtr("echo 'goodbye world'"),
										FileReference: testutils.CreateFileReference(4, 3, 2, 1),
									},
								},
							},
							FileReference: testutils.CreateFileReference(8, 7, 6, 5),
						},
					},
				},
			},
			expectedStep: []*models.Step{
				{
					Type: "shell",
					Shell: &models.Shell{
						Type:          utils.GetPtr("shell"),
						Script:        utils.GetPtr("echo 'hello world'"),
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
				{
					Type: "shell",
					Shell: &models.Shell{
						Type:          utils.GetPtr("shell"),
						Script:        utils.GetPtr("echo 'goodbye world'"),
						FileReference: testutils.CreateFileReference(4, 3, 2, 1),
					},
					FileReference: testutils.CreateFileReference(8, 7, 6, 5),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			step := parseStep(testCase.bitbucketStep)
			testutils.DeepCompare(t, testCase.expectedStep, step)
		})
	}
}

func TestStepRunnerParse(t *testing.T) {
	testCases := []struct {
		name             string
		bitbucketExeUnit *bitbucketModels.ExecutionUnit
		expectedRunner   *models.Runner
	}{
		{
			name:             "Execution Unit is nil",
			bitbucketExeUnit: nil,
			expectedRunner:   nil,
		},
		{
			name:             "image name is not defined",
			bitbucketExeUnit: &bitbucketModels.ExecutionUnit{},
			expectedRunner:   nil,
		},
		{
			name: "image name is defined",
			bitbucketExeUnit: &bitbucketModels.ExecutionUnit{
				Image: &bitbucketModels.Image{
					ImageData: &bitbucketModels.ImageData{
						Name: utils.GetPtr("node:10.15.3"),
					},
				},
			},
			expectedRunner: &models.Runner{
				DockerMetadata: &models.DockerMetadata{
					Image: utils.GetPtr("node:10.15.3"),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			runner := parseStepRunner(testCase.bitbucketExeUnit)
			testutils.DeepCompare(t, testCase.expectedRunner, runner)
		})
	}
}

func TestScriptParse(t *testing.T) {
	testCases := []struct {
		name             string
		bitbucketScripts []*bitbucketModels.Script
		expectedShell    *models.Shell
		expectedTask     *models.Task
	}{
		{
			name:             "Step is nil",
			bitbucketScripts: nil,
			expectedShell:    nil,
			expectedTask:     nil,
		},
		{
			name: "single script",
			bitbucketScripts: []*bitbucketModels.Script{
				{
					String:        utils.GetPtr("echo 'hello world'"),
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
			},
			expectedShell: &models.Shell{
				Type:          utils.GetPtr("shell"),
				Script:        utils.GetPtr("echo 'hello world'"),
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
		{
			name: "multiple scripts",
			bitbucketScripts: []*bitbucketModels.Script{
				{
					String:        utils.GetPtr("echo 'hello world'"),
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				{
					String:        utils.GetPtr("echo 'goodbye world'"),
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
			},
			expectedShell: &models.Shell{
				Type:          utils.GetPtr("shell"),
				Script:        utils.GetPtr("echo 'hello world'\necho 'goodbye world'"),
				FileReference: testutils.CreateFileReference(1, 2, 7, 8),
			},
		},
		{
			name: "pipe script",
			bitbucketScripts: []*bitbucketModels.Script{
				{
					PipeToExecute: &bitbucketModels.PipeToExecute{
						Pipe: &bitbucketModels.Pipe{
							String:        utils.GetPtr("echo 'hello world'"),
							FileReference: testutils.CreateFileReference(1, 2, 3, 4),
						},
					},
				},
			},
			expectedTask: &models.Task{
				Name:        utils.GetPtr("echo 'hello world'"),
				VersionType: "none",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			shell := parseScriptToShell(testCase.bitbucketScripts)
			task := parseScriptToTask(testCase.bitbucketScripts)
			testutils.DeepCompare(t, testCase.expectedShell, shell)
			testutils.DeepCompare(t, testCase.expectedTask, task)

		})
	}
}

func TestExecutionUnitParse(t *testing.T) {
	testCases := []struct {
		name             string
		bitbucketExeUnit *bitbucketModels.ExecutionUnitRef
		expectedStep     *models.Step
	}{
		{
			name:             "Step is nil",
			bitbucketExeUnit: nil,
			expectedStep:     nil,
		},
		{
			name: "execution unit has script",
			bitbucketExeUnit: &bitbucketModels.ExecutionUnitRef{
				ExecutionUnit: &bitbucketModels.ExecutionUnit{
					Name: utils.GetPtr("test"),
					Script: []*bitbucketModels.Script{
						{
							String:        utils.GetPtr("echo 'hello world'"),
							FileReference: testutils.CreateFileReference(1, 2, 3, 4),
						},
					},
				},
				FileReference: testutils.CreateFileReference(5, 6, 7, 8),
			},
			expectedStep: &models.Step{
				Name: utils.GetPtr("test"),
				Type: "shell",
				Shell: &models.Shell{
					Type:          utils.GetPtr("shell"),
					Script:        utils.GetPtr("echo 'hello world'"),
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				FileReference: testutils.CreateFileReference(5, 6, 7, 8),
			},
		},
		{
			name: "execution unit has image",
			bitbucketExeUnit: &bitbucketModels.ExecutionUnitRef{
				ExecutionUnit: &bitbucketModels.ExecutionUnit{
					Name: utils.GetPtr("test"),
					Image: &bitbucketModels.Image{
						ImageData: &bitbucketModels.ImageData{
							Name: utils.GetPtr("test"),
						},
					},
				},
				FileReference: testutils.CreateFileReference(5, 6, 7, 8),
			},
			expectedStep: &models.Step{
				Name: utils.GetPtr("test"),
				Runner: &models.Runner{
					DockerMetadata: &models.DockerMetadata{
						Image: utils.GetPtr("test"),
					},
				},
				FileReference: testutils.CreateFileReference(5, 6, 7, 8),
			},
		},
		{
			name: "execution unit has pipe",
			bitbucketExeUnit: &bitbucketModels.ExecutionUnitRef{
				ExecutionUnit: &bitbucketModels.ExecutionUnit{
					Name: utils.GetPtr("test"),
					Script: []*bitbucketModels.Script{
						{
							PipeToExecute: &bitbucketModels.PipeToExecute{
								Pipe: &bitbucketModels.Pipe{
									String:        utils.GetPtr("echo 'hello world'"),
									FileReference: testutils.CreateFileReference(1, 2, 3, 4),
								},
								Variables: &bitbucketModels.EnvironmentVariablesRef{
									EnvironmentVariables: models.EnvironmentVariables{
										"key":  "value",
										"key2": "value2",
									},
									FileReference: testutils.CreateFileReference(5, 6, 7, 8),
								},
							},
							FileReference: testutils.CreateFileReference(1, 2, 7, 8),
						},
					},
				},
				FileReference: testutils.CreateFileReference(9, 10, 11, 12),
			},
			expectedStep: &models.Step{
				Name: utils.GetPtr("test"),
				Type: "task",
				Task: &models.Task{
					Name:        utils.GetPtr("echo 'hello world'"),
					VersionType: "none",
				},
				EnvironmentVariables: &models.EnvironmentVariablesRef{
					EnvironmentVariables: models.EnvironmentVariables{
						"key":  "value",
						"key2": "value2",
					},
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
				FileReference: testutils.CreateFileReference(9, 10, 11, 12),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			step := parseExecutionUnitToStep(testCase.bitbucketExeUnit)
			testutils.DeepCompare(t, testCase.expectedStep, step)
		})
	}
}
