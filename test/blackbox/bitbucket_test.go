package blackbox

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestBitbucket(t *testing.T) {
	testCases := []TestCase{
		{
			Filename: "simple-pipeline.yml",
			Expected: &models.Pipeline{
				Defaults: &models.Defaults{
					Runner: &models.Runner{
						DockerMetadata: &models.DockerMetadata{
							Image: utils.GetPtr("atlassian/default-image:3"),
						},
					},
				},
				Jobs: []*models.Job{
					{
						FileReference: testutils.CreateFileReference(14, 11, 27, 36),
						ID:            utils.GetPtr("job-default"),
						Name:          utils.GetPtr("default"),
						Metadata: models.Metadata{
							Build: true,
							Test:  true,
						},
						Steps: []*models.Step{
							{
								Type: "shell",
								Name: utils.GetPtr("Build and Test"),
								Metadata: models.Metadata{
									Build: true,
									Test:  true,
								},
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("IMAGE_NAME=$BITBUCKET_REPO_SLUG\ndocker build . --file Dockerfile --tag ${IMAGE_NAME}"),
									FileReference: testutils.CreateFileReference(17, 17, 18, 69),
								},
								FileReference: testutils.CreateFileReference(14, 11, 22, 23),
							},
							{
								Name: utils.GetPtr("Lint the Dockerfile"),
								Type: "shell",
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("hadolint Dockerfile"),
									FileReference: testutils.CreateFileReference(27, 17, 27, 36),
								},
								Runner: &models.Runner{
									DockerMetadata: &models.DockerMetadata{
										Image: utils.GetPtr("hadolint/hadolint:latest-debian"),
									},
								},
								FileReference: testutils.CreateFileReference(23, 11, 27, 36),
							},
						},
					},
					{
						FileReference: testutils.CreateFileReference(30, 9, 54, 21),
						ID:            utils.GetPtr("job-master"),
						Name:          utils.GetPtr("master"),
						Metadata: models.Metadata{
							Build: true,
							Test:  true,
						},
						Steps: []*models.Step{
							{
								Type: "shell",
								Name: utils.GetPtr("Build and Test"),
								Metadata: models.Metadata{
									Build: true,
									Test:  true,
								},
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("IMAGE_NAME=$BITBUCKET_REPO_SLUG\ndocker build . --file Dockerfile --tag ${IMAGE_NAME}\ndocker save ${IMAGE_NAME} --output \"${IMAGE_NAME}.tar\""),
									FileReference: testutils.CreateFileReference(33, 15, 35, 69),
								},
								FileReference: testutils.CreateFileReference(30, 9, 41, 20),
							},
							{
								Type: "shell",
								Name: utils.GetPtr("Deploy to Production"),
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("echo ${DOCKERHUB_PASSWORD} | docker login --username \"$DOCKERHUB_USERNAME\" --password-stdin\nIMAGE_NAME=$BITBUCKET_REPO_SLUG\ndocker load --input \"${IMAGE_NAME}.tar\"\nVERSION=\"prod-0.1.${BITBUCKET_BUILD_NUMBER}\"\nIMAGE=${DOCKERHUB_NAMESPACE}/${IMAGE_NAME}\ndocker tag \"${IMAGE_NAME}\" \"${IMAGE}:${VERSION}\"\ndocker push \"${IMAGE}:${VERSION}\""),
									FileReference: testutils.CreateFileReference(46, 15, 52, 48),
								},
								FileReference: testutils.CreateFileReference(42, 9, 54, 21),
							},
						},
					},
				},
			},
		},
		{
			Filename: "alias-pipeline.yml",
			Expected: &models.Pipeline{
				Defaults: &models.Defaults{
					Runner: &models.Runner{
						DockerMetadata: &models.DockerMetadata{
							Image: utils.GetPtr("node:14.17.6"),
						},
					},
				},
				Jobs: []*models.Job{
					{
						FileReference: testutils.CreateFileReference(10, 7, 23, 61),
						ID:            utils.GetPtr("job-**"),
						Name:          utils.GetPtr("**"),
						Metadata: models.Metadata{
							Build: true,
						},
						Steps: []*models.Step{
							{
								Type: "shell",
								Name: utils.GetPtr("build"),
								Metadata: models.Metadata{
									Build: true,
								},
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("yarn\nyarn build"),
									FileReference: testutils.CreateFileReference(17, 13, 18, 23),
								},
								AfterScript: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("npx notify -s \"Install and build\" --only-failure"),
									FileReference: testutils.CreateFileReference(23, 13, 23, 61),
								},
								FileReference: testutils.CreateAliasFileReference(10, 7, 23, 61, true),
							},
						},
					},
					{
						FileReference: testutils.CreateFileReference(10, 7, 27, 24),
						ID:            utils.GetPtr("job-deploy-staging"),
						Name:          utils.GetPtr("deploy-staging"),
						Metadata: models.Metadata{
							Build: true,
						},
						Steps: []*models.Step{
							{
								Type: "shell",
								Name: utils.GetPtr("build"),
								Metadata: models.Metadata{
									Build: true,
								},
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("yarn\nyarn build"),
									FileReference: testutils.CreateFileReference(17, 13, 18, 23),
								},
								AfterScript: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("npx notify -s \"Install and build\" --only-failure"),
									FileReference: testutils.CreateFileReference(23, 13, 23, 61),
								},
								FileReference: testutils.CreateAliasFileReference(10, 7, 23, 61, true),
							},
							{
								Name: utils.GetPtr("deploy"),
								Type: "shell",
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("echo deploy"),
									FileReference: testutils.CreateFileReference(27, 13, 27, 24),
								},
								FileReference: testutils.CreateAliasFileReference(24, 7, 27, 24, true),
							},
						},
					},
					{
						FileReference: testutils.CreateFileReference(10, 7, 27, 24),
						ID:            utils.GetPtr("job-master"),
						Name:          utils.GetPtr("master"),
						Metadata: models.Metadata{
							Build: true,
						},
						Steps: []*models.Step{
							{
								Type: "shell",
								Name: utils.GetPtr("build"),
								Metadata: models.Metadata{
									Build: true,
								},
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("yarn\nyarn build"),
									FileReference: testutils.CreateFileReference(17, 13, 18, 23),
								},
								AfterScript: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("npx notify -s \"Install and build\" --only-failure"),
									FileReference: testutils.CreateFileReference(23, 13, 23, 61),
								},
								FileReference: testutils.CreateAliasFileReference(10, 7, 23, 61, true),
							},
							{
								Name: utils.GetPtr("deploy"),
								Type: "shell",
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("echo deploy"),
									FileReference: testutils.CreateFileReference(27, 13, 27, 24),
								},
								FileReference: testutils.CreateAliasFileReference(24, 7, 27, 24, true),
							},
						},
					},
				},
			},
		},
		{
			Filename: "merge-step-pipeline.yml",
			Expected: &models.Pipeline{
				Defaults: &models.Defaults{
					Runner: &models.Runner{
						DockerMetadata: &models.DockerMetadata{
							Image: utils.GetPtr("atlassian/default-image:3"),
						},
					},
				},
				Jobs: []*models.Job{
					{
						FileReference: testutils.CreateFileReference(19, 13, 23, 30),
						ID:            utils.GetPtr("job-main"),
						Name:          utils.GetPtr("main"),
						Metadata: models.Metadata{
							Test: true,
						},
						Steps: []*models.Step{
							{
								Type: "shell",
								Name: utils.GetPtr("Test"),
								Metadata: models.Metadata{
									Test: true,
								},
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("echo testing...\nnpm run test"),
									FileReference: testutils.CreateFileReference(8, 13, 9, 25),
								},
								FileReference: testutils.CreateAliasFileReference(19, 13, 21, 25, true),
							},
							{
								Type: "shell",
								Name: utils.GetPtr("Send Result"),
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("echo sending result...\nnpm run send-result"),
									FileReference: testutils.CreateFileReference(13, 13, 14, 32),
								},
								FileReference: testutils.CreateAliasFileReference(22, 13, 23, 30, true),
							},
						},
					},
				},
			},
		},
		{
			Filename: "variables-pipeline.yml",
			Expected: &models.Pipeline{
				Defaults: &models.Defaults{
					Runner: &models.Runner{
						DockerMetadata: &models.DockerMetadata{
							Image: utils.GetPtr("atlassian/default-image:3"),
						},
					},
				},
				Jobs: []*models.Job{
					{
						FileReference: testutils.CreateFileReference(6, 9, 24, 25),
						ID:            utils.GetPtr("job-master"),
						Name:          utils.GetPtr("master"),
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("Deploy to Production"),
								Type: "task",
								Task: &models.Task{
									Name:        utils.GetPtr("atlassian/aws-elasticbeanstalk-deploy:1.0.2\natlassian/aws-elasticbeanstalk-run:1.0.2"),
									VersionType: "none",
								},
								EnvironmentVariables: &models.EnvironmentVariablesRef{
									EnvironmentVariables: models.EnvironmentVariables{
										"AWS_ACCESS_KEY_ID":     "$AWS_ACCESS_KEY_ID",
										"AWS_SECRET_ACCESS_KEY": "$AWS_SECRET_ACCESS_KEY",
										"AWS_DEFAULT_REGION":    "$AWS_DEFAULT_REGION",
										"APPLICATION_NAME":      "pipes-templates-java-spring-boot-app",
										"ENVIRONMENT_NAME":      "Production",
										"S3_BUCKET":             "pipes-template-java-spring-boot-source",
										"ZIP_FILE":              "application.zip",
										"VERSION_LABEL":         "prod-0.1.$BITBUCKET_BUILD_NUMBER",
										"KEY":                   "value",
										"FOO":                   "bar",
									},
									FileReference: testutils.CreateFileReference(12, 17, 24, 25),
								},
								FileReference: testutils.CreateFileReference(6, 9, 24, 25),
							},
						},
					},
				},
			},
		},
		{
			Filename: "image-step.yml",
			Expected: &models.Pipeline{
				Defaults: &models.Defaults{},
				Jobs: []*models.Job{
					{
						FileReference: testutils.CreateFileReference(4, 11, 8, 71),
						ID:            utils.GetPtr("job-master"),
						Name:          utils.GetPtr("master"),
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("Run Aqua scanner"),
								Type: "shell",
								Shell: &models.Shell{
									Type:          utils.GetPtr("shell"),
									Script:        utils.GetPtr("trivy fs --security-checks config,vuln,secret --sast ."),
									FileReference: testutils.CreateFileReference(8, 17, 8, 71),
								},
								Runner: &models.Runner{
									DockerMetadata: &models.DockerMetadata{
										Image: utils.GetPtr("aquasec/aqua-scanner"),
									},
								},
								FileReference: testutils.CreateFileReference(4, 11, 8, 71),
							},
						},
					},
				},
			},
		},
	}

	executeTestCases(t, testCases, "bitbucket", consts.BitbucketPlatform, "", "")
}
