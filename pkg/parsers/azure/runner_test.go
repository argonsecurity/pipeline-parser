package azure

import (
	"testing"

	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/r3labs/diff/v3"
	"github.com/stretchr/testify/assert"
)

func TestParseRunner(t *testing.T) {
	testCases := []struct {
		name           string
		job            *azureModels.BaseJob
		expectedRunner *models.Runner
	}{
		{
			name:           "job is nil",
			job:            nil,
			expectedRunner: nil,
		},
		{
			name:           "job is empty",
			job:            &azureModels.BaseJob{},
			expectedRunner: nil,
		},
		{
			name: "job with pool and no container",
			job: &azureModels.BaseJob{
				Pool: &azureModels.Pool{
					VmImage: "ubuntu-18.04",
				},
			},
			expectedRunner: &models.Runner{
				OS: utils.GetPtr("linux"),
			},
		},
		{
			name: "job with container and no pool",
			job: &azureModels.BaseJob{
				Container: &azureModels.JobContainer{
					Image: "ubuntu:18.04",
				},
			},
			expectedRunner: &models.Runner{
				DockerMetadata: &models.DockerMetadata{
					Image: utils.GetPtr("ubuntu"),
					Label: utils.GetPtr("18.04"),
				},
			},
		},
		{
			name: "job with container and pool",
			job: &azureModels.BaseJob{
				Pool: &azureModels.Pool{
					VmImage: "ubuntu-18.04",
				},
				Container: &azureModels.JobContainer{
					Image: "ubuntu:18.04",
				},
			},
			expectedRunner: &models.Runner{
				OS: utils.GetPtr("linux"),
				DockerMetadata: &models.DockerMetadata{
					Image: utils.GetPtr("ubuntu"),
					Label: utils.GetPtr("18.04"),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseRunner(testCase.job)

			changelog, err := diff.Diff(testCase.expectedRunner, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}
}

func TestParsePool(t *testing.T) {
	testCases := []struct {
		name           string
		pool           *azureModels.Pool
		runner         *models.Runner
		expectedRunner *models.Runner
	}{
		{
			name:           "pool is nil",
			pool:           nil,
			expectedRunner: nil,
		},
		{
			name:           "pool is empty",
			pool:           &azureModels.Pool{},
			expectedRunner: nil,
		},
		{
			name:           "runner is nil",
			runner:         nil,
			expectedRunner: nil,
		},
		{
			name: "pool with data, runner is empty",
			pool: &azureModels.Pool{
				VmImage: "ubuntu-18.04",
			},
			runner: &models.Runner{},
			expectedRunner: &models.Runner{
				OS: utils.GetPtr("linux"),
			},
		},
		{
			name: "pool with data, runner is not empty empty",
			pool: &azureModels.Pool{
				VmImage: "ubuntu-18.04",
			},
			runner: &models.Runner{
				DockerMetadata: &models.DockerMetadata{
					Image: utils.GetPtr("ubuntu"),
					Label: utils.GetPtr("18.04"),
				},
			},
			expectedRunner: &models.Runner{
				OS: utils.GetPtr("linux"),
				DockerMetadata: &models.DockerMetadata{
					Image: utils.GetPtr("ubuntu"),
					Label: utils.GetPtr("18.04"),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parsePool(testCase.pool, testCase.runner)

			changelog, err := diff.Diff(testCase.expectedRunner, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}

}

func TestParseContainer(t *testing.T) {
	testCases := []struct {
		name           string
		container      *azureModels.JobContainer
		runner         *models.Runner
		expectedRunner *models.Runner
	}{
		{
			name:           "container is nil",
			container:      nil,
			expectedRunner: nil,
		},
		{
			name:           "container is empty",
			container:      &azureModels.JobContainer{},
			expectedRunner: nil,
		},
		{
			name:           "runner is nil",
			runner:         nil,
			expectedRunner: nil,
		},
		{
			name: "container with data, runner is empty",
			container: &azureModels.JobContainer{
				Image: "ubuntu:18.04",
			},
			runner: &models.Runner{},
			expectedRunner: &models.Runner{
				DockerMetadata: &models.DockerMetadata{
					Image: utils.GetPtr("ubuntu"),
					Label: utils.GetPtr("18.04"),
				},
			},
		},
		{
			name: "container with data, runner is not empty empty",
			container: &azureModels.JobContainer{
				Image: "ubuntu:18.04",
			},
			runner: &models.Runner{
				OS:   utils.GetPtr("linux"),
				Arch: utils.GetPtr("x86_64"),
			},
			expectedRunner: &models.Runner{
				OS:   utils.GetPtr("linux"),
				Arch: utils.GetPtr("x86_64"),
				DockerMetadata: &models.DockerMetadata{
					Image: utils.GetPtr("ubuntu"),
					Label: utils.GetPtr("18.04"),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseContainer(testCase.container, testCase.runner)

			changelog, err := diff.Diff(testCase.expectedRunner, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}
}
