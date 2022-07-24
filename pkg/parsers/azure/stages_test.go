package azure

import (
	"testing"

	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/r3labs/diff/v3"
	"github.com/stretchr/testify/assert"
)

func TestParseStages(t *testing.T) {
	testCases := []struct {
		name         string
		stages       *azureModels.Stages
		expectedJobs []*models.Job
	}{
		{
			name:         "Stages is nil",
			stages:       nil,
			expectedJobs: nil,
		},
		{
			name:         "Stages is empty",
			stages:       &azureModels.Stages{},
			expectedJobs: nil,
		},
		{
			name: "Stages with data",
			stages: &azureModels.Stages{
				Stages: []*azureModels.Stage{
					{
						Jobs: &azureModels.Jobs{
							CIJobs: []*azureModels.CIJob{
								{
									Job:           "ci-job-1",
									FileReference: testutils.CreateFileReference(1, 2, 3, 4),
								},
							},
						},
					},
					{
						Jobs: &azureModels.Jobs{
							CIJobs: []*azureModels.CIJob{
								{
									Job:           "ci-job-2",
									FileReference: testutils.CreateFileReference(5, 6, 7, 8),
								},
							},
						},
					},
				},
			},
			expectedJobs: []*models.Job{
				{
					ID:              utils.GetPtr("ci-job-1"),
					Name:            utils.GetPtr(""),
					TimeoutMS:       utils.GetPtr(defaultTimeoutMS),
					ContinueOnError: utils.GetPtr(false),
					FileReference:   testutils.CreateFileReference(1, 2, 3, 4),
				},
				{
					ID:              utils.GetPtr("ci-job-2"),
					Name:            utils.GetPtr(""),
					TimeoutMS:       utils.GetPtr(defaultTimeoutMS),
					ContinueOnError: utils.GetPtr(false),
					FileReference:   testutils.CreateFileReference(5, 6, 7, 8),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseStages(testCase.stages)

			changelog, err := diff.Diff(testCase.expectedJobs, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}
}

func TestParseStage(t *testing.T) {
	testCases := []struct {
		name         string
		stage        *azureModels.Stage
		expectedJobs []*models.Job
	}{
		{
			name:         "Stage is nil",
			stage:        nil,
			expectedJobs: nil,
		},
		{
			name:         "Stage is empty",
			stage:        &azureModels.Stage{},
			expectedJobs: nil,
		},
		{
			name: "Stage has jobs",
			stage: &azureModels.Stage{
				Jobs: &azureModels.Jobs{
					CIJobs: []*azureModels.CIJob{
						{
							Job:           "ci-job",
							FileReference: testutils.CreateFileReference(1, 2, 3, 4),
						},
					},
					DeploymentJobs: []*azureModels.DeploymentJob{
						{
							Deployment:    "deployment-job",
							FileReference: testutils.CreateFileReference(1, 2, 3, 4),
						},
					},
				},
			},
			expectedJobs: []*models.Job{
				{
					ID:              utils.GetPtr("ci-job"),
					Name:            utils.GetPtr(""),
					TimeoutMS:       utils.GetPtr(defaultTimeoutMS),
					ContinueOnError: utils.GetPtr(false),
					FileReference:   testutils.CreateFileReference(1, 2, 3, 4),
				},
				{
					ID:              utils.GetPtr("deployment-job"),
					Name:            utils.GetPtr(""),
					TimeoutMS:       utils.GetPtr(defaultTimeoutMS),
					ContinueOnError: utils.GetPtr(false),
					FileReference:   testutils.CreateFileReference(1, 2, 3, 4),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseStage(testCase.stage)

			changelog, err := diff.Diff(testCase.expectedJobs, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}
}
