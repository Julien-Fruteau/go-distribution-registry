package registry

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
)

func NewMockRegistry(u string) RegistryClient {
	return RegistryClient{
		baseUrl: u + "/v2/",
		conf: Conf{
			"host",
			"http",
			"user",
			"pwd",
			"mime",
		},
		httpHeaders: map[string]string{
			"Accept":        "mime",
			"Authorization": GetBasicAuthHeader("user", "pwd"),
		},
		httpClient: &http.Client{},
	}
}

func NewTestServer(t testing.TB) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /v2/{name}/manifests/"+string(digest.FromBytes([]byte("ok"))), func(w http.ResponseWriter, r *http.Request) {
		// NOTE: accessing a wildcard value : https://pkg.go.dev/net/http#ServeMux (patterns)
		// w.Write([]byte(r.PathValue("name")))
		w.WriteHeader(http.StatusAccepted)
	})
	server := httptest.NewServer(mux)

	t.Cleanup(server.Close)
	return server
}

func TestDeleteManifestKOWrongDigest(t *testing.T) {
	m := NewTestServer(t)
	cli := NewMockRegistry(m.URL)
	_, err := cli.DeleteManifest("hot", "wrong-digest", "")
	assert.ErrorIs(t, err, ErrInvalidDigest)
}

func TestDeleteManifestOKStatusAccepted(t *testing.T) {
	m := NewTestServer(t)
	cli := NewMockRegistry(m.URL)
	ok, err := cli.DeleteManifest("hot", digest.FromBytes([]byte("ok")), "")
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestDeleteManifestOKStatusNotFound(t *testing.T) {
	m := NewTestServer(t)
	cli := NewMockRegistry(m.URL)
	// use a valid digest, but does not match any test server handler pattern
	// the test server returns 404, one may consider to add the endpoint and return 404
	ok, err := cli.DeleteManifest("hot", digest.FromBytes([]byte("not-found")), "")
	assert.NoError(t, err)
	assert.True(t, ok)
}
