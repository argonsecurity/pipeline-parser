package bitbucket

import (
	"testing"

	bbModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestDefaultsParse(t *testing.T) {
	testCases := []struct {
		name              string
		bitbucketPipeline *bbModels.Pipeline
		expectedDefaults  *models.Defaults
	}{
		{
			name:              "Pipeline is nil",
			bitbucketPipeline: nil,
			expectedDefaults:  nil,
		},
		{
			name:              "image name is not defined",
			bitbucketPipeline: &bbModels.Pipeline{},
			expectedDefaults: &models.Defaults{
				Runner: nil,
			},
		},
		{
			name: "image name is defined",
			bitbucketPipeline: &bbModels.Pipeline{
				Image: &bbModels.Image{
					ImageData: &bbModels.ImageData{
						Name: utils.GetPtr("node:10.15.3"),
					},
				},
			},
			expectedDefaults: &models.Defaults{
				Runner: &models.Runner{
					DockerMetadata: &models.DockerMetadata{
						Image: utils.GetPtr("node:10.15.3"),
					},
				},
			},
		},
		{
			name: "pipeline has global settings",
			bitbucketPipeline: &bbModels.Pipeline{
				Image: &bbModels.Image{
					ImageData: &bbModels.ImageData{
						Name: utils.GetPtr("node:10.15.3"),
					},
				},
				Options: &bbModels.GlobalSettings{
					Docker:  utils.GetPtr(true),
					MaxTime: utils.GetPtr(int64(60)),
					Size:    utils.GetPtr(bbModels.X1),
				},
			},
			expectedDefaults: &models.Defaults{
				Runner: &models.Runner{
					DockerMetadata: &models.DockerMetadata{
						Image: utils.GetPtr("node:10.15.3"),
					},
				},
				Settings: &map[string]any{
					"docker":   utils.GetPtr(true),
					"max-time": utils.GetPtr(int64(60)),
					"size":     utils.GetPtr(bbModels.X1),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			defaults := parsePipelineDefaults(testCase.bitbucketPipeline)
			testutils.DeepCompare(t, testCase.expectedDefaults, defaults)
		})
	}
}
