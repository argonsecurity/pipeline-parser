package job

import (
	"testing"

	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/job"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/r3labs/diff/v3"
	"github.com/stretchr/testify/assert"
)

func TestParseJob(t *testing.T) {
	testCases := []struct {
		name        string
		jobID       string
		job         *gitlabModels.Job
		expectedJob *models.Job
	}{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := parseJob(testCase.jobID, testCase.job)

			assert.NoError(t, err)

			changelog, err := diff.Diff(testCase.expectedJob, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0)
		})
	}
}

func TestGetJobContinueOnError(t *testing.T) {
	testCases := []struct {
		name             string
		job              *gitlabModels.Job
		expectedContinue *bool
	}{
		{
			name:             "Job is empty",
			job:              &gitlabModels.Job{},
			expectedContinue: nil,
		},
		{
			name: "Job with AllowFailure is enabled",
			job: &gitlabModels.Job{
				AllowFailure: &job.AllowFailure{
					Enabled: utils.GetPtr(true),
				},
			},
			expectedContinue: utils.GetPtr(true),
		},
		{
			name: "Job with AllowFailure is disabled",
			job: &gitlabModels.Job{
				AllowFailure: &job.AllowFailure{
					Enabled: utils.GetPtr(false),
				},
			},
			expectedContinue: utils.GetPtr(false),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := getJobContinueOnError(testCase.job)

			assert.Equal(t, testCase.expectedContinue, got, testCase.name)
		})
	}
}

func TestGetJobConcurrencyGroup(t *testing.T) {
	testCases := []struct {
		name                     string
		job                      *gitlabModels.Job
		expectedConcurrencyGroup *models.ConcurrencyGroup
	}{
		{
			name:                     "Job is empty",
			job:                      &gitlabModels.Job{},
			expectedConcurrencyGroup: nil,
		},
		{
			name: "Job with stage",
			job: &gitlabModels.Job{
				Stage: "stage",
			},
			expectedConcurrencyGroup: utils.GetPtr(models.ConcurrencyGroup("stage")),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := getJobConcurrencyGroup(testCase.job)

			assert.Equal(t, testCase.expectedConcurrencyGroup, got)
		})
	}
}

func TestGetJobConditions(t *testing.T) {
	testCases := []struct {
		name               string
		job                *gitlabModels.Job
		expectedConditions []*models.Condition
	}{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := getJobConditions(testCase.job)

			assert.ElementsMatch(t, testCase.expectedConditions, got, testCase.name)
		})
	}
}
