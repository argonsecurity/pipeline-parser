package utils

import (
	"fmt"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/imroc/req/v3"
)

func GetHttpClient(credentials *models.Credentials) *req.Client {
	return createClient(credentials, "token")
}

func GetHttpClientWithBasicAuth(credentials *models.Credentials) *req.Client {
	return createClient(credentials, "basic")
}

func createClient(credentials *models.Credentials, auth string) *req.Client {
	client := req.C()
	if credentials == nil {
		return client
	}

	return client.SetCommonHeader("Authorization", fmt.Sprintf("%s %s", auth, credentials.Token))
}
