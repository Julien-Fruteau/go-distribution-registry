package http

import (
	"encoding/base64"
	"net/http"
	"net/url"
)

func GetBasicAuthHeader(username, password string) string {
	auth := username + ":" + password
	base64Auth := base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + base64Auth
}

func GetNewRequest(method string, u string, headers map[string]string, params map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		return nil, err
	}

	values := req.URL.Query()
	for k, v := range params {
		value := url.QueryEscape(v)
		values[k] = append(values[k], value)
	}

	req.URL.RawQuery = values.Encode()

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func SetRequestHeader(request *http.Request, headers map[string]string) {
	for k, v := range headers {
		request.Header.Set(k, v)
	}
}
