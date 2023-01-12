package bitbucket

import (
	"strings"
	"testing"

	bbModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/r3labs/diff/v3"
)

func TestLoad(t *testing.T) {
	testCases := []struct {
		name             string
		filename         string
		expectedPipeline *bbModels.Pipeline
		expectedError    error
	}{
		{
			name:     "parallel steps",
			filename: "../../../test/fixtures/bitbucket/parallel-steps.yml",
			expectedPipeline: &bbModels.Pipeline{
				Image: &bbModels.Image{
					Name: "node:16",
				},
				Pipelines: &bbModels.BuildPipelines{
					Default: []*bbModels.Step{
						{
							Parallel: []*bbModels.ParallelSteps{
								{
									Step: &bbModels.ExecutionUnitRef{
										ExecutionUnit: &bbModels.ExecutionUnit{
											Name: "Build and Test",
											Caches: []string{
												"node",
											},
											Script: []bbModels.Script{
												{
													String:        "npm install",
													FileReference: testutils.CreateFileReference(11, 17, 11, 28),
												},
												{
													String:        "npm test",
													FileReference: testutils.CreateFileReference(12, 17, 12, 25),
												},
											},
										},
										FileReference: testutils.CreateFileReference(7, 13, 12, 25),
									},
								},
								{
									Step: &bbModels.ExecutionUnitRef{
										ExecutionUnit: &bbModels.ExecutionUnit{
											Name: "Code linting",
											Caches: []string{
												"node",
											},
											Script: []bbModels.Script{
												{
													String:        "npm install eslint",
													FileReference: testutils.CreateFileReference(16, 17, 16, 35),
												},
												{
													String:        "npx eslint .",
													FileReference: testutils.CreateFileReference(17, 17, 17, 29),
												},
											},
										},
										FileReference: testutils.CreateFileReference(14, 13, 19, 21),
									},
								},
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:     "sync steps",
			filename: "../../../test/fixtures/bitbucket/sync-steps.yml",
			expectedPipeline: &bbModels.Pipeline{
				Image: &bbModels.Image{
					Name: "node:16",
				},
				Pipelines: &bbModels.BuildPipelines{
					PullRequests: &bbModels.StepMap{
						"**": {
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: "Build and Test",
										Caches: []string{
											"node",
										},
										Script: []bbModels.Script{
											{
												String:        "npm install",
												FileReference: testutils.CreateFileReference(11, 15, 11, 26),
											},
											{
												String:        "npm test",
												FileReference: testutils.CreateFileReference(12, 15, 12, 23),
											},
										},
									},
									FileReference: testutils.CreateFileReference(6, 9, 12, 23),
								},
							},
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: "Code linting",
										Caches: []string{
											"node",
										},
										Script: []bbModels.Script{
											{
												String:        "npm install eslint",
												FileReference: testutils.CreateFileReference(16, 15, 16, 33),
											},
											{
												String:        "npx eslint .",
												FileReference: testutils.CreateFileReference(17, 15, 17, 27),
											},
										},
									},
									FileReference: testutils.CreateFileReference(13, 9, 19, 19),
								},
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:     "multiple pipelines definitions",
			filename: "../../../test/fixtures/bitbucket/multiple-pipelines-types.yml",
			expectedPipeline: &bbModels.Pipeline{
				Image: &bbModels.Image{
					Name: "node:16",
				},
				Pipelines: &bbModels.BuildPipelines{
					Custom: &bbModels.StepMap{
						"notify": {
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: "Notify Teams",
										Caches: []string{
											"node",
										},
										Script: []bbModels.Script{
											{
												String:        "npx notify -s \"deployment\"",
												FileReference: testutils.CreateFileReference(11, 15, 11, 41),
											},
										},
									},
									FileReference: testutils.CreateFileReference(6, 9, 11, 41),
								},
							},
						},
					},
					Branches: &bbModels.StepMap{
						"master": {
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: "step 1",
									},
									FileReference: testutils.CreateFileReference(14, 9, 15, 23),
								},
							},
							{
								Step: &bbModels.ExecutionUnitRef{
									ExecutionUnit: &bbModels.ExecutionUnit{
										Name: "step 2",
									},
									FileReference: testutils.CreateFileReference(16, 9, 17, 23),
								},
							},
							{
								Parallel: []*bbModels.ParallelSteps{
									{
										Step: &bbModels.ExecutionUnitRef{
											ExecutionUnit: &bbModels.ExecutionUnit{
												Name: "step 3",
											},
											FileReference: testutils.CreateFileReference(20, 15, 20, 27),
										},
									},
									{
										Step: &bbModels.ExecutionUnitRef{
											ExecutionUnit: &bbModels.ExecutionUnit{
												Name: "step 4",
											},
											FileReference: testutils.CreateFileReference(22, 15, 22, 27),
										},
									},
								},
							},
						},
					},
				},
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
