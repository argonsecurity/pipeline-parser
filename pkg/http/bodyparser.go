package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/url"
	"strings"
)

const (
	ContentTypeHeader            = "Content-Type"
	FormDataContentType          = "application/x-www-form-urlencoded"
	JsonContentType              = "application/json"
	MultipartFormDataContentType = "multipart/form-data"
)

var (
	contentTypeToParserMap = map[string]func(data any, contentType string) (*bytes.Buffer, string, error){
		JsonContentType:              parseJson,
		FormDataContentType:          parseFormData,
		MultipartFormDataContentType: parseMultipartFormData,
	}
)

func parseMultipartFormData(data any, contentType string) (*bytes.Buffer, string, error) {
	values, ok := data.(map[string]string)
	if !ok {
		return nil, "", fmt.Errorf("data is not a map[string]string")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range values {
		part, err := writer.CreateFormField(k)
		if err != nil {
			return nil, "", err
		}
		_, err = part.Write([]byte(v))
		if err != nil {
			return nil, "", err
		}
	}
	writer.Close()
	return body, writer.FormDataContentType(), nil
}

func parseFormData(data any, contentType string) (*bytes.Buffer, string, error) {
	values, ok := data.(map[string]string)
	if !ok {
		return nil, "", fmt.Errorf("data is not a map[string]string")
	}
	form := url.Values{}
	for k, v := range values {
		form.Add(k, v)
	}
	return bytes.NewBuffer([]byte(form.Encode())), contentType, nil
}

func parseJson(data any, contentType string) (*bytes.Buffer, string, error) {
	parsedData, err := json.Marshal(data)
	if err != nil {
		return nil, "", err
	}
	return bytes.NewBuffer(parsedData), contentType, nil
}

func getContentTypeName(contentTypeValue string) string {
	split := strings.Split(contentTypeValue, ";")
	return split[0]
}

func parseData(data any, headers Headers) (*bytes.Buffer, string, error) {
	if data == nil {
		return nil, "", nil
	}

	contentTypeValue, ok := headers[ContentTypeHeader]
	if !ok {
		return parseJson(data, contentTypeValue)
	}

	if data == nil {
		return nil, contentTypeValue, nil
	}

	contentType := getContentTypeName(contentTypeValue)
	parser, ok := contentTypeToParserMap[contentType]
	if !ok {
		return parseJson(data, contentTypeValue)
	}

	return parser(data, contentType)
}
