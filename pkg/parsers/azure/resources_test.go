package azure

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestParseResource(t *testing.T) {
	testCases := []struct {
		name              string
		resources         *azureModels.Resources
		expectedResources *models.Resources
	}{
		{
			name:              "Resources is nil",
			resources:         nil,
			expectedResources: nil,
		},
		{
			name: "Resources is nil",
			resources: &azureModels.Resources{
				Resources: nil,
			},
			expectedResources: nil,
		},
		{
			name: "Resources repositories is nil",
			resources: &azureModels.Resources{
				Resources: []*azureModels.Resource{
					{
						Repositories: nil,
					},
				},
			},
			expectedResources: nil,
		},
		{
			name: "Resources repositories is the local repository",
			resources: &azureModels.Resources{
				Resources: []*azureModels.Resource{
					{
						Repositories: []*azureModels.RepositoryRef{
							{
								Repository: &azureModels.Repository{
									Repository: "test",
									Name:       "org/repo",
								},
								FileReference: testutils.CreateFileReference(43, 3, 45, 13),
							},
						},
						FileReference: testutils.CreateFileReference(43, 3, 45, 13),
					},
				},
			},
			expectedResources: &models.Resources{
				Repositories: []*models.ImportSource{
					{
						RepositoryAlias: utils.GetPtr("test"),
						Reference:       utils.GetPtr(""),
						Type:            models.SourceTypeLocal,
						SCM:             consts.AzurePlatform,
						Repository:      utils.GetPtr("org/repo"),
					},
				},
				FileReference: testutils.CreateFileReference(43, 3, 45, 13),
			},
		},
		{
			name: "Resources repositories is github repository",
			resources: &azureModels.Resources{
				Resources: []*azureModels.Resource{
					{
						Repositories: []*azureModels.RepositoryRef{
							{
								Repository: &azureModels.Repository{
									Repository: "test",
									Ref:        "ref",
									Type:       "github",
									Name:       "org/repo",
								},
							},
						},
						FileReference: testutils.CreateFileReference(43, 3, 45, 13),
					},
				},
				FileReference: testutils.CreateFileReference(43, 3, 45, 13),
			},
			expectedResources: &models.Resources{
				Repositories: []*models.ImportSource{
					{
						RepositoryAlias: utils.GetPtr("test"),
						Reference:       utils.GetPtr("ref"),
						Type:            models.SourceTypeRemote,
						SCM:             consts.GitHubPlatform,
						Repository:      utils.GetPtr("org/repo"),
					},
				},
				FileReference: testutils.CreateFileReference(43, 3, 45, 13),
			},
		},
		{
			name: "Resources repositories is another azure repository",
			resources: &azureModels.Resources{
				Resources: []*azureModels.Resource{
					{
						Repositories: []*azureModels.RepositoryRef{
							{
								Repository: &azureModels.Repository{
									Repository: "test",
									Ref:        "ref",
									Type:       "git",
									Name:       "org/repo",
								},
							},
						},
						FileReference: testutils.CreateFileReference(43, 3, 45, 13),
					},
				},
			},
			expectedResources: &models.Resources{
				Repositories: []*models.ImportSource{
					{
						RepositoryAlias: utils.GetPtr("test"),
						Reference:       utils.GetPtr("ref"),
						Type:            models.SourceTypeRemote,
						SCM:             consts.AzurePlatform,
						Repository:      utils.GetPtr("org/repo"),
					},
				},
				FileReference: testutils.CreateFileReference(43, 3, 45, 13),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseResources(testCase.resources)
			testutils.DeepCompare(t, testCase.expectedResources, got)
		})
	}
}
