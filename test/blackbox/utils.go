package blackbox

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"testing"

	githubEnhancer "github.com/argonsecurity/pipeline-parser/pkg/enhancers/github"
	gitlabEnhancer "github.com/argonsecurity/pipeline-parser/pkg/enhancers/gitlab"
	"github.com/argonsecurity/pipeline-parser/pkg/handler"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/go-test/deep"
)

func readFile(filename string) []byte {
	b, _ := os.ReadFile(filename)
	return b
}

func executeTestCases(t *testing.T, testCases []TestCase, folder string, platform models.Platform) {
	for _, testCase := range testCases {
		if testCase.TestdataDir != "" {
			h := http.FileServer(http.Dir(testCase.TestdataDir))
			ts := httptest.NewServer(h)
			githubEnhancer.GithubBaseURL = ts.URL
			gitlabEnhancer.GITLAB_BASE_URL = ts.URL
		}

		buf := readFile(filepath.Join("../fixtures", folder, testCase.Filename))
		pipeline, err := handler.Handle(buf, platform, &models.Credentials{})
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

		pipeline = SortPipeline(pipeline)

		if diffs := deep.Equal(pipeline, testCase.Expected); diffs != nil {
			t.Errorf("there are %d differences in file %s", len(diffs), testCase.Filename)
			for _, diff := range diffs {
				t.Errorf(diff)
			}
		}
	}
}

func SortPipeline(pipeline *models.Pipeline) *models.Pipeline {
	if pipeline.Jobs != nil {
		pipeline.Jobs = SortJobs(pipeline.Jobs)
	}
	if pipeline.Triggers != nil {
		pipeline.Triggers = &models.Triggers{
			Triggers:      SortTriggers(pipeline.Triggers.Triggers),
			FileReference: pipeline.Triggers.FileReference,
		}
	}

	for _, importData := range pipeline.Imports {
		if importData.Pipeline != nil {
			importData.Pipeline = SortPipeline(importData.Pipeline)
		}
	}
	return pipeline
}

func SortJobs(jobs []*models.Job) []*models.Job {
	sort.Slice(jobs, func(i, j int) bool {
		return *jobs[i].ID < *jobs[j].ID
	})

	for _, job := range jobs {
		if job.Imports != nil && job.Imports.Pipeline != nil && job.Imports.Pipeline.Jobs != nil {
			job.Imports.Pipeline.Jobs = SortJobs(job.Imports.Pipeline.Jobs)
		}
	}
	return jobs
}

func SortTriggers(triggers []*models.Trigger) []*models.Trigger {
	sort.Slice(triggers, func(i, j int) bool {
		return triggers[i].Event < triggers[j].Event
	})
	return triggers
}

func SortParameters(parameters *[]models.Parameter) *[]models.Parameter {
	sort.Slice(*parameters, func(i, j int) bool {
		return *(*parameters)[i].Name < *(*parameters)[j].Name
	})
	return parameters
}
