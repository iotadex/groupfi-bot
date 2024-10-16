package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func HttpRequest(url string, method string, postParams []byte, headers map[string]string) ([]byte, error) {
	httpClient := &http.Client{}
	var reader io.Reader
	if len(postParams) > 0 {
		reader = strings.NewReader(string(postParams))
		if headers == nil {
			headers = map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
		}
	} else {
		reader = nil
	}

	request, err := http.NewRequest(method, url, reader)
	if nil != err {
		return nil, fmt.Errorf("NewRequest error. %v", err)
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	response, err := httpClient.Do(request)
	if nil != err {
		return nil, fmt.Errorf("do the request error. %v", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if nil != err {
		return nil, fmt.Errorf("readAll response.Body error. %v", err)
	}

	return body, nil
}

func HttpGet(url string) ([]byte, error) {
	return HttpRequest(url, "GET", nil, nil)
}

func HttpGetWithHeader(url string, headers map[string]string) ([]byte, error) {
	return HttpRequest(url, "GET", nil, headers)
}

func HttpPost(url string, postParams interface{}, headers map[string]string) ([]byte, error) {
	httpClient := http.Client{}

	dataByte, err := json.Marshal(postParams)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(dataByte)
	request, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if nil != err {
		return nil, fmt.Errorf("NewRequest error. %v", err)
	}
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	response, err := httpClient.Do(request)
	if nil != err {
		return nil, fmt.Errorf("do the request error. %v", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if nil != err {
		return nil, fmt.Errorf("readAll response.Body error. %v", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("response error. %s", string(body))
	}

	return body, nil
}
