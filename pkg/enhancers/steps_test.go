package enhancers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/argonsecurity/pipeline-parser/pkg/enhancers/config"
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
			name: "Step name doesn't contain build",
			step: &models.Step{
				Name: utils.GetPtr("some stage"),
			},
			config: config.CommonConfiguration,
			expectedStep: &models.Step{
				Name: utils.GetPtr("some stage"),
				Metadata: models.Metadata{
					Build: false,
				},
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
			name: "Shell no build",
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
				Metadata: models.Metadata{
					Build: false,
				},
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

			assert.Equal(t, testCase.expectedStep, enhancedStep)
		})

	}
}
