package azure

import (
	"strings"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
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
		// {
		// 	name:     "all-triggers",
		// 	filename: "../../../test/fixtures/azure/all-triggers.yaml",
		// 	expectedPipeline: &models.Pipeline{
		// 		Name: "all-triggers",
		// 		Trigger: &models.TriggerRef{
		// 			Trigger: &models.Trigger{
		// 				Batch: true,
		// 				Branches: &models.Filter{
		// 					Include: []string{"master", "main"},
		// 					Exclude: []string{"test/*"},
		// 				},
		// 				Paths: &models.Filter{
		// 					Include: []string{"path/to/file", "another/path/to/file"},
		// 					Exclude: []string{"all/*"},
		// 				},
		// 				Tags: &models.Filter{
		// 					Include: []string{"v1.0.*"},
		// 					Exclude: []string{"v2.0.*"},
		// 				},
		// 			},
		// 			FileReference: testutils.CreateFileReference(2, 1, 20, 15),
		// 		},
		// 		PR: &models.PRRef{
		// 			PR: &models.PR{
		// 				AutoCancel: true,
		// 				Branches: &models.Filter{
		// 					Include: []string{"features/*"},
		// 					Exclude: []string{"features/experimental/*"},
		// 				},
		// 				Paths: &models.Filter{
		// 					Include: []string{"path/to/file"},
		// 					Exclude: []string{"README.md"},
		// 				},
		// 				Drafts: true,
		// 			},
		// 			FileReference: testutils.CreateFileReference(21, 1, 33, 15),
		// 		},
		// 	},
		// },
		// {
		// 	name:     "no-trigger",
		// 	filename: "../../../test/fixtures/azure/no-trigger.yaml",
		// 	expectedPipeline: &models.Pipeline{
		// 		Name:    "no-trigger",
		// 		Trigger: &models.TriggerRef{FileReference: testutils.CreateFileReference(2, 10, 2, 14)},
		// 		PR:      &models.PRRef{FileReference: testutils.CreateFileReference(3, 5, 3, 9)},
		// 	},
		// },
		{
			name:     "branch-list-trigger",
			filename: "../../../test/fixtures/azure/branch-list-trigger.yaml",
			expectedPipeline: &models.Pipeline{
				Name: "branch-list-trigger",
				Trigger: &models.TriggerRef{
					Trigger: &models.Trigger{
						Branches: &models.Filter{
							Include: []string{"main", "development"},
						},
					},
					FileReference: testutils.CreateFileReference(2, 1, 4, 16),
				},
				PR: &models.PRRef{
					PR: &models.PR{
						Branches: &models.Filter{
							Include: []string{"main", "develop"},
						},
					},
					FileReference: testutils.CreateFileReference(5, -1, 7, 10),
				},
			},
		},
		// {
		// 	name:     "variables",
		// 	filename: "../../../test/fixtures/azure/variables.yaml",
		// 	expectedPipeline: &models.Pipeline{
		// 		Name: "variables",
		// 		Variables: &models.Variables{
		// 			{
		// 				Name:  "MY_VAR",
		// 				Value: "my value",
		// 			},
		// 			{
		// 				Name:  "ANOTHER_VAR",
		// 				Value: "another value",
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name:     "pool",
		// 	filename: "../../../test/fixtures/azure/pool.yaml",
		// 	expectedPipeline: &models.Pipeline{
		// 		Name: "pool",
		// 		Pool: &models.Pool{
		// 			Name:    "MyPool",
		// 			Demands: []string{"demand1", "demand2"},
		// 			VmImage: "ubuntu-latest",
		// 		},
		// 	},
		// },
		// {
		// 	name:             "",
		// 	filename:         "../../../test/fixtures/azure/parameters.yaml",
		// 	expectedPipeline: &models.Pipeline{},
		// },
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			loader := &AzureLoader{}
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
