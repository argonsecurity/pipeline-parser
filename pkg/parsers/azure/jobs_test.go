package azure

import (
	"testing"

	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
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
							Pool: &azureModels.Pool{
								VmImage: "ubuntu-18.04",
							},
							Container: &azureModels.JobContainer{
								Image: "ubuntu:18.04",
							},
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
					ContinueOnError: utils.GetPtr("true"),
					TimeoutMS:       utils.GetPtr(6000000),
					Conditions:      []*models.Condition{{Statement: "job-1-condition"}},
					Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job-2")}},
					Runner: &models.Runner{
						OS: utils.GetPtr("linux"),
						DockerMetadata: &models.DockerMetadata{
							Image: utils.GetPtr("ubuntu"),
							Label: utils.GetPtr("18.04"),
						},
					},
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
					ContinueOnError: utils.GetPtr("true"),
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
					ContinueOnError: utils.GetPtr("true"),
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
					ContinueOnError: utils.GetPtr("true"),
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

			testutils.DeepCompare(t, testCase.expectedJobs, got)
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
				ContinueOnError: utils.GetPtr("false"),
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
					Pool: &azureModels.Pool{
						VmImage: "ubuntu-18.04",
					},
					Container: &azureModels.JobContainer{
						Image: "ubuntu:18.04",
					},
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
				ContinueOnError: utils.GetPtr("true"),
				TimeoutMS:       utils.GetPtr(6000000),
				Conditions:      []*models.Condition{{Statement: "job-1-condition"}},
				Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job-2")}},
				Runner: &models.Runner{
					OS: utils.GetPtr("linux"),
					DockerMetadata: &models.DockerMetadata{
						Image: utils.GetPtr("ubuntu"),
						Label: utils.GetPtr("18.04"),
					},
				},
				EnvironmentVariables: &models.EnvironmentVariablesRef{
					EnvironmentVariables: models.EnvironmentVariables{
						"var1": "value1",
						"var2": "value2",
					},
					FileReference: testutils.CreateFileReference(1, 2, 11, 12),
				},
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
			testutils.DeepCompare(t, testCase.expectedJob, got)
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
				ContinueOnError: utils.GetPtr("false"),
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
					Pool: &azureModels.Pool{
						VmImage: "ubuntu-18.04",
					},
					Container: &azureModels.JobContainer{
						Image: "ubuntu:18.04",
					},
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
				ContinueOnError: utils.GetPtr("true"),
				TimeoutMS:       utils.GetPtr(6000000),
				Conditions:      []*models.Condition{{Statement: "job-1-condition"}},
				Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job-2")}},
				Runner: &models.Runner{
					OS: utils.GetPtr("linux"),
					DockerMetadata: &models.DockerMetadata{
						Image: utils.GetPtr("ubuntu"),
						Label: utils.GetPtr("18.04"),
					},
				},
				EnvironmentVariables: &models.EnvironmentVariablesRef{
					EnvironmentVariables: models.EnvironmentVariables{
						"var1": "value1",
						"var2": "value2",
					},
					FileReference: testutils.CreateFileReference(1, 2, 11, 12),
				},
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
			testutils.DeepCompare(t, testCase.expectedJob, got)
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
				ContinueOnError: utils.GetPtr("false"),
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
				Pool: &azureModels.Pool{
					VmImage: "ubuntu-18.04",
				},
				Container: &azureModels.JobContainer{
					Image: "ubuntu:18.04",
				},
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
				ContinueOnError: utils.GetPtr("true"),
				TimeoutMS:       utils.GetPtr(6000000),
				Conditions:      []*models.Condition{{Statement: "job-1-condition"}},
				Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job-2")}},
				Runner: &models.Runner{
					OS: utils.GetPtr("linux"),
					DockerMetadata: &models.DockerMetadata{
						Image: utils.GetPtr("ubuntu"),
						Label: utils.GetPtr("18.04"),
					},
				},
				EnvironmentVariables: &models.EnvironmentVariablesRef{
					EnvironmentVariables: models.EnvironmentVariables{
						"var1": "value1",
						"var2": "value2",
					},
					FileReference: testutils.CreateFileReference(1, 2, 11, 12),
				},
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
			testutils.DeepCompare(t, testCase.expectedJob, got)
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
			testutils.DeepCompare(t, testCase.expectedDependencies, got)
		})
	}
}
