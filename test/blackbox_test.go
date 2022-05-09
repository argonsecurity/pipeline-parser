package test

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/handler"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"

	"github.com/go-test/deep"
)

const (
	workflowFolderPath = "fixtures/github/"
)

type TestCase struct {
	Filename   string
	Expected   *models.Pipeline
	ShouldFail bool
}

func readFile(filename string) []byte {
	b, _ := ioutil.ReadFile(filename)
	return b
}

func getAllGitHubPermissions(permission models.Permission) *models.TokenPermissions {
	allPermissions := map[string]models.Permission{
		"run-pipeline":        permission,
		"checks":              permission,
		"contents":            permission,
		"deployments":         permission,
		"discussions":         permission,
		"id-token":            permission,
		"issues":              permission,
		"packages":            permission,
		"pages":               permission,
		"pull-request":        permission,
		"repository-projects": permission,
		"security-events":     permission,
		"statuses":            permission,
	}
	return &models.TokenPermissions{
		Permissions: allPermissions,
	}
}

func createFileLocation(l1, c1, l2, c2 int) *models.FileLocation {
	return &models.FileLocation{
		StartRef: &models.FileRef{
			Line:   l1,
			Column: c1,
		},
		EndRef: &models.FileRef{
			Line:   l2,
			Column: c2,
		},
	}
}

func Test_GitHubParser(t *testing.T) {
	testCases := []TestCase{
		{
			Filename: "steps.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("steps"),
				Jobs: SortJobs(&[]models.Job{
					{
						ID:   utils.GetPtr("job1"),
						Name: utils.GetPtr("Job 1"),
						Steps: &[]models.Step{
							{
								Name: utils.GetPtr("task without params"),
								Type: "task",
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("v1"),
									VersionType: "tag",
								},
								FileLocation: createFileLocation(7, 9, 8, 15),
							},
							{
								Name: utils.GetPtr("task with params"),
								Type: "task",
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("v1"),
									VersionType: "tag",
									Inputs: &[]models.Parameter{
										{
											Name:  utils.GetPtr("repo"),
											Value: "repository",
										},
									},
								},
								FileLocation: createFileLocation(10, 9, 13, 17),
							},
							{
								Name: utils.GetPtr("task with commit ID version"),
								Type: "task",
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("c44948622e1b6bb0eb0cec5b813c1ac561158e1e"),
									VersionType: "commit",
								},
								FileLocation: createFileLocation(15, 9, 16, 15),
							},
							{
								Name: utils.GetPtr("task with branch version"),
								Type: "task",
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("master"),
									VersionType: "branch",
								},
								FileLocation: createFileLocation(18, 9, 19, 15),
							},
							{
								Name: utils.GetPtr("task with tag version"),
								Type: "task",
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("v1.1.1"),
									VersionType: "tag",
								},
								FileLocation: createFileLocation(21, 9, 22, 15),
							},
							{
								Name: utils.GetPtr("shell"),
								Type: "shell",
								Shell: &models.Shell{
									Script: utils.GetPtr("command line"),
								},
								FileLocation: createFileLocation(24, 9, 25, 14),
							},
							{
								Name: utils.GetPtr("custom shell"),
								Type: "shell",
								Shell: &models.Shell{
									Script: utils.GetPtr("command line"),
									Type:   utils.GetPtr("cmd"),
								},
								FileLocation: createFileLocation(27, 9, 29, 14),
							},
						},
						TimeoutMS:       utils.GetPtr(21600000),
						ContinueOnError: utils.GetPtr(false),
						FileLocation:    createFileLocation(4, 3, 29, 14),
					},
				}),
			},
		},
		{
			Filename: "dependant-jobs.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("dependable jobs"),
				Jobs: SortJobs(&[]models.Job{
					{
						ID:              utils.GetPtr("dependable-job"),
						Name:            utils.GetPtr("Dependable Job"),
						ContinueOnError: utils.GetPtr(false),
						TimeoutMS:       utils.GetPtr(21600000),
						FileLocation:    createFileLocation(4, 3, 5, 11),
					},
					{
						ID:              utils.GetPtr("dependant-job"),
						Name:            utils.GetPtr("Dependant Job"),
						ContinueOnError: utils.GetPtr(false),
						TimeoutMS:       utils.GetPtr(21600000),
						Dependencies:    &[]string{"dependable-job"},
						FileLocation:    createFileLocation(7, 3, 9, 13),
					},
				}),
			},
		},
		{
			Filename: "all-triggers.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("all-triggers"),
				Triggers: &models.Triggers{
					FileLocation: createFileLocation(2, 3, 30, 13),
					Triggers: SortTriggers(&[]models.Trigger{
						{
							Event:        models.ScheduledEvent,
							Schedules:    &[]string{"30 2 * * *"},
							FileLocation: createFileLocation(3, 3, 4, 13),
						},
						{
							Event: models.PushEvent,
							Branches: &models.Filter{
								AllowList: []string{"master"},
							},
							FileLocation: createFileLocation(5, 3, 7, 9),
						},
						{
							Event: models.PipelineRunEvent,
							Branches: &models.Filter{
								DenyList: []string{"master"},
							},
							FileLocation: createFileLocation(26, 3, 28, 9),
						},
						{
							Event: models.PullRequestEvent,
							Paths: &models.Filter{
								DenyList: []string{"*/test/*"},
							},
							FileLocation: createFileLocation(8, 3, 10, 9),
						},
						{
							Event: "pull_request_target",
							Paths: &models.Filter{
								AllowList: []string{"*/test/*"},
							},
							FileLocation: createFileLocation(11, 3, 13, 9),
						},
						{
							Event: models.ManualEvent,
							Parameters: []models.Parameter{
								{
									Name:        utils.GetPtr("workflow-input"),
									Description: utils.GetPtr("The workflow input"),
									Default:     "default-value",
								},
							},
							FileLocation: createFileLocation(14, 3, 19, 19),
						},
						{
							Event: models.PipelineTriggerEvent,
							Parameters: []models.Parameter{
								{
									Name:        utils.GetPtr("workflow-input"),
									Description: utils.GetPtr("The workflow input"),
									Default:     "default-value",
								},
							},
							FileLocation: createFileLocation(20, 3, 25, 19),
						},
						{
							Event: "label",
							Filters: map[string]any{
								"types": []string{"created"},
							},
							FileLocation: createFileLocation(29, 3, 30, 13),
						},
					}),
				},
			},
		},
		{
			Filename: "token-permissions.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("permissions"),
				Jobs: SortJobs(&[]models.Job{
					{
						FileLocation:     createFileLocation(8, 3, 10, 18),
						ID:               utils.GetPtr("job1"),
						Name:             utils.GetPtr("Job 1"),
						ContinueOnError:  utils.GetPtr(false),
						TokenPermissions: getAllGitHubPermissions(models.Permission{Read: true}),
						TimeoutMS:        utils.GetPtr(21600000),
					},
				}),
				Defaults: &models.Defaults{
					TokenPermissions: &models.TokenPermissions{
						FileLocation: createFileLocation(2, 3, 5, 18),
						Permissions: map[string]models.Permission{
							"run-pipeline": {
								Read: true,
							},
							"statuses": {
								Write: true,
							},
							"pull-request": {
								Read: true,
							},
						},
					},
				},
			},
		},
		{
			Filename: "runners.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("runners"),
				Jobs: SortJobs(&[]models.Job{
					{
						ID:              utils.GetPtr("job1"),
						Name:            utils.GetPtr("Job 1"),
						ContinueOnError: utils.GetPtr(false),
						TimeoutMS:       utils.GetPtr(21600000),
						Runner: &models.Runner{
							OS:           utils.GetPtr("linux"),
							Labels:       &[]string{"ubuntu-latest"},
							SelfHosted:   utils.GetPtr(false),
							FileLocation: createFileLocation(6, 14, 6, 14),
						},
						FileLocation: createFileLocation(4, 3, 6, 14),
					},
					{
						ID:              utils.GetPtr("job2"),
						Name:            utils.GetPtr("Job 2"),
						TimeoutMS:       utils.GetPtr(21600000),
						ContinueOnError: utils.GetPtr(false),
						Runner: &models.Runner{
							OS:           utils.GetPtr("windows"),
							Labels:       &[]string{"self-hosted", "windows-latest"},
							SelfHosted:   utils.GetPtr(true),
							FileLocation: createFileLocation(9, 14, 9, 28),
						},
						FileLocation: createFileLocation(7, 3, 9, 28),
					},
					{
						ID:              utils.GetPtr("job3"),
						Name:            utils.GetPtr("Job 3"),
						TimeoutMS:       utils.GetPtr(21600000),
						ContinueOnError: utils.GetPtr(false),
						Runner: &models.Runner{
							OS:           utils.GetPtr("linux"),
							Arch:         utils.GetPtr("x64"),
							Labels:       &[]string{"self-hosted", "linux", "x64"},
							SelfHosted:   utils.GetPtr(true),
							FileLocation: createFileLocation(12, 14, 12, 35),
						},
						FileLocation: createFileLocation(10, 3, 12, 35),
					},
				}),
			},
		},
		{
			Filename: "environment-variables.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("environment-variables"),
				Jobs: SortJobs(&[]models.Job{
					{
						ID:   utils.GetPtr("job1"),
						Name: utils.GetPtr("Job 1"),
						EnvironmentVariables: &models.EnvironmentVariablesRef{
							EnvironmentVariables: models.EnvironmentVariables{
								"STRING": "string",
								"NUMBER": 1,
							},
							FileLocation: createFileLocation(10, 7, 12, 15),
						},
						FileLocation:    createFileLocation(8, 3, 18, 19),
						ContinueOnError: utils.GetPtr(false),
						TimeoutMS:       utils.GetPtr(21600000),
						Steps: &[]models.Step{
							{
								Name: utils.GetPtr("Step 1"),
								Type: "shell",
								Shell: &models.Shell{
									Script: utils.GetPtr("command line"),
								},
								FileLocation: createFileLocation(14, 9, 18, 19),
								EnvironmentVariables: &models.EnvironmentVariablesRef{
									EnvironmentVariables: models.EnvironmentVariables{
										"STRING": "string",
										"NUMBER": 1,
									},
									FileLocation: createFileLocation(16, 11, 18, 19),
								},
							},
						},
					},
				}),
				Defaults: &models.Defaults{
					EnvironmentVariables: &models.EnvironmentVariablesRef{
						EnvironmentVariables: models.EnvironmentVariables{
							"STRING": "string",
							"NUMBER": 1,
						},
						FileLocation: createFileLocation(3, 3, 5, 11),
					},
				},
			},
		},
		{
			Filename: "concurrent-jobs.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("concurrent-jobs"),
				Jobs: SortJobs(&[]models.Job{
					{
						ID:               utils.GetPtr("job1"),
						Name:             utils.GetPtr("Job 1"),
						ContinueOnError:  utils.GetPtr(false),
						TimeoutMS:        utils.GetPtr(21600000),
						ConcurrencyGroup: utils.GetPtr("ci"),
						FileLocation:     createFileLocation(3, 3, 5, 18),
					},
					{
						ID:               utils.GetPtr("job2"),
						Name:             utils.GetPtr("Job 2"),
						ContinueOnError:  utils.GetPtr(false),
						TimeoutMS:        utils.GetPtr(21600000),
						ConcurrencyGroup: utils.GetPtr("ci"),
						FileLocation:     createFileLocation(7, 3, 9, 18),
					},
				}),
			},
		},
	}

	for _, testCase := range testCases {
		buf := readFile(filepath.Join(workflowFolderPath, testCase.Filename))
		pipeline, err := handler.Handle(buf, consts.GitHubPlatform)
		if err != nil {
			if !testCase.ShouldFail {
				t.Errorf("%s: %s", testCase.Filename, err)
			}
			continue
		}

		if testCase.ShouldFail {
			t.Errorf("%s: expected error, but got none", testCase.Filename)
			continue
		}

		if pipeline.Jobs != nil {
			pipeline.Jobs = SortJobs(pipeline.Jobs)
		}
		if pipeline.Triggers != nil {
			pipeline.Triggers = &models.Triggers{
				Triggers:     SortTriggers(pipeline.Triggers.Triggers),
				FileLocation: pipeline.Triggers.FileLocation,
			}
		}

		if diffs := deep.Equal(pipeline, testCase.Expected); diffs != nil {
			t.Errorf("there are %d differences in file %s", len(diffs), testCase.Filename)
			for _, diff := range diffs {
				t.Errorf(diff)
			}
		}
	}
}

func SortJobs(jobs *[]models.Job) *[]models.Job {
	sort.Slice(*jobs, func(i, j int) bool {
		return *(*jobs)[i].ID < *(*jobs)[j].ID
	})
	return jobs
}

func SortTriggers(triggers *[]models.Trigger) *[]models.Trigger {
	sort.Slice(*triggers, func(i, j int) bool {
		return (*triggers)[i].Event < (*triggers)[j].Event
	})
	return triggers
}

func SortParameters(parameters *[]models.Parameter) *[]models.Parameter {
	sort.Slice(*parameters, func(i, j int) bool {
		return *(*parameters)[i].Name < *(*parameters)[j].Name
	})
	return parameters
}
