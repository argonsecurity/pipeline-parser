package azure

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_loadLocalFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "happy flow - relative path",
			args: args{
				path: "../azure/testdata/file",
			},
			want:    []byte("file content"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadLocalFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadLocalFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_generateRequestUrl(t *testing.T) {
	type args struct {
		proj         string
		repo         string
		path         string
		version      string
		organization string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "happy flow - w/o version",
			args: args{
				proj:         "proj",
				repo:         "repo",
				path:         "file",
				version:      "",
				organization: "azure-org",
			},
			want: "https://dev.azure.com/azure-org/proj/_apis/git/repositories/repo/items?path=file",
		},
		{
			name: "happy flow - w/ string",
			args: args{
				proj:         "proj",
				repo:         "repo",
				path:         "file",
				version:      "refs/head/3",
				organization: "azure-org",
			},
			want: "https://dev.azure.com/azure-org/proj/_apis/git/repositories/repo/items?path=file&versionDescriptor.versionType=tag&version=refs/head/3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateRequestUrl(tt.args.proj, tt.args.repo, tt.args.path, tt.args.version, tt.args.organization, "https://dev.azure.com/")
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_extractRemoteParams(t *testing.T) {
	type args struct {
		jobImport *models.Import
		resources *models.Resources
	}
	tests := []struct {
		name         string
		args         args
		wantProj     string
		wantRepo     string
		wantPath     string
		wantVersion  string
		wantPlatform models.Platform
	}{
		{
			name: "happy flow",
			args: args{
				jobImport: &models.Import{
					Source: &models.ImportSource{
						Path:            utils.GetPtr("file"),
						SCM:             consts.AzurePlatform,
						RepositoryAlias: utils.GetPtr("templates"),
					},
				},
				resources: &models.Resources{
					Repositories: []*models.ImportSource{
						{
							RepositoryAlias: utils.GetPtr("templates"),
							Repository:      utils.GetPtr("proj/repo"),
							Reference:       utils.GetPtr("refs/head/3"),
							SCM:             consts.AzurePlatform,
						},
					},
				},
			},
			wantProj:     "proj",
			wantRepo:     "repo",
			wantPath:     "file",
			wantVersion:  "refs/head/3",
			wantPlatform: consts.AzurePlatform,
		},
		{
			name: "happy flow - multiple repos",
			args: args{
				jobImport: &models.Import{
					Source: &models.ImportSource{
						Path:            utils.GetPtr("file"),
						SCM:             consts.AzurePlatform,
						RepositoryAlias: utils.GetPtr("templates"),
					},
				},
				resources: &models.Resources{
					Repositories: []*models.ImportSource{
						{
							RepositoryAlias: utils.GetPtr("templates2"),
							Repository:      utils.GetPtr("proj/repo"),
							Reference:       utils.GetPtr("refs/head/3"),
							SCM:             consts.GitHubPlatform,
						},
						{
							RepositoryAlias: utils.GetPtr("templates"),
							Repository:      utils.GetPtr("proj/repo"),
							Reference:       utils.GetPtr("refs/head/3"),
							SCM:             consts.AzurePlatform,
						},
					},
				},
			},
			wantProj:     "proj",
			wantRepo:     "repo",
			wantPath:     "file",
			wantVersion:  "refs/head/3",
			wantPlatform: consts.AzurePlatform,
		},
		{
			name: "no resources",
			args: args{
				jobImport: &models.Import{
					Source: &models.ImportSource{
						Path:            utils.GetPtr("file"),
						SCM:             consts.AzurePlatform,
						RepositoryAlias: utils.GetPtr("templates"),
					},
				},
			},
			wantProj:     "",
			wantRepo:     "",
			wantPath:     "file",
			wantVersion:  "",
			wantPlatform: consts.AzurePlatform,
		},
		{
			name: "no matches",
			args: args{
				jobImport: &models.Import{
					Source: &models.ImportSource{
						Path:            utils.GetPtr("file"),
						SCM:             consts.AzurePlatform,
						RepositoryAlias: utils.GetPtr("templates-2"),
					},
				},
				resources: &models.Resources{
					Repositories: []*models.ImportSource{
						{
							RepositoryAlias: utils.GetPtr("templates"),
							Organization:    utils.GetPtr("proj"),
							Repository:      utils.GetPtr("repo"),
							Reference:       utils.GetPtr("refs/head/3"),
							SCM:             consts.AzurePlatform,
						},
					},
				},
			},
			wantProj:     "",
			wantRepo:     "",
			wantPath:     "file",
			wantVersion:  "",
			wantPlatform: consts.AzurePlatform,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProj, gotRepo, gotPath, gotVersion, gotPlatform := extractRemoteParams(tt.args.jobImport, tt.args.resources)
			assert.Equal(t, tt.wantProj, gotProj)
			assert.Equal(t, tt.wantRepo, gotRepo)
			assert.Equal(t, tt.wantPath, gotPath)
			assert.Equal(t, tt.wantVersion, gotVersion)
			assert.Equal(t, tt.wantPlatform, gotPlatform)
		})
	}
}

func Test_getTemplates(t *testing.T) {
	type args struct {
		pipeline    *models.Pipeline
		credentials *models.Credentials
	}
	tests := []struct {
		name    string
		args    args
		want    []*enhancers.ImportedPipeline
		wantErr bool
	}{
		{
			name: "happy flow",
			args: args{
				pipeline: &models.Pipeline{
					Defaults: &models.Defaults{
						Resources: &models.Resources{
							Repositories: []*models.ImportSource{
								{
									RepositoryAlias: utils.GetPtr("templates"),
									Repository:      utils.GetPtr("proj/repo"),
									SCM:             consts.AzurePlatform,
								},
							},
						},
					},
					Imports: []*models.Import{
						{
							FileReference: testutils.CreateFileReference(0, 1, 2, 3),
							Source: &models.ImportSource{
								RepositoryAlias: utils.GetPtr("templates"),
								Path:            utils.GetPtr("file"),
								Type:            models.SourceTypeRemote,
							},
						},
						{
							FileReference: testutils.CreateFileReference(3, 2, 1, 0),
							Source: &models.ImportSource{
								Path: utils.GetPtr("testdata/file"),
								Type: models.SourceTypeLocal,
							},
						},
					},
				},
			},
			want: []*enhancers.ImportedPipeline{
				{
					OriginFileReference: testutils.CreateFileReference(0, 1, 2, 3),
					Data:                []byte("file content"),
					JobName:             "file",
				},
				{
					OriginFileReference: testutils.CreateFileReference(3, 2, 1, 0),
					Data:                []byte("file content"),
					JobName:             "testdata/file",
				},
			},
		},
		{
			name: "local does not exist",
			args: args{
				pipeline: &models.Pipeline{
					Defaults: &models.Defaults{
						Resources: &models.Resources{
							Repositories: []*models.ImportSource{
								{
									RepositoryAlias: utils.GetPtr("templates"),
									Repository:      utils.GetPtr("proj/repo"),
									SCM:             consts.AzurePlatform,
								},
							},
						},
					},
					Imports: []*models.Import{
						{
							Source: &models.ImportSource{
								RepositoryAlias: utils.GetPtr("templates"),
								Path:            utils.GetPtr("file"),
								Type:            models.SourceTypeRemote,
							},
						},
						{
							Source: &models.ImportSource{
								Path: utils.GetPtr("testdata/other"),
								Type: models.SourceTypeLocal,
							},
						},
					},
				},
			},
			want: []*enhancers.ImportedPipeline{
				{
					JobName: "remote",
					Data:    []byte("file content"),
				},
				{
					JobName: "local",
					Data:    nil,
				},
			},
			wantErr: true,
		},
		{
			name: "remote does not exist",
			args: args{
				pipeline: &models.Pipeline{
					Defaults: &models.Defaults{
						Resources: &models.Resources{
							Repositories: []*models.ImportSource{
								{
									RepositoryAlias: utils.GetPtr("templates"),
									Repository:      utils.GetPtr("proj/other"),
									SCM:             consts.AzurePlatform,
								},
							},
						},
					},
					Imports: []*models.Import{
						{
							Source: &models.ImportSource{
								RepositoryAlias: utils.GetPtr("templates"),
								Path:            utils.GetPtr("file"),
								Type:            models.SourceTypeRemote,
							},
						},
						{
							Source: &models.ImportSource{
								Path: utils.GetPtr("testdata/file"),
								Type: models.SourceTypeLocal,
							},
						},
					},
				},
			},
			want: []*enhancers.ImportedPipeline{
				{
					JobName: "remote",
					Data:    nil,
				},
				{
					JobName: "local",
					Data:    []byte("file content"),
				},
			},
			wantErr: true,
		},
		{
			name: "happy flow - all import locations (imports, defaults, jobs, steps, envs)",
			args: args{
				pipeline: &models.Pipeline{
					Defaults: &models.Defaults{
						EnvironmentVariables: &models.EnvironmentVariablesRef{
							Imports: &models.Import{
								Source: &models.ImportSource{
									Path: utils.GetPtr("testdata/file"),
									Type: models.SourceTypeLocal,
								},
								FileReference: testutils.CreateFileReference(1, 1, 1, 1),
							},
						},
					},
					Imports: []*models.Import{
						{
							Source: &models.ImportSource{
								Path: utils.GetPtr("testdata/file"),
								Type: models.SourceTypeLocal,
							},
							FileReference: testutils.CreateFileReference(1, 1, 1, 2),
						},
					},
					Jobs: []*models.Job{
						{
							Imports: &models.Import{
								Source: &models.ImportSource{
									Path: utils.GetPtr("testdata/file"),
									Type: models.SourceTypeLocal,
								},
								FileReference: testutils.CreateFileReference(1, 1, 1, 3),
							},
							EnvironmentVariables: &models.EnvironmentVariablesRef{
								Imports: &models.Import{
									Source: &models.ImportSource{
										Path: utils.GetPtr("testdata/file"),
										Type: models.SourceTypeLocal,
									},
									FileReference: testutils.CreateAliasFileReference(1, 1, 1, 4, false),
								},
							},
							PreSteps: []*models.Step{
								{
									Imports: &models.Import{
										Source: &models.ImportSource{
											Path: utils.GetPtr("testdata/file"),
											Type: models.SourceTypeLocal,
										},
										FileReference: testutils.CreateFileReference(1, 1, 1, 5),
									},
									EnvironmentVariables: &models.EnvironmentVariablesRef{
										Imports: &models.Import{
											Source: &models.ImportSource{
												Path: utils.GetPtr("testdata/file"),
												Type: models.SourceTypeLocal,
											},
											FileReference: testutils.CreateFileReference(1, 1, 1, 6),
										},
									},
								},
							},
							Steps: []*models.Step{
								{
									Imports: &models.Import{
										Source: &models.ImportSource{
											Path: utils.GetPtr("testdata/file"),
											Type: models.SourceTypeLocal,
										},
										FileReference: testutils.CreateFileReference(1, 1, 1, 7),
									},
									EnvironmentVariables: &models.EnvironmentVariablesRef{
										Imports: &models.Import{
											Source: &models.ImportSource{
												Path: utils.GetPtr("testdata/file"),
												Type: models.SourceTypeLocal,
											},
											FileReference: testutils.CreateFileReference(1, 1, 1, 8),
										},
									},
								},
							},
							PostSteps: []*models.Step{
								{
									Imports: &models.Import{
										Source: &models.ImportSource{
											Path: utils.GetPtr("testdata/file"),
											Type: models.SourceTypeLocal,
										},
										FileReference: testutils.CreateFileReference(1, 1, 1, 9),
									},
									EnvironmentVariables: &models.EnvironmentVariablesRef{
										Imports: &models.Import{
											Source: &models.ImportSource{
												Path: utils.GetPtr("testdata/file"),
												Type: models.SourceTypeLocal,
											},
											FileReference: testutils.CreateFileReference(1, 1, 1, 10),
										},
									},
								},
							},
						},
					},
				},
			},
			want: []*enhancers.ImportedPipeline{
				{
					JobName:             "testdata/file",
					Data:                []byte("file content"),
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
				{
					JobName:             "testdata/file",
					Data:                []byte("file content"),
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 2),
				},
				{
					JobName:             "testdata/file",
					Data:                []byte("file content"),
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 3),
				},
				{
					JobName:             "testdata/file",
					Data:                []byte("file content"),
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 4),
				},
				{
					JobName:             "testdata/file",
					Data:                []byte("file content"),
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 5),
				},
				{
					JobName:             "testdata/file",
					Data:                []byte("file content"),
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 6),
				},
				{
					JobName:             "testdata/file",
					Data:                []byte("file content"),
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 7),
				},
				{
					JobName:             "testdata/file",
					Data:                []byte("file content"),
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 8),
				},
				{
					JobName:             "testdata/file",
					Data:                []byte("file content"),
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 9),
				},
				{
					JobName:             "testdata/file",
					Data:                []byte("file content"),
					OriginFileReference: testutils.CreateFileReference(1, 1, 1, 10),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.FileServer(http.Dir("testdata"))
			ts := httptest.NewServer(h)
			defer ts.Close()
			AZURE_SAAS_BASE_URL = ts.URL

			got, err := getTemplates(tt.args.pipeline, tt.args.credentials, utils.GetPtr("azure-org"), &AZURE_SAAS_BASE_URL)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
