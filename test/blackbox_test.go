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
		// {
		// 	Filename: "plain-scheduled.yaml",
		// 	Expected: &models.Pipeline{
		// 		Name: utils.GetPtr("plain scheduled"),
		// 		Triggers: &[]models.Trigger{
		// 			{
		// 				Event:     models.ScheduledEvent,
		// 				Scheduled: utils.GetPtr("30 2 * * *"),
		// 			},
		// 		},
		// 		Jobs: utils.GetPtr([]models.Job{
		// 			{
		// 				ID:   utils.GetPtr("GenerateSourcecred"),
		// 				Name: utils.GetPtr("GenerateSourcecred"),
		// 				Runner: &models.Runner{
		// 					OS:     utils.GetPtr("linux"),
		// 					Labels: &[]string{"ubuntu-latest"},
		// 				},
		// 				Steps: &[]models.Step{
		// 					{
		// 						Name: utils.GetPtr("Checkout Repository"),
		// 						Type: "task",
		// 						Task: &models.Task{
		// 							Name:        utils.GetPtr("actions/checkout"),
		// 							Version:     utils.GetPtr("v1"),
		// 							VersionType: "tag",
		// 						},
		// 					},
		// 				},
		// 				TimeoutMS:       utils.GetPtr(21600000),
		// 				ContinueOnError: utils.GetPtr(false),
		// 			},
		// 		}),
		// 	},
		// },
		{
			Filename: "dependant-jobs.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("manual"),
				Triggers: SortTriggers(&[]models.Trigger{
					{
						Event: models.ManualEvent,
					},
					{
						Event: models.PullRequestEvent,
						Paths: &models.Filter{
							DenyList: []string{"**.md"},
						},
					},
					{
						Event: models.PushEvent,
						Branches: &models.Filter{
							AllowList: []string{"develop"},
						},
					},
				}),
				Jobs: SortJobs(&[]models.Job{
					{
						ID:              utils.GetPtr("release-ubuntu"),
						Name:            utils.GetPtr("release-ubuntu"),
						ContinueOnError: utils.GetPtr(false),
						Runner: &models.Runner{
							OS:     utils.GetPtr("linux"),
							Labels: &[]string{"ubuntu-latest"},
						},
						Steps: &[]models.Step{
							{
								Name: utils.GetPtr("Checkout"),
								Type: "task",
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("v2"),
									VersionType: "tag",
									Inputs: &[]models.Parameter{
										{Name: utils.GetPtr("submodules"), Value: true},
									},
								},
							},
						},
						TimeoutMS:    utils.GetPtr(21600000),
						Dependencies: &[]string{"build-core"},
					},
					{
						ID:              utils.GetPtr("build-core"),
						Name:            utils.GetPtr("build-core"),
						ContinueOnError: utils.GetPtr(false),
						Runner: &models.Runner{
							OS:     utils.GetPtr("linux"),
							Labels: &[]string{"ubuntu-latest"},
						},
						Steps: &[]models.Step{
							{
								Name: utils.GetPtr("Upload artifacts"),
								Type: "task",
								Task: &models.Task{
									Name:        utils.GetPtr("actions/upload-artifact"),
									Version:     utils.GetPtr("v1"),
									VersionType: "tag",
									Inputs: SortParameters(&[]models.Parameter{
										{Name: utils.GetPtr("path"), Value: "ocaml-build-artifacts.tgz"},
										{Name: utils.GetPtr("name"), Value: "ocaml-build-artifacts"},
									}),
								},
							},
						},
						TimeoutMS: utils.GetPtr(21600000),
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

		pipeline.Jobs = SortJobs(pipeline.Jobs)
		pipeline.Triggers = SortTriggers(pipeline.Triggers)

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
