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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseExtends(testCase.extends)
			testutils.DeepCompare(t, testCase.expectedImports, got)
		})
	}
}
