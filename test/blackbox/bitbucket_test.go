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
								FileReference: testutils.CreateFileReference(15, 13, 22, 23),
							},
							{
								Name: utils.GetPtr("Lint the Dockerfile"),
								Shell: &models.Shell{
									Script:        utils.GetPtr("- hadolint Dockerfile \n"),
									FileReference: testutils.CreateFileReference(27, 17, 27, 36),
								},
								FileReference: testutils.CreateFileReference(24, 13, 27, 36),
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
	}

	executeTestCases(t, testCases, "bitbucket", consts.BitbucketPlatform)
}
