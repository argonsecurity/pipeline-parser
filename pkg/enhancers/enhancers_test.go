package enhancers

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestEnhance(t *testing.T) {
	testCases := []struct {
		name             string
		pipeline         *models.Pipeline
		platform         consts.Platform
		expectedPipeline *models.Pipeline
	}{
		{
			name: "Pipeline contains job with name contains build (lowercase)",
			pipeline: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("build app"),
					},
				},
			},
			platform: consts.GitHubPlatform,
			expectedPipeline: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("build app"),
						Metadata: models.Metadata{
							Build: true,
						},
					},
				},
			},
		},
		{
			name: "Pipeline contains job with name contains Build (uppercase)",
			pipeline: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("Build app"),
					},
				},
			},
			platform: consts.GitHubPlatform,
			expectedPipeline: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("Build app"),
						Metadata: models.Metadata{
							Build: true,
						},
					},
				},
			},
		},
		{
			name: "Pipeline contains job with name doesn't contain build",
			pipeline: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("some job"),
					},
				},
			},
			platform: consts.GitHubPlatform,
			expectedPipeline: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("some job"),
						Metadata: models.Metadata{
							Build: false,
						},
					},
				},
			},
		},
		{
			name: "Pipeline contains job with step with name contains build (lowercase)",
			pipeline: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("build app"),
							},
						},
					},
				},
			},
			platform: consts.GitHubPlatform,
			expectedPipeline: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("build app"),
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
			},
		},
		{
			name: "Pipeline contains job with step with name contains Build (uppercase)",
			pipeline: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("Build app"),
							},
						},
					},
				},
			},
			platform: consts.GitHubPlatform,
			expectedPipeline: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("Build app"),
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
			},
		},
		{
			name: "Pipeline contains job with step with name doesn't contain build",
			pipeline: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("some job"),
							},
						},
					},
				},
			},
			platform: consts.GitHubPlatform,
			expectedPipeline: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("some job"),
								Metadata: models.Metadata{
									Build: false,
								},
							},
						},
						Metadata: models.Metadata{
							Build: false,
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			enhancedPipeline, err := Enhance(testCase.pipeline, testCase.platform)

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedPipeline, enhancedPipeline)
		})

	}
}
