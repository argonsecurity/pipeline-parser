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

func TestParseJobs(t *testing.T) {
	testCases := []struct {
		name         string
		jobs         *azureModels.Jobs
		expectedJobs []*models.Job
	}{
		{
			name:         "Jobs is nil",
			jobs:         nil,
			expectedJobs: nil,
		},
		{
			name:         "Empty jobs",
			jobs:         &azureModels.Jobs{},
			expectedJobs: nil,
		},
		{
			name: "Jobs with data",
			jobs: &azureModels.Jobs{
				CIJobs: []*azureModels.CIJob{
					{
						Job: "job-1",
						BaseJob: azureModels.BaseJob{
							DisplayName:      "job-1",
							DependsOn:        &azureModels.DependsOn{"job-2"},
							Condition:        "job-1-condition",
							ContinueOnError:  true,
							TimeoutInMinutes: 100,
							Steps: &azureModels.Steps{
								{
									Name:      "step-1",
									Bash:      "script",
									Condition: "step-1-condition",
								},
							},
						},
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
					{
						Job: "job-2",
						BaseJob: azureModels.BaseJob{
							DisplayName:      "job-2",
							DependsOn:        &azureModels.DependsOn{"job-3"},
							Condition:        "job-2-condition",
							ContinueOnError:  true,
							TimeoutInMinutes: 100,
							Steps: &azureModels.Steps{
								{
									Name:      "step-2",
									Bash:      "script",
									Condition: "step-2-condition",
								},
							},
						},
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
				},
				DeploymentJobs: []*azureModels.DeploymentJob{
					{
						Deployment: "deployment-1",
						BaseJob: azureModels.BaseJob{
							DisplayName:      "job-1",
							DependsOn:        &azureModels.DependsOn{"job-2"},
							Condition:        "job-1-condition",
							ContinueOnError:  true,
							TimeoutInMinutes: 100,
							Steps: &azureModels.Steps{
								{
									Name:      "step-1",
									Bash:      "script",
									Condition: "step-1-condition",
								},
							},
						},
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
					{
						Deployment: "deployment-2",
						BaseJob: azureModels.BaseJob{
							DisplayName:      "job-2",
							DependsOn:        &azureModels.DependsOn{"job-3"},
							Condition:        "job-2-condition",
							ContinueOnError:  true,
							TimeoutInMinutes: 100,
							Steps: &azureModels.Steps{
								{
									Name:      "step-2",
									Bash:      "script",
									Condition: "step-2-condition",
								},
							},
						},
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
				},
			},
			expectedJobs: []*models.Job{
				{
					ID:              utils.GetPtr("job-1"),
					Name:            utils.GetPtr("job-1"),
					ContinueOnError: utils.GetPtr(true),
					TimeoutMS:       utils.GetPtr(6000000),
					Conditions:      []*models.Condition{{Statement: "job-1-condition"}},
					Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job-2")}},
					Steps: []*models.Step{
						{
							ID:   utils.GetPtr("step-1"),
							Name: utils.GetPtr(""),
							Type: models.ShellStepType,
							Shell: &models.Shell{
								Type:   utils.GetPtr("bash"),
								Script: utils.GetPtr("script"),
							},
							Conditions: &[]models.Condition{{Statement: "step-1-condition"}},
						},
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				{
					ID:              utils.GetPtr("job-2"),
					Name:            utils.GetPtr("job-2"),
					ContinueOnError: utils.GetPtr(true),
					TimeoutMS:       utils.GetPtr(6000000),
					Conditions:      []*models.Condition{{Statement: "job-2-condition"}},
					Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job-3")}},
					Steps: []*models.Step{
						{
							ID:   utils.GetPtr("step-2"),
							Name: utils.GetPtr(""),
							Type: models.ShellStepType,
							Shell: &models.Shell{
								Type:   utils.GetPtr("bash"),
								Script: utils.GetPtr("script"),
							},
							Conditions: &[]models.Condition{{Statement: "step-2-condition"}},
						},
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				{
					ID:              utils.GetPtr("deployment-1"),
					Name:            utils.GetPtr("job-1"),
					ContinueOnError: utils.GetPtr(true),
					TimeoutMS:       utils.GetPtr(6000000),
					Conditions:      []*models.Condition{{Statement: "job-1-condition"}},
					Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job-2")}},
					Steps: []*models.Step{
						{
							ID:   utils.GetPtr("step-1"),
							Name: utils.GetPtr(""),
							Type: models.ShellStepType,
							Shell: &models.Shell{
								Type:   utils.GetPtr("bash"),
								Script: utils.GetPtr("script"),
							},
							Conditions: &[]models.Condition{{Statement: "step-1-condition"}},
						},
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				{
					ID:              utils.GetPtr("deployment-2"),
					Name:            utils.GetPtr("job-2"),
					ContinueOnError: utils.GetPtr(true),
					TimeoutMS:       utils.GetPtr(6000000),
					Conditions:      []*models.Condition{{Statement: "job-2-condition"}},
					Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job-3")}},
					Steps: []*models.Step{
						{
							ID:   utils.GetPtr("step-2"),
							Name: utils.GetPtr(""),
							Type: models.ShellStepType,
							Shell: &models.Shell{
								Type:   utils.GetPtr("bash"),
								Script: utils.GetPtr("script"),
							},
							Conditions: &[]models.Condition{{Statement: "step-2-condition"}},
						},
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseJobs(testCase.jobs)

			changelog, err := diff.Diff(testCase.expectedJobs, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}
}

func TestParseCIJob(t *testing.T) {
	testCases := []struct {
		name        string
		ciJob       *azureModels.CIJob
		expectedJob *models.Job
	}{
		{
			name:        "CIJob is nil",
			ciJob:       nil,
			expectedJob: nil,
		},
		{
			name:  "Empty CIJob",
			ciJob: &azureModels.CIJob{},
			expectedJob: &models.Job{
				ID:              utils.GetPtr(""),
				Name:            utils.GetPtr(""),
				ContinueOnError: utils.GetPtr(false),
				TimeoutMS:       utils.GetPtr(defaultTimeoutMS),
			},
		},
		{
			name: "CIJob with data",
			ciJob: &azureModels.CIJob{
				Job: "job-1",
				BaseJob: azureModels.BaseJob{
					DisplayName:      "job-1",
					DependsOn:        &azureModels.DependsOn{"job-2"},
					Condition:        "job-1-condition",
					ContinueOnError:  true,
					TimeoutInMinutes: 100,
					Steps: &azureModels.Steps{
						{
							Name:      "step-1",
							Bash:      "script",
							Condition: "step-1-condition",
						},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedJob: &models.Job{
				ID:              utils.GetPtr("job-1"),
				Name:            utils.GetPtr("job-1"),
				ContinueOnError: utils.GetPtr(true),
				TimeoutMS:       utils.GetPtr(6000000),
				Conditions:      []*models.Condition{{Statement: "job-1-condition"}},
				Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job-2")}},
				Steps: []*models.Step{
					{
						ID:   utils.GetPtr("step-1"),
						Name: utils.GetPtr(""),
						Type: models.ShellStepType,
						Shell: &models.Shell{
							Type:   utils.GetPtr("bash"),
							Script: utils.GetPtr("script"),
						},
						Conditions: &[]models.Condition{{Statement: "step-1-condition"}},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseCIJob(testCase.ciJob)

			changelog, err := diff.Diff(testCase.expectedJob, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}
}

func TestParseDeploymentJob(t *testing.T) {
	testCases := []struct {
		name          string
		deploymentJob *azureModels.DeploymentJob
		expectedJob   *models.Job
	}{
		{
			name:          "DeploymentJob is nil",
			deploymentJob: nil,
			expectedJob:   nil,
		},
		{
			name:          "Empty deploymentJob",
			deploymentJob: &azureModels.DeploymentJob{},
			expectedJob: &models.Job{
				ID:              utils.GetPtr(""),
				Name:            utils.GetPtr(""),
				ContinueOnError: utils.GetPtr(false),
				TimeoutMS:       utils.GetPtr(defaultTimeoutMS),
			},
		},
		{
			name: "DeploymentJob with data",
			deploymentJob: &azureModels.DeploymentJob{
				Deployment: "deployment-1",
				BaseJob: azureModels.BaseJob{
					DisplayName:      "job-1",
					DependsOn:        &azureModels.DependsOn{"job-2"},
					Condition:        "job-1-condition",
					ContinueOnError:  true,
					TimeoutInMinutes: 100,
					Steps: &azureModels.Steps{
						{
							Name:      "step-1",
							Bash:      "script",
							Condition: "step-1-condition",
						},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedJob: &models.Job{
				ID:              utils.GetPtr("deployment-1"),
				Name:            utils.GetPtr("job-1"),
				ContinueOnError: utils.GetPtr(true),
				TimeoutMS:       utils.GetPtr(6000000),
				Conditions:      []*models.Condition{{Statement: "job-1-condition"}},
				Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job-2")}},
				Steps: []*models.Step{
					{
						ID:   utils.GetPtr("step-1"),
						Name: utils.GetPtr(""),
						Type: models.ShellStepType,
						Shell: &models.Shell{
							Type:   utils.GetPtr("bash"),
							Script: utils.GetPtr("script"),
						},
						Conditions: &[]models.Condition{{Statement: "step-1-condition"}},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseDeploymentJob(testCase.deploymentJob)

			changelog, err := diff.Diff(testCase.expectedJob, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}
}

func TestParseBaseJob(t *testing.T) {
	testCases := []struct {
		name        string
		baseJob     *azureModels.BaseJob
		expectedJob *models.Job
	}{
		{
			name:        "BaseJob is nil",
			baseJob:     nil,
			expectedJob: nil,
		},
		{
			name:    "Empty BaseJob",
			baseJob: &azureModels.BaseJob{},
			expectedJob: &models.Job{
				Name:            utils.GetPtr(""),
				ContinueOnError: utils.GetPtr(false),
				TimeoutMS:       utils.GetPtr(defaultTimeoutMS),
			},
		},
		{
			name: "BaseJob with data",
			baseJob: &azureModels.BaseJob{
				DisplayName:      "job-1",
				DependsOn:        &azureModels.DependsOn{"job-2"},
				Condition:        "job-1-condition",
				ContinueOnError:  true,
				TimeoutInMinutes: 100,
				Steps: &azureModels.Steps{
					{
						Name:      "step-1",
						Bash:      "script",
						Condition: "step-1-condition",
					},
				},
			},
			expectedJob: &models.Job{
				Name:            utils.GetPtr("job-1"),
				ContinueOnError: utils.GetPtr(true),
				TimeoutMS:       utils.GetPtr(6000000),
				Conditions:      []*models.Condition{{Statement: "job-1-condition"}},
				Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job-2")}},
				Steps: []*models.Step{
					{
						ID:   utils.GetPtr("step-1"),
						Name: utils.GetPtr(""),
						Type: models.ShellStepType,
						Shell: &models.Shell{
							Type:   utils.GetPtr("bash"),
							Script: utils.GetPtr("script"),
						},
						Conditions: &[]models.Condition{{Statement: "step-1-condition"}},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseBaseJob(testCase.baseJob)

			changelog, err := diff.Diff(testCase.expectedJob, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}
}

func TestParseDependencies(t *testing.T) {
	testCases := []struct {
		name                 string
		dependsOn            *azureModels.DependsOn
		expectedDependencies []*models.JobDependency
	}{
		{
			name:                 "DependsOn is nill",
			dependsOn:            &azureModels.DependsOn{},
			expectedDependencies: []*models.JobDependency{},
		},
		{
			name:                 "Empty dependsOn",
			dependsOn:            &azureModels.DependsOn{},
			expectedDependencies: []*models.JobDependency{},
		},
		{
			name: "DependsOn with one dependency",
			dependsOn: &azureModels.DependsOn{
				"job-1",
			},
			expectedDependencies: []*models.JobDependency{
				{
					JobID: utils.GetPtr("job-1"),
				},
			},
		},
		{
			name: "DependsOn with some dependencies",
			dependsOn: &azureModels.DependsOn{
				"job-1",
				"job-2",
				"job-3",
			},
			expectedDependencies: []*models.JobDependency{
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
			got := parseDependencies(testCase.dependsOn)

			changelog, err := diff.Diff(testCase.expectedDependencies, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0)
		})
	}
}
