package http

import (
	"encoding/base64"
	"net/http"
)

func GetBasicAuthHeader(username, password string) string {
	auth := username + ":" + password
	base64Auth := base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + base64Auth
}


func GetNewRequest(method string, url string, params map[string]string) (*http.Request, error) {
  req, err := http.NewRequest(method, url, nil)
  if err != nil {
    return nil, err
  }

  values := req.URL.Query()
  for k, v := range params {
    values.Add(k, v)
  }

  req.URL.RawQuery = values.Encode()

  return req, nil
}
