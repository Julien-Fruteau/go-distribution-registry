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

func TestDeleteManifestKOWrongDigest(t *testing.T) {
  // NOTE: httptest server 
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusAccepted) }))
	defer mockServer.Close()

	cli := NewMockRegistry(mockServer.URL)
  _, err := cli.DeleteManifest("hot",  "wrong", "potato")
  assert.ErrorIs(t, err, ErrInvalidDigest )
}

func TestDeleteManifestOKStatusAccepted(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusAccepted) }))
	defer mockServer.Close()

	cli := NewMockRegistry(mockServer.URL)
  ok, err := cli.DeleteManifest("hot",  digest.FromBytes([]byte("abc")), "potato")
  assert.NoError(t, err)
  assert.True(t, ok)
}

func TestDeleteManifestOKStatusNotFound(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) }))
	defer mockServer.Close()

	cli := NewMockRegistry(mockServer.URL)
  ok, err := cli.DeleteManifest("hot",  digest.FromBytes([]byte("abc")), "potato")
  assert.NoError(t, err)
  assert.True(t, ok)
}
