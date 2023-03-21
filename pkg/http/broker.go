package http

import (
	"net/http"
	"net/url"
	"path"
)

type BrokerRequestConfig struct {
	Headers Headers `json:"headers,omitempty"`
	Params  Params  `json:"params,omitempty"`
}

type BrokerRequestBody struct {
	Url    string               `json:"url"`
	Method string               `json:"method"`
	Data   interface{}          `json:"data,omitempty"`
	Config *BrokerRequestConfig `json:"config,omitempty"`
}

type BrokerHTTPClient struct {
	BrokerClientUrl string
	client          *HTTPClient
}

const brokerClientEndpoint = "proxy"

var (
	brokerRequestHeaders = Headers{
		"Content-Type": "application/json",
	}
)

func NewBrokerHTTPClient(brokerClientUrl string, headers Headers) *BrokerHTTPClient {
	return &BrokerHTTPClient{
		BrokerClientUrl: brokerClientUrl,
		client:          NewHTTPClient(headers),
	}
}

func createBrokerRequestBody(url string, method string, headers Headers, params Params, data interface{}) *BrokerRequestBody {
	return &BrokerRequestBody{
		Url:    url,
		Method: method,
		Data:   data,
		Config: &BrokerRequestConfig{
			Headers: headers,
			Params:  params,
		},
	}
}

func (brokerHttpClient *BrokerHTTPClient) sendRequest(body *BrokerRequestBody) ([]byte, error) {
	urlObj, err := url.Parse(brokerHttpClient.BrokerClientUrl)
	if err != nil {
		return nil, err
	}

	urlObj.Path = path.Join(urlObj.Path, brokerClientEndpoint)
	return brokerHttpClient.client.Post(urlObj.String(), brokerRequestHeaders, body)
}

func (brokerHttpClient *BrokerHTTPClient) Get(url string, headers Headers, params Params) ([]byte, error) {
	body := createBrokerRequestBody(url, http.MethodGet, headers, params, nil)
	return brokerHttpClient.sendRequest(body)
}

func (brokerHttpClient *BrokerHTTPClient) Post(url string, headers Headers, data interface{}) ([]byte, error) {
	body := createBrokerRequestBody(url, http.MethodPost, headers, nil, data)
	return brokerHttpClient.sendRequest(body)
}

func (brokerHttpClient *BrokerHTTPClient) Put(url string, headers Headers, data interface{}) ([]byte, error) {
	body := createBrokerRequestBody(url, http.MethodPut, headers, nil, data)
	return brokerHttpClient.sendRequest(body)
}

func (brokerHttpClient *BrokerHTTPClient) Delete(url string, headers Headers, data interface{}) ([]byte, error) {
	body := createBrokerRequestBody(url, http.MethodDelete, headers, nil, data)
	return brokerHttpClient.sendRequest(body)
}
