package blackbox

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/handler"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/go-test/deep"
)

func readFile(filename string) []byte {
	b, _ := ioutil.ReadFile(filename)
	return b
}

func executeTestCases(t *testing.T, testCases []TestCase, folder string, platform consts.Platform) {
	for _, testCase := range testCases {
		buf := readFile(filepath.Join("../fixtures", folder, testCase.Filename))
		pipeline, err := handler.Handle(buf, platform)
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
				Triggers:      SortTriggers(pipeline.Triggers.Triggers),
				FileReference: pipeline.Triggers.FileReference,
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

func SortJobs(jobs []*models.Job) []*models.Job {
	sort.Slice(jobs, func(i, j int) bool {
		return *jobs[i].ID < *jobs[j].ID
	})
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
