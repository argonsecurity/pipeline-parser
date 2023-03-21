package httpmock

import (
	"errors"

	"github.com/argonsecurity/pipeline-parser/pkg/http"
)

type MockHttpClient struct {
	responses map[string]string
}

// Setters
func (m *MockHttpClient) SetResponse(req string, res string) *MockHttpClient {
	if m.responses == nil {
		m.responses = make(map[string]string)
	}
	m.responses[req] = res
	return m
}

// Implementations
func (m *MockHttpClient) Get(url string, headers http.Headers, params http.Params) ([]byte, error) {
	if m.responses[url] == "" {
		return nil, errors.New("not found")
	}
	return []byte(m.responses[url]), nil
}

func (m *MockHttpClient) Post(url string, headers http.Headers, data interface{}) ([]byte, error) {
	if m.responses[url] == "" {
		return nil, errors.New("not found")
	}
	return []byte(m.responses[url]), nil
}

func (m *MockHttpClient) Put(url string, headers http.Headers, data interface{}) ([]byte, error) {
	if m.responses[url] == "" {
		return nil, errors.New("not found")
	}
	return []byte(m.responses[url]), nil
}

func (m *MockHttpClient) Delete(url string, headers http.Headers, data interface{}) ([]byte, error) {
	if m.responses[url] == "" {
		return nil, errors.New("not found")
	}
	return []byte(m.responses[url]), nil
}

func (m *MockHttpClient) GetHTTPClient(url string, headers http.Headers) http.HTTPService {
	return m
}
