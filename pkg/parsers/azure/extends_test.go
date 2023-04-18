package azure

import (
	"testing"

	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestParseExtends(t *testing.T) {
	testCases := []struct {
		name            string
		extends         *azureModels.Extends
		expectedImports []*models.Import
	}{
		{
			name:            "Extends is nil",
			extends:         nil,
			expectedImports: nil,
		},
		{
			name: "Extends with relative path",
			extends: &azureModels.Extends{
				Template: azureModels.Template{
					Template: "template1",
				},
			},
			expectedImports: []*models.Import{{
				Source: &models.ImportSource{
					Path:            utils.GetPtr("template1"),
					RepositoryAlias: utils.GetPtr(""),
				},
			}},
		},
		{
			name: "Extends with repository alias",
			extends: &azureModels.Extends{
				Template: azureModels.Template{
					Template: "template1@repo1",
				},
			},
			expectedImports: []*models.Import{{
				Source: &models.ImportSource{
					Path:            utils.GetPtr("template1"),
					RepositoryAlias: utils.GetPtr("repo1"),
				},
			}},
		},
		{
			name: "Extends with parameter template",
			extends: &azureModels.Extends{
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				Template: azureModels.Template{
					Template: "template1@repo1",
					Parameters: map[string]any{
						"foo": "bar",
						"testSteps": azureModels.Template{
							Template: "template2@repo2",
							Parameters: map[string]any{
								"foo2": "bar2",
							},
						},
					},
				},
			},
			expectedImports: []*models.Import{
				{
					FileReference: testutils.CreateAliasFileReference(1, 2, 3, 4, false),
					Parameters: map[string]any{
						"foo": "bar",
					},
					Source: &models.ImportSource{
						Path:            utils.GetPtr("template1"),
						RepositoryAlias: utils.GetPtr("repo1"),
					},
				},
				{
					Parameters: map[string]any{
						"foo2": "bar2",
					},
					Source: &models.ImportSource{
						Path:            utils.GetPtr("template2"),
						RepositoryAlias: utils.GetPtr("repo2"),
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseExtends(testCase.extends)
			testutils.DeepCompare(t, testCase.expectedImports, got)
		})
	}
}
