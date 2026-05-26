package registry

import (
	"fmt"
	"net/http"
	"strings"

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
	username string
	password string
	mime     string
}

func NewRegistryClient(optHost string) (RegistryClient, error) {
	host, err := resolveHost(dockerConfigPath, optHost)
	if err != nil {
		return RegistryClient{}, err
	}
	scheme := env.GetEnvOrDefault("REG_SCHEME", "http")
	username, password, err := resolveCredentials(dockerConfigPath, host)
	if err != nil {
		return RegistryClient{}, err
	}
	mime := env.GetEnvOrDefault("REG_MIME", fmt.Sprintf("%s, %s, %s, %s", MIME_V2_MANIFEST, MIME_V2_LIST, MIME_OCI_LIST, MIME_OCI_MANIFEST))

	return RegistryClient{
		baseUrl: scheme + "://" + host + "/v2/",
		conf: Conf{
			host:     host,
			scheme:   scheme,
			username: username,
			password: password,
			mime:     mime,
		},
		httpHeaders: map[string]string{
			"Accept":        mime,
			"Authorization": GetBasicAuthHeader(username, password),
		},
		httpClient: &http.Client{},
	}, nil
}

// NormalizeName strips the registry host prefix from a repository name if present.
// For example, "docker.mine.com/group/backend" becomes "group/backend" when the
// configured host is "docker.mine.com".
func (r *RegistryClient) NormalizeName(name string) string {
	prefix := r.conf.host + "/"
	if strings.HasPrefix(name, prefix) {
		return name[len(prefix):]
	}
	return name
}

// if needing to provide multiple accept header, contatenate
// them separated by coma
func (r *RegistryClient) GetCustomHeader(mediaType string) map[string]string {
	return map[string]string{
		"Accept":        mediaType,
		"Authorization": GetBasicAuthHeader(r.conf.username, r.conf.password),
	}
}
