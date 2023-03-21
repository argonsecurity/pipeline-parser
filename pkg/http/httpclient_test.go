package http

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testUrl  = "https://www.google.com"
	testData = "data"
)

var (
	testHeaders Headers = Headers{"test": "test"}
	testParams  Params  = Params{"test": "test"}
)

type ClientMock struct{}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK}, nil
}

func GetFakeHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &ClientMock{},

		RetryWaitMin: defaultRetryWaitMin,
		RetryWaitMax: defaultRetryWaitMax,
		RetryMax:     defaultRetryMax,

		CheckForRetry: DefaultRetryPolicy,
		Backoff:       DefaultBackoff,
	}
}

func testRequestWithData(t *testing.T, requestFunc func(url string, headers Headers, data interface{}) ([]byte, error)) {
	_, err := requestFunc(testUrl, testHeaders, testData)
	assert.NoError(t, err, "with headers and data")

	_, err = requestFunc(testUrl, nil, nil)
	assert.NoError(t, err, "without headers and data")

	_, err = requestFunc(testUrl, testHeaders, nil)
	assert.NoError(t, err, "without data")

	_, err = requestFunc(testUrl, nil, testData)
	assert.NoError(t, err, "without headers")
}

func Test_Get(t *testing.T) {
	client := GetFakeHTTPClient()
	_, err := client.Get(testUrl, nil, nil)
	assert.NoError(t, err)

	_, err = client.Get(testUrl, testHeaders, nil)
	assert.NoError(t, err)

	_, err = client.Get(testUrl, nil, testParams)
	assert.NoError(t, err)

	_, err = client.Get(testUrl, testHeaders, testParams)
	assert.NoError(t, err)
}

func Test_Post(t *testing.T) {
	client := GetFakeHTTPClient()
	testRequestWithData(t, client.Post)
}

func Test_Delete(t *testing.T) {
	client := GetFakeHTTPClient()
	testRequestWithData(t, client.Delete)
}

func Test_Put(t *testing.T) {
	client := GetFakeHTTPClient()
	testRequestWithData(t, client.Put)
}

func Test_readBody(t *testing.T) {
	type args struct {
		resp *http.Response
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "nil response",
			args: args{
				resp: nil,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "nil body",
			args: args{
				resp: &http.Response{
					Body: nil,
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "body with data",
			args: args{
				resp: &http.Response{
					Body: io.NopCloser(strings.NewReader("test")),
				},
			},
			want:    []byte("test"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readBody(tt.args.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("readBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readBody() = %v, want %v", got, tt.want)
			}
		})
	}
}
