package registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

// wrapper function to do http call over the registry
func HttpDo[T any](client *http.Client, method, path string, headers, params map[string]string) (response T, respHeaders http.Header, err error) {
	req, err := GetNewRequest(method, path, headers, params)
	if err != nil {
		return response, nil, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return response, nil, fmt.Errorf("error performing http request %s %s: %v", method, req.URL, err)
	}
	defer resp.Body.Close()

	respHeaders = resp.Header

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, respHeaders, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		var respErr RegistryError
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return response, respHeaders, fmt.Errorf("%d, error getting response: %v", resp.StatusCode, string(body))
		}
		return response, respHeaders, fmt.Errorf("%d, error getting response: %v", resp.StatusCode, respErr)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, respHeaders, fmt.Errorf("error unmarshal response: %v", err)
	}

	return response, resp.Header, nil
}
