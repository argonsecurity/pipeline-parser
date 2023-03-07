package bitbucket

import (
	"testing"

	bbModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	bitbucketModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestJobsParse(t *testing.T) {
	testCases := []struct {
		name              string
		bitbucketPipeline *bbModels.Pipeline
		expectedJobs      []*models.Job
	}{
		{
			name:              "Pipeline is nil",
			bitbucketPipeline: nil,
			expectedJobs:      nil,
		},
		{
			name: "Pipeline has no jobs",
			bitbucketPipeline: &bbModels.Pipeline{
				Image: &bbModels.Image{
					ImageData: &bbModels.ImageData{
						Name: utils.GetPtr("node:10.15.3"),
					},
				},
			},
			expectedJobs: nil,
		},
		{
			name: "Pipeline has default job",
			bitbucketPipeline: &bbModels.Pipeline{
				Pipelines: &bbModels.BuildPipelines{
					Default: []*bbModels.Step{
						{
							Step: &bbModels.ExecutionUnitRef{
								ExecutionUnit: &bbModels.ExecutionUnit{
									Script: []*bbModels.Script{
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
			bitbucketPipeline: &bbModels.Pipeline{
				Pipelines: &bbModels.BuildPipelines{
					PullRequests: &bbModels.StepMap{
						"*": []*bbModels.Step{
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Script: []*bbModels.Script{
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
			bitbucketPipeline: &bbModels.Pipeline{
				Pipelines: &bbModels.BuildPipelines{
					Branches: &bbModels.StepMap{
						"master": []*bbModels.Step{
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Script: []*bbModels.Script{
											{
												PipeToExecute: &bbModels.PipeToExecute{
													Pipe: &bbModels.Pipe{
														String:        utils.GetPtr("echo 'hello world'"),
														FileReference: testutils.CreateFileReference(3, 4, 5, 6),
													},
													Variables: &bbModels.EnvironmentVariablesRef{
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
			bitbucketPipeline: &bbModels.Pipeline{
				Pipelines: &bbModels.BuildPipelines{
					Tags: &bbModels.StepMap{
						"test:1.2.3": []*bbModels.Step{
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{},
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
			bitbucketPipeline: &bbModels.Pipeline{
				Pipelines: &bbModels.BuildPipelines{
					Custom: &bbModels.StepMap{
						"on-push": []*bbModels.Step{
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{},
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
			bitbucketStep: &bbModels.Step{
				Step: &bbModels.ExecutionUnitRef{
					ExecutionUnit: &bbModels.ExecutionUnit{
						Script: []*bbModels.Script{
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
			bitbucketStep: &bbModels.Step{
				Parallel: []*bbModels.ParallelSteps{
					{
						Step: &bbModels.ExecutionUnitRef{
							ExecutionUnit: &bbModels.ExecutionUnit{
								Script: []*bbModels.Script{
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
						Step: &bbModels.ExecutionUnitRef{
							ExecutionUnit: &bbModels.ExecutionUnit{
								Script: []*bbModels.Script{
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
		bitbucketExeUnit *bbModels.ExecutionUnit
		expectedRunner   *models.Runner
	}{
		{
			name:              "Execution Unit is nil",
			bitbucketExeUnit: nil,
			expectedRunner:    nil,
		},
		{
			name:              "image name is not defined",
			bitbucketExeUnit: &bbModels.ExecutionUnit{},
			expectedRunner:    nil,
		},
		{
			name: "image name is defined",
			bitbucketExeUnit: &bbModels.ExecutionUnit{
				Image: &bbModels.Image{
					ImageData: &bbModels.ImageData{
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
			bitbucketScripts: []*bbModels.Script{
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
			bitbucketScripts: []*bbModels.Script{
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
			bitbucketScripts: []*bbModels.Script{
				{
					PipeToExecute: &bbModels.PipeToExecute{
						Pipe: &bbModels.Pipe{
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
			bitbucketExeUnit: &bbModels.ExecutionUnitRef{
				ExecutionUnit: &bbModels.ExecutionUnit{
					Name: utils.GetPtr("test"),
					Script: []*bbModels.Script{
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
			bitbucketExeUnit: &bbModels.ExecutionUnitRef{
				ExecutionUnit: &bbModels.ExecutionUnit{
					Name: utils.GetPtr("test"),
					Image: &bbModels.Image{
						ImageData: &bbModels.ImageData{
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
			bitbucketExeUnit: &bbModels.ExecutionUnitRef{
				ExecutionUnit: &bbModels.ExecutionUnit{
					Name: utils.GetPtr("test"),
					Script: []*bbModels.Script{
						{
							PipeToExecute: &bbModels.PipeToExecute{
								Pipe: &bbModels.Pipe{
									String:        utils.GetPtr("echo 'hello world'"),
									FileReference: testutils.CreateFileReference(1, 2, 3, 4),
								},
								Variables: &bbModels.EnvironmentVariablesRef{
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
