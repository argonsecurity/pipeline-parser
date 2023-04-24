package azure

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
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

func Test_InheritParentPipelineData(t *testing.T) {
	type args struct {
		Parent *models.Pipeline
		Child  *models.Pipeline
	}
	tests := []struct {
		name string
		args args
		want *models.Pipeline
	}{
		{
			name: "parent is nil - return child",
			args: args{
				Parent: nil,
				Child: &models.Pipeline{
					Name: utils.GetPtr("child"),
					Defaults: &models.Defaults{
						Resources: &models.Resources{
							Repositories: []*models.ImportSource{
								{
									Path:            utils.GetPtr("path"),
									Type:            models.SourceTypeLocal,
									SCM:             consts.AzurePlatform,
									Organization:    utils.GetPtr("org"),
									Repository:      utils.GetPtr("repo"),
									RepositoryAlias: utils.GetPtr("repoAlias"),
									Reference:       utils.GetPtr("ref"),
								},
							},
						},
					},
				},
			},
			want: &models.Pipeline{
				Name: utils.GetPtr("child"),
				Defaults: &models.Defaults{
					Resources: &models.Resources{
						Repositories: []*models.ImportSource{
							{
								Path:            utils.GetPtr("path"),
								Type:            models.SourceTypeLocal,
								SCM:             consts.AzurePlatform,
								Organization:    utils.GetPtr("org"),
								Repository:      utils.GetPtr("repo"),
								RepositoryAlias: utils.GetPtr("repoAlias"),
								Reference:       utils.GetPtr("ref"),
							},
						},
					},
				},
			},
		},
		{
			name: "child is nil - return nil",
			args: args{
				Parent: nil,
				Child:  nil,
			},
			want: nil,
		},
		{
			name: "child includes all parent repositories",
			args: args{
				Parent: &models.Pipeline{
					Defaults: &models.Defaults{
						Resources: &models.Resources{
							Repositories: []*models.ImportSource{
								{
									Path:            utils.GetPtr("path"),
									Type:            models.SourceTypeLocal,
									SCM:             consts.AzurePlatform,
									Organization:    utils.GetPtr("org"),
									Repository:      utils.GetPtr("repo"),
									RepositoryAlias: utils.GetPtr("repoAlias"),
									Reference:       utils.GetPtr("ref"),
								},
							},
						},
					},
				},
				Child: &models.Pipeline{
					Defaults: &models.Defaults{
						Resources: &models.Resources{
							Repositories: []*models.ImportSource{
								{
									Path:            utils.GetPtr("path"),
									Type:            models.SourceTypeLocal,
									SCM:             consts.AzurePlatform,
									Organization:    utils.GetPtr("org"),
									Repository:      utils.GetPtr("repo"),
									RepositoryAlias: utils.GetPtr("repoAlias"),
									Reference:       utils.GetPtr("ref"),
								},
							},
						},
					},
				},
			},
			want: &models.Pipeline{
				Defaults: &models.Defaults{
					Resources: &models.Resources{
						Repositories: []*models.ImportSource{
							{
								Path:            utils.GetPtr("path"),
								Type:            models.SourceTypeLocal,
								SCM:             consts.AzurePlatform,
								Organization:    utils.GetPtr("org"),
								Repository:      utils.GetPtr("repo"),
								RepositoryAlias: utils.GetPtr("repoAlias"),
								Reference:       utils.GetPtr("ref"),
							},
						},
					},
				},
			},
		},
		{
			name: "child inherits some parents repositories",
			args: args{
				Parent: &models.Pipeline{
					Defaults: &models.Defaults{
						Resources: &models.Resources{
							Repositories: []*models.ImportSource{
								{
									Path:            utils.GetPtr("path"),
									Type:            models.SourceTypeLocal,
									SCM:             consts.AzurePlatform,
									Organization:    utils.GetPtr("org"),
									Repository:      utils.GetPtr("repo"),
									RepositoryAlias: utils.GetPtr("repoAlias"),
									Reference:       utils.GetPtr("ref"),
								},
								{
									Path:            utils.GetPtr("path"),
									Type:            models.SourceTypeLocal,
									SCM:             consts.AzurePlatform,
									Organization:    utils.GetPtr("org"),
									Repository:      utils.GetPtr("repo2"),
									RepositoryAlias: utils.GetPtr("repoAlias"),
									Reference:       utils.GetPtr("ref"),
								},
							},
						},
					},
				},
				Child: &models.Pipeline{
					Defaults: &models.Defaults{
						Resources: &models.Resources{
							Repositories: []*models.ImportSource{
								{
									Path:            utils.GetPtr("path"),
									Type:            models.SourceTypeLocal,
									SCM:             consts.AzurePlatform,
									Organization:    utils.GetPtr("org"),
									Repository:      utils.GetPtr("repo"),
									RepositoryAlias: utils.GetPtr("repoAlias"),
									Reference:       utils.GetPtr("ref"),
								},
							},
						},
					},
				},
			},
			want: &models.Pipeline{
				Defaults: &models.Defaults{
					Resources: &models.Resources{
						Repositories: []*models.ImportSource{
							{
								Path:            utils.GetPtr("path"),
								Type:            models.SourceTypeLocal,
								SCM:             consts.AzurePlatform,
								Organization:    utils.GetPtr("org"),
								Repository:      utils.GetPtr("repo"),
								RepositoryAlias: utils.GetPtr("repoAlias"),
								Reference:       utils.GetPtr("ref"),
							},
							{
								Path:            utils.GetPtr("path"),
								Type:            models.SourceTypeLocal,
								SCM:             consts.AzurePlatform,
								Organization:    utils.GetPtr("org"),
								Repository:      utils.GetPtr("repo2"),
								RepositoryAlias: utils.GetPtr("repoAlias"),
								Reference:       utils.GetPtr("ref"),
							},
						},
					},
				},
			},
		},
		{
			name: "child inherits parents repositories but does not define Defaults",
			args: args{
				Parent: &models.Pipeline{
					Defaults: &models.Defaults{
						Resources: &models.Resources{
							Repositories: []*models.ImportSource{
								{
									Path:            utils.GetPtr("path"),
									Type:            models.SourceTypeLocal,
									SCM:             consts.AzurePlatform,
									Organization:    utils.GetPtr("org"),
									Repository:      utils.GetPtr("repo"),
									RepositoryAlias: utils.GetPtr("repoAlias"),
									Reference:       utils.GetPtr("ref"),
								},
								{
									Path:            utils.GetPtr("path"),
									Type:            models.SourceTypeLocal,
									SCM:             consts.AzurePlatform,
									Organization:    utils.GetPtr("org"),
									Repository:      utils.GetPtr("repo2"),
									RepositoryAlias: utils.GetPtr("repoAlias"),
									Reference:       utils.GetPtr("ref"),
								},
							},
						},
					},
				},
				Child: &models.Pipeline{
					Defaults: nil,
				},
			},
			want: &models.Pipeline{
				Defaults: &models.Defaults{
					Resources: &models.Resources{
						Repositories: []*models.ImportSource{
							{
								Path:            utils.GetPtr("path"),
								Type:            models.SourceTypeLocal,
								SCM:             consts.AzurePlatform,
								Organization:    utils.GetPtr("org"),
								Repository:      utils.GetPtr("repo"),
								RepositoryAlias: utils.GetPtr("repoAlias"),
								Reference:       utils.GetPtr("ref"),
							},
							{
								Path:            utils.GetPtr("path"),
								Type:            models.SourceTypeLocal,
								SCM:             consts.AzurePlatform,
								Organization:    utils.GetPtr("org"),
								Repository:      utils.GetPtr("repo2"),
								RepositoryAlias: utils.GetPtr("repoAlias"),
								Reference:       utils.GetPtr("ref"),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var azureEnhancer AzureEnhancer
			got := azureEnhancer.InheritParentPipelineData(tt.args.Parent, tt.args.Child)
			assert.Equal(t, tt.want, got)
		})
	}
}
