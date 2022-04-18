package test

import (
	"io/ioutil"
	"path/filepath"
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

		if diff := deep.Equal(pipeline, testCase.Expected); diff != nil {
			t.Errorf("%s: %s", testCase.Filename, diff)
		}
	}
}
