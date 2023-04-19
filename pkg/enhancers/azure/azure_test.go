package azure

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_mergePipelines(t *testing.T) {
	type args struct {
		Pipeline          *models.Pipeline
		ImportedPipelines *enhancers.ImportedPipeline
	}
	tests := []struct {
		name string
		args args
		want *models.Pipeline
	}{
		{
			name: "nil input",
			args: args{
				Pipeline:          nil,
				ImportedPipelines: nil,
			},
			want: nil,
		},
		{
			name: "imported pipeline not found",
			args: args{
				Pipeline: &models.Pipeline{
					Imports: []*models.Import{
						{
							Source: &models.ImportSource{
								Path: utils.GetPtr("test"),
								Type: models.SourceTypeLocal,
							},
							FileReference: testutils.CreateFileReference(1, 1, 1, 1),
						},
					},
				},
				ImportedPipelines: &enhancers.ImportedPipeline{
					Pipeline: &models.Pipeline{
						Name: utils.GetPtr("test"),
					},
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 2),
				},
			},
			want: &models.Pipeline{
				Imports: []*models.Import{
					{
						Source: &models.ImportSource{
							Path: utils.GetPtr("test"),
							Type: models.SourceTypeLocal,
						},
						FileReference: testutils.CreateFileReference(1, 1, 1, 1),
					},
				},
			},
		},
		{
			name: "map imports",
			args: args{
				Pipeline: &models.Pipeline{
					Imports: []*models.Import{
						{
							Source: &models.ImportSource{
								Path: utils.GetPtr("test"),
								Type: models.SourceTypeLocal,
							},
							FileReference: testutils.CreateFileReference(1, 1, 1, 1),
						},
					},
				},
				ImportedPipelines: &enhancers.ImportedPipeline{
					Pipeline: &models.Pipeline{
						Name: utils.GetPtr("test"),
					},
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 1),
					JobName:             "test",
				},
			},
			want: &models.Pipeline{
				Imports: []*models.Import{
					{
						Source: &models.ImportSource{
							Path: utils.GetPtr("test"),
							Type: models.SourceTypeLocal,
						},
						FileReference: testutils.CreateFileReference(1, 1, 1, 1),
						Pipeline: &models.Pipeline{
							Name: utils.GetPtr("test"),
						},
					},
				},
			},
		},
		{
			name: "map jobs",
			args: args{
				Pipeline: &models.Pipeline{
					Jobs: []*models.Job{
						{
							Imports: &models.Import{
								Source: &models.ImportSource{
									Path: utils.GetPtr("test"),
									Type: models.SourceTypeLocal,
								},
								FileReference: testutils.CreateFileReference(1, 1, 1, 1),
							},
						},
					},
				},
				ImportedPipelines: &enhancers.ImportedPipeline{
					Pipeline: &models.Pipeline{
						Name: utils.GetPtr("test"),
					},
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 1),
					JobName:             "test",
				},
			},
			want: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Imports: &models.Import{
							Source: &models.ImportSource{
								Path: utils.GetPtr("test"),
								Type: models.SourceTypeLocal,
							},
							FileReference: testutils.CreateFileReference(1, 1, 1, 1),
							Pipeline: &models.Pipeline{
								Name: utils.GetPtr("test"),
							},
						},
					},
				},
			},
		},
		{
			name: "map job envs",
			args: args{
				Pipeline: &models.Pipeline{
					Jobs: []*models.Job{
						{
							EnvironmentVariables: &models.EnvironmentVariablesRef{
								Imports: &models.Import{
									Source: &models.ImportSource{
										Path: utils.GetPtr("test"),
										Type: models.SourceTypeLocal,
									},
									FileReference: testutils.CreateFileReference(1, 1, 1, 1),
								},
							},
						},
					},
				},
				ImportedPipelines: &enhancers.ImportedPipeline{
					Pipeline: &models.Pipeline{
						Name: utils.GetPtr("test"),
					},
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 1),
					JobName:             "test",
				},
			},
			want: &models.Pipeline{
				Jobs: []*models.Job{
					{
						EnvironmentVariables: &models.EnvironmentVariablesRef{
							Imports: &models.Import{
								Source: &models.ImportSource{
									Path: utils.GetPtr("test"),
									Type: models.SourceTypeLocal,
								},
								FileReference: testutils.CreateFileReference(1, 1, 1, 1),
								Pipeline: &models.Pipeline{
									Name: utils.GetPtr("test"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "map steps",
			args: args{
				Pipeline: &models.Pipeline{
					Jobs: []*models.Job{
						{
							Steps: []*models.Step{
								{
									Imports: &models.Import{
										Source: &models.ImportSource{
											Path: utils.GetPtr("test"),
											Type: models.SourceTypeLocal,
										},
										FileReference: testutils.CreateFileReference(1, 1, 1, 1),
									},
								},
							},
						},
					},
				},
				ImportedPipelines: &enhancers.ImportedPipeline{
					Pipeline: &models.Pipeline{
						Name: utils.GetPtr("test"),
					},
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 1),
					JobName:             "test",
				},
			},
			want: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Steps: []*models.Step{
							{
								Imports: &models.Import{
									Source: &models.ImportSource{
										Path: utils.GetPtr("test"),
										Type: models.SourceTypeLocal,
									},
									FileReference: testutils.CreateFileReference(1, 1, 1, 1),
									Pipeline: &models.Pipeline{
										Name: utils.GetPtr("test"),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "map step envs",
			args: args{
				Pipeline: &models.Pipeline{
					Jobs: []*models.Job{
						{
							Steps: []*models.Step{
								{
									EnvironmentVariables: &models.EnvironmentVariablesRef{
										Imports: &models.Import{
											Source: &models.ImportSource{
												Path: utils.GetPtr("test"),
												Type: models.SourceTypeLocal,
											},
											FileReference: testutils.CreateFileReference(1, 1, 1, 1),
										},
									},
								},
							},
						},
					},
				},
				ImportedPipelines: &enhancers.ImportedPipeline{
					Pipeline: &models.Pipeline{
						Name: utils.GetPtr("test"),
					},
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 1),
					JobName:             "test",
				},
			},
			want: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Steps: []*models.Step{
							{
								EnvironmentVariables: &models.EnvironmentVariablesRef{
									Imports: &models.Import{
										Source: &models.ImportSource{
											Path: utils.GetPtr("test"),
											Type: models.SourceTypeLocal,
										},
										FileReference: testutils.CreateFileReference(1, 1, 1, 1),
										Pipeline: &models.Pipeline{
											Name: utils.GetPtr("test"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "map pre steps",
			args: args{
				Pipeline: &models.Pipeline{
					Jobs: []*models.Job{
						{
							PreSteps: []*models.Step{
								{
									Imports: &models.Import{
										Source: &models.ImportSource{
											Path: utils.GetPtr("test"),
											Type: models.SourceTypeLocal,
										},
										FileReference: testutils.CreateFileReference(1, 1, 1, 1),
									},
								},
							},
						},
					},
				},
				ImportedPipelines: &enhancers.ImportedPipeline{
					Pipeline: &models.Pipeline{
						Name: utils.GetPtr("test"),
					},
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 1),
					JobName:             "test",
				},
			},
			want: &models.Pipeline{
				Jobs: []*models.Job{
					{
						PreSteps: []*models.Step{
							{
								Imports: &models.Import{
									Source: &models.ImportSource{
										Path: utils.GetPtr("test"),
										Type: models.SourceTypeLocal,
									},
									FileReference: testutils.CreateFileReference(1, 1, 1, 1),
									Pipeline: &models.Pipeline{
										Name: utils.GetPtr("test"),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "map pre step envs",
			args: args{
				Pipeline: &models.Pipeline{
					Jobs: []*models.Job{
						{
							PreSteps: []*models.Step{
								{
									EnvironmentVariables: &models.EnvironmentVariablesRef{
										Imports: &models.Import{
											Source: &models.ImportSource{
												Path: utils.GetPtr("test"),
												Type: models.SourceTypeLocal,
											},
											FileReference: testutils.CreateFileReference(1, 1, 1, 1),
										},
									},
								},
							},
						},
					},
				},
				ImportedPipelines: &enhancers.ImportedPipeline{
					Pipeline: &models.Pipeline{
						Name: utils.GetPtr("test"),
					},
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 1),
					JobName:             "test",
				},
			},
			want: &models.Pipeline{
				Jobs: []*models.Job{
					{
						PreSteps: []*models.Step{
							{
								EnvironmentVariables: &models.EnvironmentVariablesRef{
									Imports: &models.Import{
										Source: &models.ImportSource{
											Path: utils.GetPtr("test"),
											Type: models.SourceTypeLocal,
										},
										FileReference: testutils.CreateFileReference(1, 1, 1, 1),
										Pipeline: &models.Pipeline{
											Name: utils.GetPtr("test"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "map post steps",
			args: args{
				Pipeline: &models.Pipeline{
					Jobs: []*models.Job{
						{
							PostSteps: []*models.Step{
								{
									Imports: &models.Import{
										Source: &models.ImportSource{
											Path: utils.GetPtr("test"),
											Type: models.SourceTypeLocal,
										},
										FileReference: testutils.CreateFileReference(1, 1, 1, 1),
									},
								},
							},
						},
					},
				},
				ImportedPipelines: &enhancers.ImportedPipeline{
					Pipeline: &models.Pipeline{
						Name: utils.GetPtr("test"),
					},
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 1),
					JobName:             "test",
				},
			},
			want: &models.Pipeline{
				Jobs: []*models.Job{
					{
						PostSteps: []*models.Step{
							{
								Imports: &models.Import{
									Source: &models.ImportSource{
										Path: utils.GetPtr("test"),
										Type: models.SourceTypeLocal,
									},
									FileReference: testutils.CreateFileReference(1, 1, 1, 1),
									Pipeline: &models.Pipeline{
										Name: utils.GetPtr("test"),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "map post step envs",
			args: args{
				Pipeline: &models.Pipeline{
					Jobs: []*models.Job{
						{
							PostSteps: []*models.Step{
								{
									EnvironmentVariables: &models.EnvironmentVariablesRef{
										Imports: &models.Import{
											Source: &models.ImportSource{
												Path: utils.GetPtr("test"),
												Type: models.SourceTypeLocal,
											},
											FileReference: testutils.CreateFileReference(1, 1, 1, 1),
										},
									},
								},
							},
						},
					},
				},
				ImportedPipelines: &enhancers.ImportedPipeline{
					Pipeline: &models.Pipeline{
						Name: utils.GetPtr("test"),
					},
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 1),
					JobName:             "test",
				},
			},
			want: &models.Pipeline{
				Jobs: []*models.Job{
					{
						PostSteps: []*models.Step{
							{
								EnvironmentVariables: &models.EnvironmentVariablesRef{
									Imports: &models.Import{
										Source: &models.ImportSource{
											Path: utils.GetPtr("test"),
											Type: models.SourceTypeLocal,
										},
										FileReference: testutils.CreateFileReference(1, 1, 1, 1),
										Pipeline: &models.Pipeline{
											Name: utils.GetPtr("test"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergePipelines(tt.args.Pipeline, tt.args.ImportedPipelines)
			assert.Equal(t, tt.want, got)
		})
	}
}
