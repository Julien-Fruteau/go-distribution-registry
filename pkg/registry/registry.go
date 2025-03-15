package registry

import (
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

// 💥 manifest content type should be in Conf
//
// The client should include an Accept header indicating which manifest content types it supports. For more details on the manifest format and content types, see Image Manifest Version 2, Schema 2. In a successful response, the Content-Type header will indicate which manifest type is being returned.

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
