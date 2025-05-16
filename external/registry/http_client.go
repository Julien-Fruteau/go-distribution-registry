package registry

import (
	"fmt"
	"net/http"

	"github.com/distribution/distribution/v3/manifest/schema2"
	"github.com/julien-fruteau/go-distribution-registry/internal/env"
)

var (
	// NOTE: some media type already present in distribution
	test     = schema2.MediaTypeManifest
	mime_map = map[string]string{
		"MIME_V2_MANIFEST":    MIME_V2_MANIFEST,
		"MIME_V2_LIST":        MIME_V2_LIST,
		"MIME_V2_CONFIG":      MIME_V2_CONFIG,
		"MIME_V2_LAYER_GZIP":  MIME_V2_LAYER_GZIP,
		"MIME_V2_PLUGIN_JSON": MIME_V2_PLUGIN_JSON,
		"MIME_OCI_MANIFEST":   MIME_OCI_MANIFEST,
		"MIME_OCI_LIST":       MIME_OCI_LIST,
		"MIME_OCI_CONFIG":     MIME_OCI_CONFIG,
	}
)

type RegistryClient struct {
	baseUrl     string
	conf        Conf
	httpHeaders map[string]string
	httpClient  *http.Client
}

type Conf struct {
	host     string
	scheme   string
	username  string
	password string
	mime     string
}

// The client should include an Accept header indicating which manifest content types it supports. For more details on the manifest format and content types, see Image Manifest Version 2, Schema 2. In a successful response, the Content-Type header will indicate which manifest type is being returned.
func NewRegistryClient() RegistryClient {
	host     := env.GetEnvOrDefault("REG_HOST", "localhost")
	scheme   := env.GetEnvOrDefault("REG_SCHEME", "http")
	username := env.GetEnvOrDefault("REG_USER", "admin")
	password := env.GetEnvOrDefault("REG_PASSWORD", "")
	mime     := env.GetEnvOrDefault("REG_MIME", fmt.Sprintf("%s, %s, %s, %s", MIME_V2_MANIFEST, MIME_V2_LIST, MIME_OCI_LIST, MIME_OCI_MANIFEST))
	// mime := env.GetEnvOrDefault("REG_MIME", MIME_V2)

	return RegistryClient{
		baseUrl: scheme + "://" + host + "/v2/",
		conf: Conf{
			host:     host,
			scheme:   scheme,
			username:  username,
			password: password,
			mime:     mime,
		},
		httpHeaders: map[string]string{
			"Accept":        mime,
			"Authorization": GetBasicAuthHeader(username, password),
		},
		httpClient: &http.Client{},
	}
}

// if needing to provide multiple accept header, contatenate
// them separated by coma
func (r *RegistryClient) GetCustomHeader(mediaType string) map[string]string {
	return map[string]string{
		"Accept":        mediaType,
		"Authorization": GetBasicAuthHeader(r.conf.username, r.conf.password),
	}
}
