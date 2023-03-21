package general

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/enhancers/general/config"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestEnhanceJob(t *testing.T) {
	testCases := []struct {
		name        string
		job         *models.Job
		config      *config.EnhancementConfiguration
		expectedJob *models.Job
	}{
		{
			name: "Job name contains build (lowercase)",
			job: &models.Job{
				Name: utils.GetPtr("build app"),
			},
			config: config.CommonConfiguration,
			expectedJob: &models.Job{
				Name: utils.GetPtr("build app"),
				Metadata: models.Metadata{
					Build: true,
				},
			},
		},
		{
			name: "Job name contains Build (uppercase)",
			job: &models.Job{
				Name: utils.GetPtr("Build app"),
			},
			config: config.CommonConfiguration,
			expectedJob: &models.Job{
				Name: utils.GetPtr("Build app"),
				Metadata: models.Metadata{
					Build: true,
				},
			},
		},
		{
			name: "Job contains step with build",
			job: &models.Job{
				Steps: []*models.Step{
					{
						Metadata: models.Metadata{
							Build: true,
						},
					},
				},
			},
			config: config.CommonConfiguration,
			expectedJob: &models.Job{
				Steps: []*models.Step{
					{
						Metadata: models.Metadata{
							Build: true,
						},
					},
				},
				Metadata: models.Metadata{
					Build: true,
				},
			},
		},
		{
			name: "Job name contains test (lowercase)",
			job: &models.Job{
				Name: utils.GetPtr("test app"),
			},
			config: config.CommonConfiguration,
			expectedJob: &models.Job{
				Name: utils.GetPtr("test app"),
				Metadata: models.Metadata{
					Test: true,
				},
			},
		},
		{
			name: "Job name contains tests (lowercase)",
			job: &models.Job{
				Name: utils.GetPtr("tests app"),
			},
			config: config.CommonConfiguration,
			expectedJob: &models.Job{
				Name: utils.GetPtr("tests app"),
				Metadata: models.Metadata{
					Test: true,
				},
			},
		},
		{
			name: "Job name contains Test (uppercase)",
			job: &models.Job{
				Name: utils.GetPtr("Test app"),
			},
			config: config.CommonConfiguration,
			expectedJob: &models.Job{
				Name: utils.GetPtr("Test app"),
				Metadata: models.Metadata{
					Test: true,
				},
			},
		},
		{
			name: "Job name contains Tests (uppercase)",
			job: &models.Job{
				Name: utils.GetPtr("Tests app"),
			},
			config: config.CommonConfiguration,
			expectedJob: &models.Job{
				Name: utils.GetPtr("Tests app"),
				Metadata: models.Metadata{
					Test: true,
				},
			},
		},
		{
			name: "Job contains step with test",
			job: &models.Job{
				Steps: []*models.Step{
					{
						Metadata: models.Metadata{
							Test: true,
						},
					},
				},
			},
			config: config.CommonConfiguration,
			expectedJob: &models.Job{
				Steps: []*models.Step{
					{
						Metadata: models.Metadata{
							Test: true,
						},
					},
				},
				Metadata: models.Metadata{
					Test: true,
				},
			},
		},
		{
			name: "Job name doesn't contain build test or deploy",
			job: &models.Job{
				Name: utils.GetPtr("some job"),
			},
			config: config.CommonConfiguration,
			expectedJob: &models.Job{
				Name:     utils.GetPtr("some job"),
				Metadata: models.Metadata{},
			},
		},
		{
			name: "Job doesn't contain step with build test or deploy",
			job: &models.Job{
				Steps: []*models.Step{
					{
						Metadata: models.Metadata{},
					},
				},
			},
			config: config.CommonConfiguration,
			expectedJob: &models.Job{
				Steps: []*models.Step{
					{
						Metadata: models.Metadata{},
					},
				},
				Metadata: models.Metadata{},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			enhancedJob := enhanceJob(testCase.job, testCase.config)

			assert.Equal(t, testCase.expectedJob, enhancedJob, testCase.name)
		})

	}
}
