package bitbucket

import (
	"strings"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/r3labs/diff/v3"
)

func TestLoad(t *testing.T) {
	testCases := []struct {
		name             string
		filename         string
		expectedPipeline *models.Pipeline
		expectedError    error
	}{
		{
			name:             "simple",
			filename:         "../../../test/fixtures/bitbucket/simple.yml",
			expectedPipeline: &models.Pipeline{
				// Image: &models.Image{
				// 	Name: "node:16",
				// },
				// Pipelines: &models.BuildPipelines{
				// 	Default: []*models.Step{
				// 		{
				// 			Parallel: []*models.ParallelSteps{
				// 				{
				// 					Step: &models.BuildExecutionUnit{
				// 						Name: "Build and Test",
				// 						Caches: []string{
				// 							"node",
				// 						},
				// 						Script: []models.Script{
				// 							{
				// 								String: "npm install",
				// 							},
				// 							{String: "npm test"},
				// 						},
				// 					},
				// 				},
				// 				{
				// 					Step: &models.BuildExecutionUnit{
				// 						Name: "Code linting",
				// 						Caches: []string{
				// 							"node",
				// 						},
				// 						Script: []models.Script{
				// 							{
				// 								String: "npm install eslint",
				// 							},
				// 							{String: "npx eslint ."},
				// 						},
				// 					},
				// 				},
				// 			},
				// 		},
				// 	},
				// },
			},
			expectedError: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			loader := &BitbucketLoader{}
			pipeline, err := loader.Load(testutils.ReadFile(testCase.filename))

			if err != testCase.expectedError {
				t.Errorf("Expected error: %v, got: %v", testCase.expectedError, err)
			}

			changelog, _ := diff.Diff(pipeline, testCase.expectedPipeline)

			if len(changelog) > 0 {
				t.Errorf("Loader result is not as expected:")
				for _, change := range changelog {
					t.Errorf("field: %s, got: %v, expected: %v", strings.Join(change.Path, "."), change.From, change.To)
				}
			}
		})
	}
}
