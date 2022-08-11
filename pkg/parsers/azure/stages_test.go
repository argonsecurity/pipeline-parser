package azure

import (
	"testing"

	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
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

			testutils.DeepCompare(t, testCase.expectedJobs, got)
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
				Variables: &azureModels.Variables{
					{
						Name:          "var1",
						Value:         "value1",
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
					{
						Group:         "group1",
						FileReference: testutils.CreateFileReference(5, 6, 7, 8),
					},
					{
						Name:          "var2",
						Value:         "value2",
						FileReference: testutils.CreateFileReference(9, 10, 11, 12),
					},
				},
				Jobs: &azureModels.Jobs{
					CIJobs: []*azureModels.CIJob{
						{
							Job: "ci-job",
							BaseJob: azureModels.BaseJob{
								Variables: &azureModels.Variables{
									{
										Name:          "varjob1",
										Value:         "valuejob1",
										FileReference: testutils.CreateFileReference(31, 23, 33, 34),
									},
									{
										Group:         "groupjob1",
										FileReference: testutils.CreateFileReference(35, 36, 37, 83),
									},
									{
										Name:          "varjob2",
										Value:         "valuejob2",
										FileReference: testutils.CreateFileReference(39, 310, 131, 132),
									},
								},
							},
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
					EnvironmentVariables: &models.EnvironmentVariablesRef{
						EnvironmentVariables: models.EnvironmentVariables{
							"var1":    "value1",
							"var2":    "value2",
							"varjob1": "valuejob1",
							"varjob2": "valuejob2",
						},
						FileReference: testutils.CreateFileReference(31, 23, 131, 132),
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				{
					ID:              utils.GetPtr("deployment-job"),
					Name:            utils.GetPtr(""),
					TimeoutMS:       utils.GetPtr(defaultTimeoutMS),
					ContinueOnError: utils.GetPtr(false),
					EnvironmentVariables: &models.EnvironmentVariablesRef{
						EnvironmentVariables: models.EnvironmentVariables{
							"var1": "value1",
							"var2": "value2",
						},
						FileReference: testutils.CreateFileReference(1, 2, 11, 12),
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseStage(testCase.stage)

			testutils.DeepCompare(t, testCase.expectedJobs, got)
		})
	}
}
