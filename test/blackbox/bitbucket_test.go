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
						ID:   utils.GetPtr("job-default"),
						Name: utils.GetPtr("default"),
						Metadata: models.Metadata{
							Build: true,
							Test:  true,
						},
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("Build and Test"),
								Metadata: models.Metadata{
									Build: true,
									Test:  true,
								},
								Shell: &models.Shell{
									Script:        utils.GetPtr("- IMAGE_NAME=$BITBUCKET_REPO_SLUG \n- docker build . --file Dockerfile --tag ${IMAGE_NAME} \n"),
									FileReference: testutils.CreateFileReference(17, 17, 18, 69),
								},
								FileReference: testutils.CreateFileReference(14, 11, 22, 23),
							},
							{
								Name: utils.GetPtr("Lint the Dockerfile"),
								Shell: &models.Shell{
									Script:        utils.GetPtr("- hadolint Dockerfile \n"),
									FileReference: testutils.CreateFileReference(27, 17, 27, 36),
								},
								FileReference: testutils.CreateFileReference(23, 11, 27, 36),
							},
						},
					},
					{
						ID:   utils.GetPtr("job-master"),
						Name: utils.GetPtr("master"),
						Metadata: models.Metadata{
							Build: true,
							Test:  true,
						},
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("Build and Test"),
								Metadata: models.Metadata{
									Build: true,
									Test:  true,
								},
								Shell: &models.Shell{
									Script:        utils.GetPtr("- IMAGE_NAME=$BITBUCKET_REPO_SLUG \n- docker build . --file Dockerfile --tag ${IMAGE_NAME} \n- docker save ${IMAGE_NAME} --output \"${IMAGE_NAME}.tar\" \n"),
									FileReference: testutils.CreateFileReference(33, 15, 35, 69),
								},
								FileReference: testutils.CreateFileReference(30, 9, 41, 20),
							},
							{
								Name: utils.GetPtr("Deploy to Production"),
								Shell: &models.Shell{
									Script:        utils.GetPtr("- echo ${DOCKERHUB_PASSWORD} | docker login --username \"$DOCKERHUB_USERNAME\" --password-stdin \n- IMAGE_NAME=$BITBUCKET_REPO_SLUG \n- docker load --input \"${IMAGE_NAME}.tar\" \n- VERSION=\"prod-0.1.${BITBUCKET_BUILD_NUMBER}\" \n- IMAGE=${DOCKERHUB_NAMESPACE}/${IMAGE_NAME} \n- docker tag \"${IMAGE_NAME}\" \"${IMAGE}:${VERSION}\" \n- docker push \"${IMAGE}:${VERSION}\" \n"),
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
						ID:   utils.GetPtr("job-**"),
						Name: utils.GetPtr("**"),
						Metadata: models.Metadata{
							Build: true,
						},
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("build"),
								Metadata: models.Metadata{
									Build: true,
								},
								Shell: &models.Shell{
									Script:        utils.GetPtr("- yarn \n- yarn build \n"),
									FileReference: testutils.CreateFileReference(17, 13, 18, 23),
								},
								AfterScript: &models.Shell{
									Script:        utils.GetPtr("- npx notify -s \"Install and build\" --only-failure \n"),
									FileReference: testutils.CreateFileReference(23, 13, 23, 61),
								},
								FileReference: testutils.CreateAliasFileReference(10, 7, 23, 61, true),
							},
						},
					},
					{
						ID:   utils.GetPtr("job-deploy-staging"),
						Name: utils.GetPtr("deploy-staging"),
						Metadata: models.Metadata{
							Build: true,
						},
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("build"),
								Metadata: models.Metadata{
									Build: true,
								},
								Shell: &models.Shell{
									Script:        utils.GetPtr("- yarn \n- yarn build \n"),
									FileReference: testutils.CreateFileReference(17, 13, 18, 23),
								},
								AfterScript: &models.Shell{
									Script:        utils.GetPtr("- npx notify -s \"Install and build\" --only-failure \n"),
									FileReference: testutils.CreateFileReference(23, 13, 23, 61),
								},
								FileReference: testutils.CreateAliasFileReference(10, 7, 23, 61, true),
							},
							{
								Name: utils.GetPtr("deploy"),
								Shell: &models.Shell{
									Script:        utils.GetPtr("- echo deploy \n"),
									FileReference: testutils.CreateFileReference(27, 13, 27, 24),
								},
								FileReference: testutils.CreateAliasFileReference(24, 7, 27, 24, true),
							},
						},
					},
					{
						ID:   utils.GetPtr("job-master"),
						Name: utils.GetPtr("master"),
						Metadata: models.Metadata{
							Build: true,
						},
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("build"),
								Metadata: models.Metadata{
									Build: true,
								},
								Shell: &models.Shell{
									Script:        utils.GetPtr("- yarn \n- yarn build \n"),
									FileReference: testutils.CreateFileReference(17, 13, 18, 23),
								},
								AfterScript: &models.Shell{
									Script:        utils.GetPtr("- npx notify -s \"Install and build\" --only-failure \n"),
									FileReference: testutils.CreateFileReference(23, 13, 23, 61),
								},
								FileReference: testutils.CreateAliasFileReference(10, 7, 23, 61, true),
							},
							{
								Name: utils.GetPtr("deploy"),
								Shell: &models.Shell{
									Script:        utils.GetPtr("- echo deploy \n"),
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
						ID:   utils.GetPtr("job-main"),
						Name: utils.GetPtr("main"),
						Metadata: models.Metadata{
							Test: true,
						},
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("Test"),
								Metadata: models.Metadata{
									Test: true,
								},
								Shell: &models.Shell{
									Script:        utils.GetPtr("- echo testing... \n- npm run test \n"),
									FileReference: testutils.CreateFileReference(8, 13, 9, 25),
								},
								FileReference: testutils.CreateAliasFileReference(19, 13, 21, 25, true),
							},
							{
								Name: utils.GetPtr("Send Result"),
								Shell: &models.Shell{
									Script:        utils.GetPtr("- echo sending result... \n- npm run send-result \n"),
									FileReference: testutils.CreateFileReference(13, 13, 14, 32),
								},
								FileReference: testutils.CreateAliasFileReference(22, 13, 23, 30, true),
							},
						},
					},
				},
			},
		},
	}

	executeTestCases(t, testCases, "bitbucket", consts.BitbucketPlatform)
}
