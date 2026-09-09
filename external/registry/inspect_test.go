package registry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/distribution/distribution/v3/configuration"
	"github.com/distribution/distribution/v3/registry/handlers"
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/require"
)

func newInspectRegistry(t *testing.T) *httptest.Server {
	t.Helper()
	config := &configuration.Configuration{}
	config.Storage = configuration.Storage{
		"inmemory":    {},
		"maintenance": {"uploadpurging": map[interface{}]interface{}{"enabled": false}},
	}
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	server := httptest.NewServer(handlers.NewApp(ctx, config))
	t.Cleanup(server.Close)
	return server
}

func registryRequest(t *testing.T, method, url, mediaType, body string, status int) http.Header {
	t.Helper()
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", mediaType)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, status, resp.StatusCode, "%s", responseBody)
	return resp.Header
}

func pushInspectImage(t *testing.T, server *httptest.Server, tag, manifestType, configType, architecture string) string {
	t.Helper()
	config := fmt.Sprintf(`{"architecture":%q,"os":"linux","created":"2025-01-01T00:00:00Z","config":{"Env":["APP=test"]},"rootfs":{"type":"layers","diff_ids":[]}}`, architecture)
	configDigest := digest.FromString(config).String()
	headers := registryRequest(t, http.MethodPost, server.URL+"/v2/test/image/blobs/uploads/", "", "", http.StatusAccepted)
	registryRequest(t, http.MethodPut, headers.Get("Location")+"&digest="+configDigest, "application/octet-stream", config, http.StatusCreated)
	manifest := fmt.Sprintf(`{"schemaVersion":2,"mediaType":%q,"config":{"mediaType":%q,"digest":%q,"size":%d},"layers":[]}`, manifestType, configType, configDigest, len(config))
	registryRequest(t, http.MethodPut, server.URL+"/v2/test/image/manifests/"+tag, manifestType, manifest, http.StatusCreated)
	return manifest
}

func TestInspectSingleArchitecture(t *testing.T) {
	for _, tc := range []struct {
		name, manifestType, configType string
	}{
		{"docker", MIME_V2_MANIFEST, MIME_V2_CONFIG},
		{"oci", MIME_OCI_MANIFEST, MIME_OCI_CONFIG},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := newInspectRegistry(t)
			pushInspectImage(t, server, "single", tc.manifestType, tc.configType, "amd64")
			client := NewMockRegistry(server.URL)
			infos, err := client.Inspect("test/image", "single")
			require.NoError(t, err)
			require.Len(t, infos, 1)
			require.Equal(t, "amd64", infos[0].Architecture)
			require.Equal(t, "linux", infos[0].OS)
			require.Equal(t, "2025-01-01T00:00:00Z", infos[0].Created)
			require.Equal(t, []string{"APP=test"}, infos[0].Config.Env)
		})
	}
}

func TestInspectMultipleArchitectures(t *testing.T) {
	for _, tc := range []struct {
		name, indexType, manifestType, configType string
	}{
		{"docker", MIME_V2_LIST, MIME_V2_MANIFEST, MIME_V2_CONFIG},
		{"oci", MIME_OCI_LIST, MIME_OCI_MANIFEST, MIME_OCI_CONFIG},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := newInspectRegistry(t)
			var descriptors []string
			for _, architecture := range []string{"amd64", "arm64"} {
				manifest := pushInspectImage(t, server, architecture, tc.manifestType, tc.configType, architecture)
				descriptors = append(descriptors, fmt.Sprintf(`{"mediaType":%q,"digest":%q,"size":%d,"platform":{"architecture":%q,"os":"linux"}}`, tc.manifestType, digest.FromString(manifest).String(), len(manifest), architecture))
			}
			index := fmt.Sprintf(`{"schemaVersion":2,"mediaType":%q,"manifests":[%s]}`, tc.indexType, strings.Join(descriptors, ","))
			registryRequest(t, http.MethodPut, server.URL+"/v2/test/image/manifests/multi", tc.indexType, index, http.StatusCreated)
			client := NewMockRegistry(server.URL)
			infos, err := client.Inspect("test/image", "multi")
			require.NoError(t, err)
			require.Len(t, infos, 2)
			var architectures []string
			for _, info := range infos {
				architectures = append(architectures, info.Architecture)
				require.Equal(t, "linux", info.OS)
				require.Equal(t, "2025-01-01T00:00:00Z", info.Created)
				require.Equal(t, []string{"APP=test"}, info.Config.Env)
			}
			require.ElementsMatch(t, []string{"amd64", "arm64"}, architectures)
		})
	}
}
