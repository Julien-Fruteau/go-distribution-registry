package registry

import "net/http"

type Registry struct {
	baseUrl     string
	conf        Conf
	httpHeaders map[string]string
	httpClient  *http.Client
}

type Conf struct {
	host     string
	scheme   string
	usename  string
	password string
	mime     string
}
