package registry

import (

	// "net/http"

	"git.isi.nc/go/dtb-tool/pkg/env"
)

type Registry struct {
	BaseUrl string
	Conf    Conf
}

type Conf struct {
	Host     string
	Scheme   string
	Username string
	Password string
}

var (
	host     = env.GetEnvOrDefault("REGISTRY_HOST", "dkr.isi")
	scheme   = env.GetEnvOrDefault("REGISTRY_SCHEME", "http")
	username = env.GetEnvOrDefault("REGISTRY_USER", "admin")
	password = env.GetEnvOrDefault("REGISTRY_PASSWORD", "")
)

func NewRegistry() Registry {
	return Registry{
		// Client:     &http.Client{},
		BaseUrl: scheme + "://" + host + "/v2/",
		Conf: Conf{
			Host:     host,
			Scheme:   scheme,
			Username: username,
			Password: password,
		},
		// AuthHeader: getAuthHeader(username, password),
		// pagination: Pagination{
		// 	n:    "100",
		// 	last: "",
		// },
	}
}
