package general

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/argonsecurity/pipeline-parser/pkg/enhancers/general/config"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestEnhanceStep(t *testing.T) {
	testCases := []struct {
		name         string
		step         *models.Step
		config       *config.EnhancementConfiguration
		expectedStep *models.Step
	}{
		{
			name: "Step name contains build (lowercase)",
			step: &models.Step{
				Name: utils.GetPtr("build app"),
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Name: utils.GetPtr("build app"),
				Metadata: models.Metadata{
					Build: true,
				},
			},
		},
		{
			name: "Step name contains Build (Uppercase)",
			step: &models.Step{
				Name: utils.GetPtr("Build app"),
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Name: utils.GetPtr("Build app"),
				Metadata: models.Metadata{
					Build: true,
				},
			},
		},
		{
			name: "Step name contains test (lowercase)",
			step: &models.Step{
				Name: utils.GetPtr("test app"),
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Name: utils.GetPtr("test app"),
				Metadata: models.Metadata{
					Test: true,
				},
			},
		},
		{
			name: "Step name contains Test (Uppercase)",
			step: &models.Step{
				Name: utils.GetPtr("Test app"),
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Name: utils.GetPtr("Test app"),
				Metadata: models.Metadata{
					Test: true,
				},
			},
		},
		{
			name: "Step name contains tests (lowercase)",
			step: &models.Step{
				Name: utils.GetPtr("tests app"),
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Name: utils.GetPtr("tests app"),
				Metadata: models.Metadata{
					Test: true,
				},
			},
		},
		{
			name: "Step name contains Tests (Uppercase)",
			step: &models.Step{
				Name: utils.GetPtr("Tests app"),
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Name: utils.GetPtr("Tests app"),
				Metadata: models.Metadata{
					Test: true,
				},
			},
		},
		{
			name: "Step name doesn't contain build test or deploy",
			step: &models.Step{
				Name: utils.GetPtr("some step"),
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Name:     utils.GetPtr("some step"),
				Metadata: models.Metadata{},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			enhancedStep := enhanceStep(testCase.step, testCase.config)

			assert.Equal(t, testCase.expectedStep, enhancedStep)
		})

	}
}

func TestEnhanceShellStep(t *testing.T) {
	testCases := []struct {
		name         string
		step         *models.Step
		config       *config.EnhancementConfiguration
		expectedStep *models.Step
	}{
		{
			name: "Shell docker build",
			step: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("docker build -t ${IMAGE_NAME} ."),
				},
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("docker build -t ${IMAGE_NAME} ."),
				},
				Metadata: models.Metadata{
					Build: true,
				},
			},
		},
		{
			name: "Shell docker-compose build",
			step: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("docker-compose build -t ${IMAGE_NAME} ."),
				},
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("docker-compose build -t ${IMAGE_NAME} ."),
				},
				Metadata: models.Metadata{
					Build: true,
				},
			},
		},
		{
			name: "Shell npm run build",
			step: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("npm run build"),
				},
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("npm run build"),
				},
				Metadata: models.Metadata{
					Build: true,
				},
			},
		},
		{
			name: "Shell yarn run build",
			step: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("yarn run build"),
				},
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("yarn run build"),
				},
				Metadata: models.Metadata{
					Build: true,
				},
			},
		},
		{
			name: "Shell yarn build",
			step: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("yarn build"),
				},
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("yarn build"),
				},
				Metadata: models.Metadata{
					Build: true,
				},
			},
		},
		{
			name: "Shell go build",
			step: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("go build -ldflags \"-s -w -X=main.version=$(VERSION)\" ./cmd/pipeline-parser"),
				},
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("go build -ldflags \"-s -w -X=main.version=$(VERSION)\" ./cmd/pipeline-parser"),
				},
				Metadata: models.Metadata{
					Build: true,
				},
			},
		},
		{
			name: "Shell go test",
			step: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("go test"),
				},
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("go test"),
				},
				Metadata: models.Metadata{
					Test: true,
				},
			},
		},
		{
			name: "Shell npm run test",
			step: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("npm run test"),
				},
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("npm run test"),
				},
				Metadata: models.Metadata{
					Test: true,
				},
			},
		},
		{
			name: "Shell yarn run test",
			step: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("yarn run test"),
				},
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("yarn run test"),
				},
				Metadata: models.Metadata{
					Test: true,
				},
			},
		},
		{
			name: "Shell yarn test",
			step: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("yarn test"),
				},
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("yarn test"),
				},
				Metadata: models.Metadata{
					Test: true,
				},
			},
		},
		{
			name: "Shell no metadata",
			step: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("echo hello"),
				},
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Type: models.ShellStepType,
				Shell: &models.Shell{
					Script: utils.GetPtr("echo hello"),
				},
				Metadata: models.Metadata{},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			enhancedStep := enhanceShellStep(testCase.step, testCase.config)

			assert.Equal(t, testCase.expectedStep, enhancedStep)
		})

	}
}

func TestEnhanceTaskStep(t *testing.T) {
	testCases := []struct {
		name         string
		step         *models.Step
		config       *config.EnhancementConfiguration
		expectedStep *models.Step
	}{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			enhancedStep := enhanceTaskStep(testCase.step, testCase.config)

			assert.Equal(t, testCase.expectedStep, enhancedStep, testCase.name)
		})

	}
}
