package http

import "net/http"

func GetHTTPClient(brokerClientUrl string, headers Headers) HTTPService {
	if brokerClientUrl != "" {
		return NewBrokerHTTPClient(brokerClientUrl, headers)
	}
	return NewHTTPClient(headers)
}

type GoHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
