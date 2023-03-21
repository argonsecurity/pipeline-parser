package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
)

var (
	defaultRetryWaitMin = 20 * time.Millisecond
	defaultRetryWaitMax = 5 * time.Second
	defaultRetryMax     = 5
)

type Request struct {
	body io.ReadSeeker
	*http.Request
}

type HTTPClient struct {
	client  GoHTTPClient
	headers Headers

	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
	RetryMax     int

	CheckForRetry CheckForRetry
	Backoff       Backoff
}

func NewHTTPClient(headers Headers) *HTTPClient {
	return &HTTPClient{
		client:  http.DefaultClient,
		headers: headers,

		RetryWaitMin: defaultRetryWaitMin,
		RetryWaitMax: defaultRetryWaitMax,
		RetryMax:     defaultRetryMax,

		CheckForRetry: DefaultRetryPolicy,
		Backoff:       DefaultBackoff,
	}
}

func (httpClient *HTTPClient) createHttpRequest(method string, url string, headers Headers, data interface{}) (*Request, error) {
	var err error

	allHeaders := (map[string]string)(httpClient.headers)
	if allHeaders == nil {
		allHeaders = make(map[string]string)
	}
	maps.Copy(allHeaders, headers)

	buf, contentType, err := parseData(data, allHeaders)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		allHeaders[ContentTypeHeader] = contentType
	}

	var request *Request
	if buf != nil {
		request, err = NewRequest(method, url, bytes.NewReader(buf.Bytes()))
	} else {
		request, err = NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	for k, v := range allHeaders {
		request.Header.Add(k, v)
	}

	return request, nil
}

func (httpClient *HTTPClient) sendRequest(req *Request) ([]byte, error) {
	var err error
	for i := 0; i < httpClient.RetryMax; i++ {
		// If the body is not nil, we need to reset the reader to the beginning
		if req.body != nil {
			if _, err = req.body.Seek(0, 0); err != nil {
				return nil, err
			}
		}

		// Attempt the request
		resp, err := httpClient.client.Do(req.Request)

		// Check if we should retry
		checkOk, checkErr := httpClient.CheckForRetry(resp, err)
		if checkErr != nil {
			return nil, checkErr
		}

		// Decide if we should continue
		if !checkOk {
			if err != nil {
				return nil, err
			}

			// Read the body
			body, err := readBody(resp)
			if err != nil {
				return nil, err
			}

			// If the status code is not 2xx, return an error
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return nil, fmt.Errorf("got a response with status code %d, message: %s", resp.StatusCode, string(body))
			}

			return body, nil
		}

		// We're going to retry, consume any response to reuse the connection
		if err == nil {
			readBody(resp)
		}

		wait := httpClient.Backoff(httpClient.RetryWaitMin, httpClient.RetryWaitMax, i)
		time.Sleep(wait)
	}

	if err != nil {
		return nil, errors.Wrap(err, "max retries exceeded")
	}

	return nil, errors.New("max retries exceeded")
}

func (httpClient *HTTPClient) Get(url string, headers Headers, params Params) ([]byte, error) {
	req, err := httpClient.createHttpRequest(http.MethodGet, url, headers, nil)
	if err != nil {
		return nil, err
	}

	// Append query params
	query := req.URL.Query()
	for k, v := range params {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()

	return httpClient.sendRequest(req)
}

func (httpClient *HTTPClient) Post(url string, headers Headers, data interface{}) ([]byte, error) {
	req, err := httpClient.createHttpRequest(http.MethodPost, url, headers, data)
	if err != nil {
		return nil, err
	}

	return httpClient.sendRequest(req)
}

func (httpClient *HTTPClient) Put(url string, headers Headers, data interface{}) ([]byte, error) {
	req, err := httpClient.createHttpRequest(http.MethodPut, url, headers, data)
	if err != nil {
		return nil, err
	}

	return httpClient.sendRequest(req)
}

func (httpClient *HTTPClient) Delete(url string, headers Headers, data interface{}) ([]byte, error) {
	req, err := httpClient.createHttpRequest(http.MethodDelete, url, headers, data)
	if err != nil {
		return nil, err
	}

	return httpClient.sendRequest(req)
}

func NewRequest(method string, url string, body io.ReadSeeker) (*Request, error) {
	var rcBody io.ReadCloser
	if body != nil {
		rcBody = io.NopCloser(body)
	}

	request, err := http.NewRequest(method, url, rcBody)
	if err != nil {
		return nil, err
	}

	return &Request{body, request}, nil
}

func readBody(resp *http.Response) ([]byte, error) {
	if resp == nil || resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
