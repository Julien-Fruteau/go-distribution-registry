package registry

import (
	"net/http"

	"github.com/julien-fruteau/go/distctl/internal/env"
	httpUtils "github.com/julien-fruteau/go/distctl/internal/http"
)

var (
	mime_map = map[string]string{
		"MIME_V2":                       MIME_V2,
		"MIME_V2_LIST":                  MIME_V2_LIST,
		"MIME_V2_CONTAINER_CONFIG_JSON": MIME_V2_CONTAINER_CONFIG_JSON,
		"MIME_V2_LAYER_GZIP":            MIME_V2_LAYER_GZIP,
		"MIME_V2_PLUGIN_JSON":           MIME_V2_PLUGIN_JSON,
		"MIME_V1":                       MIME_V1,
	}
	host     = env.GetEnvOrDefault("REG_HOST", "localhost")
	scheme   = env.GetEnvOrDefault("REG_SCHEME", "http")
	username = env.GetEnvOrDefault("REG_USER", "admin")
	password = env.GetEnvOrDefault("REG_PASSWORD", "")
	mime     = env.LookupEnvOrDefault(mime_map, "REG_MIME", MIME_V2)
	// mime = env.GetEnvOrDefault("REG_MIME", MIME_V2)
)

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

//
// The client should include an Accept header indicating which manifest content types it supports. For more details on the manifest format and content types, see Image Manifest Version 2, Schema 2. In a successful response, the Content-Type header will indicate which manifest type is being returned.

func NewRegistry() Registry {
	return Registry{
		// Client:     &http.Client{},
		baseUrl: scheme + "://" + host + "/v2/",
		conf: Conf{
			host:     host,
			scheme:   scheme,
			usename:  username,
			password: password,
			mime:     mime,
		},
		httpHeaders: map[string]string{
			"Accept":        mime,
			"Authorization": httpUtils.GetBasicAuthHeader(username, password),
		},
		httpClient: &http.Client{},
	}
}
