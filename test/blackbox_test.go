package test

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/parsers/github"
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

func getAllGitHubPermissions(permission models.Permission) *map[string]models.Permission {
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
	return &allPermissions
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
							},
							{
								Name: utils.GetPtr("task with commit ID version"),
								Type: "task",
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("c44948622e1b6bb0eb0cec5b813c1ac561158e1e"),
									VersionType: "commit",
								},
							},
							{
								Name: utils.GetPtr("task with branch version"),
								Type: "task",
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("master"),
									VersionType: "branch",
								},
							},
							{
								Name: utils.GetPtr("task with tag version"),
								Type: "task",
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("v1.1.1"),
									VersionType: "tag",
								},
							},
							{
								Name: utils.GetPtr("shell"),
								Type: "shell",
								Shell: &models.Shell{
									Script: utils.GetPtr("command line"),
								},
							},
							{
								Name: utils.GetPtr("custom shell"),
								Type: "shell",
								Shell: &models.Shell{
									Script: utils.GetPtr("command line"),
									Type:   utils.GetPtr("cmd"),
								},
							},
						},
						TimeoutMS:       utils.GetPtr(21600000),
						ContinueOnError: utils.GetPtr(false),
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
						ID:              utils.GetPtr("dependant-job"),
						Name:            utils.GetPtr("Dependant Job"),
						ContinueOnError: utils.GetPtr(false),
						TimeoutMS:       utils.GetPtr(21600000),
						Dependencies:    &[]string{"dependable-job"},
					},
					{
						ID:              utils.GetPtr("dependable-job"),
						Name:            utils.GetPtr("Dependable Job"),
						ContinueOnError: utils.GetPtr(false),
						TimeoutMS:       utils.GetPtr(21600000),
					},
				}),
			},
		},
		{
			Filename: "all-triggers.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("all-triggers"),
				Triggers: SortTriggers(&[]models.Trigger{
					{
						Event:     models.ScheduledEvent,
						Scheduled: utils.GetPtr("30 2 * * *"),
					},
					{
						Event: models.PushEvent,
						Branches: &models.Filter{
							AllowList: []string{"master"},
						},
					},
					{
						Event: models.PipelineRunEvent,
						Branches: &models.Filter{
							DenyList: []string{"master"},
						},
					},
					{
						Event: models.PullRequestEvent,
						Paths: &models.Filter{
							DenyList: []string{"*/test/*"},
						},
					},
					{
						Event: "pull_request_target",
						Paths: &models.Filter{
							AllowList: []string{"*/test/*"},
						},
					},
					{
						Event: models.ManualEvent,
						Paramters: []models.Parameter{
							{
								Name:        utils.GetPtr("workflow-input"),
								Description: utils.GetPtr("The workflow input"),
								Default:     "default-value",
							},
						},
					},
					{
						Event: models.PipelineTriggerEvent,
						Paramters: []models.Parameter{
							{
								Name:        utils.GetPtr("workflow-input"),
								Description: utils.GetPtr("The workflow input"),
								Default:     "default-value",
							},
						},
					},
				}),
			},
		},
		{
			Filename: "token-permissions.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("permissions"),
				Jobs: SortJobs(&[]models.Job{
					{
						ID:               utils.GetPtr("job1"),
						Name:             utils.GetPtr("Job 1"),
						ContinueOnError:  utils.GetPtr(false),
						TokenPermissions: getAllGitHubPermissions(models.Permission{Read: true}),
						TimeoutMS:        utils.GetPtr(21600000),
					},
				}),
				Defaults: &models.Defaults{
					TokenPermissions: &map[string]models.Permission{

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
	}

	for _, testCase := range testCases {
		pipeline, err := github.Parse(readFile(filepath.Join(workflowFolderPath, testCase.Filename)))
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
			pipeline.Triggers = SortTriggers(pipeline.Triggers)
		}

		if diff := deep.Equal(pipeline, testCase.Expected); diff != nil {
			t.Errorf("%s: %s", testCase.Filename, diff)
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
