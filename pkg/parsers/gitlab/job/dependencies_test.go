package job

import (
	"testing"

	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/job"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestParseDependencies(t *testing.T) {
	testCases := []struct {
		name                    string
		job                     *gitlabModels.Job
		expectedJobDependencies []*models.JobDependency
	}{
		{
			name:                    "Job is empty",
			job:                     &gitlabModels.Job{},
			expectedJobDependencies: nil,
		},
		{
			name: "Job with dependencies and no needs",
			job: &gitlabModels.Job{
				Dependencies: []string{"job-1", "job-2", "job-3"},
			},
			expectedJobDependencies: []*models.JobDependency{
				{
					JobID: utils.GetPtr("job-1"),
				},
				{
					JobID: utils.GetPtr("job-2"),
				},
				{
					JobID: utils.GetPtr("job-3"),
				},
			},
		},
		{
			name: "Job with needs no dependencies",
			job: &gitlabModels.Job{
				Needs: &job.Needs{
					{
						Job:      "job-1",
						Pipeline: "pipeline-1",
					},
					{
						Job:      "job-2",
						Pipeline: "pipeline-2",
					},
				},
			},
			expectedJobDependencies: []*models.JobDependency{
				{
					JobID:    utils.GetPtr("job-1"),
					Pipeline: utils.GetPtr("pipeline-1"),
				},
				{
					JobID:    utils.GetPtr("job-2"),
					Pipeline: utils.GetPtr("pipeline-2"),
				},
			},
		},
		{
			name: "Job with dependencies and needs",
			job: &gitlabModels.Job{
				Dependencies: []string{"job-1", "job-2", "job-3"},
				Needs: &job.Needs{
					{
						Job:      "job-1",
						Pipeline: "pipeline-1",
					},
					{
						Job:      "job-2",
						Pipeline: "pipeline-2",
					},
				},
			},
			expectedJobDependencies: []*models.JobDependency{
				{
					JobID: utils.GetPtr("job-1"),
				},
				{
					JobID: utils.GetPtr("job-2"),
				},
				{
					JobID: utils.GetPtr("job-3"),
				},
				{
					JobID:    utils.GetPtr("job-1"),
					Pipeline: utils.GetPtr("pipeline-1"),
				},
				{
					JobID:    utils.GetPtr("job-2"),
					Pipeline: utils.GetPtr("pipeline-2"),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseDependencies(testCase.job)

			testutils.DeepCompare(t, testCase.expectedJobDependencies, got)
		})
	}
}

func TestParseJobDependencies(t *testing.T) {
	testCases := []struct {
		name                    string
		job                     *gitlabModels.Job
		expectedJobDependencies []*models.JobDependency
	}{
		{
			name:                    "Job is empty",
			job:                     &gitlabModels.Job{},
			expectedJobDependencies: nil,
		},
		{
			name: "Job dependencies has one element",
			job: &gitlabModels.Job{
				Dependencies: []string{"job-1"},
			},
			expectedJobDependencies: []*models.JobDependency{
				{
					JobID: utils.GetPtr("job-1"),
				},
			},
		},
		{
			name: "Job dependencies has some elements",
			job: &gitlabModels.Job{
				Dependencies: []string{"job-1", "job-2", "job-3"},
			},
			expectedJobDependencies: []*models.JobDependency{
				{
					JobID: utils.GetPtr("job-1"),
				},
				{
					JobID: utils.GetPtr("job-2"),
				},
				{
					JobID: utils.GetPtr("job-3"),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseJobDependencies(testCase.job)

			testutils.DeepCompare(t, testCase.expectedJobDependencies, got)
		})
	}
}

func TestParseJobNeeds(t *testing.T) {
	testCases := []struct {
		name                    string
		job                     *gitlabModels.Job
		expectedJobDependencies []*models.JobDependency
	}{
		{
			name:                    "Job is empty",
			job:                     &gitlabModels.Job{},
			expectedJobDependencies: nil,
		},
		{
			name: "Job needs has one element with job and pipeline",
			job: &gitlabModels.Job{
				Needs: &job.Needs{
					{
						Job:      "job-1",
						Pipeline: "pipeline-1",
					},
				},
			},
			expectedJobDependencies: []*models.JobDependency{
				{
					JobID:    utils.GetPtr("job-1"),
					Pipeline: utils.GetPtr("pipeline-1"),
				},
			},
		},
		{
			name: "Job needs has one element with job and no pipeline",
			job: &gitlabModels.Job{
				Needs: &job.Needs{
					{
						Job:      "job-1",
						Pipeline: "",
					},
				},
			},
			expectedJobDependencies: []*models.JobDependency{
				{
					JobID:    utils.GetPtr("job-1"),
					Pipeline: nil,
				},
			},
		},
		{
			name: "Job needs has one element with pipeline and no job",
			job: &gitlabModels.Job{
				Needs: &job.Needs{
					{
						Job:      "",
						Pipeline: "pipeline-1",
					},
				},
			},
			expectedJobDependencies: []*models.JobDependency{
				{
					JobID:    nil,
					Pipeline: utils.GetPtr("pipeline-1"),
				},
			},
		},
		{
			name: "Job needs has one element with pipeline and no job",
			job: &gitlabModels.Job{
				Needs: &job.Needs{
					{
						Job:      "",
						Pipeline: "pipeline-1",
					},
					{
						Job:      "job-1",
						Pipeline: "pipeline-1",
					},
				},
			},
			expectedJobDependencies: []*models.JobDependency{
				{
					JobID:    nil,
					Pipeline: utils.GetPtr("pipeline-1"),
				},
				{
					JobID:    utils.GetPtr("job-1"),
					Pipeline: utils.GetPtr("pipeline-1"),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseJobNeeds(testCase.job)

			testutils.DeepCompare(t, testCase.expectedJobDependencies, got)
		})
	}
}
