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

func Test_GitHubParser(t *testing.T) {
	testCases := []TestCase{
		{
			Filename: "plain-scheduled.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("plain scheduled"),
				Triggers: &[]models.Trigger{
					{
						Event:     models.ScheduledEvent,
						Scheduled: utils.GetPtr("30 2 * * *"),
					},
				},
				Jobs: utils.GetPtr([]models.Job{
					{
						ID:   utils.GetPtr("GenerateSourcecred"),
						Name: utils.GetPtr("GenerateSourcecred"),
						Runner: &models.Runner{
							OS:     utils.GetPtr("linux"),
							Labels: &[]string{"ubuntu-latest"},
						},
						Steps: &[]models.Step{
							{
								Name: utils.GetPtr("Checkout Repository"),
								Type: "task",
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("v1"),
									VersionType: "tag",
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
