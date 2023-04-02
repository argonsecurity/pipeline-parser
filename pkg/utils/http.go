package utils

import (
	"fmt"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/imroc/req/v3"
)

func GetHttpClient(credentials *models.Credentials) *req.Client {
	client := req.C()
	if credentials == nil {
		return client
	}

	return client.SetCommonHeader("Authorization", fmt.Sprintf("token %s", credentials.Token))
}
