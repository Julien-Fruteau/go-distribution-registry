package registry

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

// Delete Manifest section
func NewTestServerDeleteManifest(t testing.TB) *httptest.Server {
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
	m := NewTestServerDeleteManifest(t)
	cli := NewMockRegistry(m.URL)
	_, err := cli.DeleteManifest("hot", "wrong-digest", "")
	assert.ErrorIs(t, err, ErrInvalidDigest)
}

func TestDeleteManifestOKStatusAccepted(t *testing.T) {
	m := NewTestServerDeleteManifest(t)
	cli := NewMockRegistry(m.URL)
	ok, err := cli.DeleteManifest("hot", digest.FromBytes([]byte("ok")), "")
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestDeleteManifestOKStatusNotFound(t *testing.T) {
	m := NewTestServerDeleteManifest(t)
	cli := NewMockRegistry(m.URL)
	// use a valid digest, but does not match any test server handler pattern
	// the test server returns 404, one may consider to add the endpoint and return 404
	ok, err := cli.DeleteManifest("hot", digest.FromBytes([]byte("not-found")), "")
	assert.NoError(t, err)
	assert.True(t, ok)
}

// Get Manifests, Manifest, ConfigInfo (inspect), Blobs
// setup a mimal comprehensive repo tag chain

var testTags = map[string][]struct {
	Arch         string
	ManifestDgst digest.Digest
	ConfigDgst   digest.Digest
	Layers       []ManifestInfo
	Config       ConfigInfo
}{
	"1": {
		{
			"amd64",
			digest.FromBytes([]byte("amd64")),
			digest.FromBytes([]byte("amd64Conf")),
			[]ManifestInfo{
				{MediaType: MIME_V2_LAYER_GZIP, Digest: string(digest.FromBytes([]byte("1"))), Size: 100},
				{MediaType: MIME_V2_LAYER_GZIP, Digest: string(digest.FromBytes([]byte("2"))), Size: 200},
			},
			ConfigInfo{Architecture: "amd64", Created: "2025-01-01T00:00:00.000Z"},
		},
		{
			"arm64",
			digest.FromBytes([]byte("arm64")),
			digest.FromBytes([]byte("arm64Conf")),
			[]ManifestInfo{
				{MediaType: MIME_V2_LAYER_GZIP, Digest: string(digest.FromBytes([]byte("a"))), Size: 10},
				{MediaType: MIME_V2_LAYER_GZIP, Digest: string(digest.FromBytes([]byte("b"))), Size: 20},
			},
			ConfigInfo{Architecture: "arm64", Created: "2025-01-01T00:42:00.000Z"},
		},
	},
	// "2": {},
}

func HelperTestTags(t testing.TB) []string {
	t.Helper()
	tags := make([]string, len(testTags))
	i := 0
	for v := range testTags {
		tags[i] = v
		i++
	}
	return tags
}

func NewTestServerRepoTag(t testing.TB) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v2/{name}/tags/list", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		b, _ := json.Marshal(TagsResponse{name, HelperTestTags(t)})
		w.Write(b)
	})

	mux.HandleFunc("GET /v2/{name}/manifests/{reference}", func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Accept")
		switch {
		case strings.Contains(ct, MIME_OCI_LIST), strings.Contains(ct, MIME_V2_LIST):
			b, _ := json.Marshal(ManifestsResponse{MediaType: ct, Manifests: []ManifestInfo{
				{MediaType: MIME_OCI_MANIFEST, Digest: string(testTags["1"][0].ManifestDgst), Platform: Platform{testTags["1"][0].Arch, "Linux"}},
				{MediaType: MIME_OCI_MANIFEST, Digest: string(testTags["1"][1].ManifestDgst), Platform: Platform{testTags["1"][1].Arch, "Linux"}},
			}})
			// w.Header().Add("Docker-Content-Digest", string(testTags["1"][0].ManifestDgst))
			w.Write(b)

		case strings.Contains(ct, MIME_OCI_MANIFEST), strings.Contains(ct, MIME_V2_MANIFEST):
			r := r.PathValue("reference")

			var configDgst string
			var layers []ManifestInfo

			for _, v := range testTags["1"] {
				if string(v.ManifestDgst) == r {
					configDgst = string(v.ConfigDgst)
					layers = v.Layers
					break
				}
			}
			if configDgst == "" {
				t.Fatalf("unexpected reference: %v", r)
			}
			b, _ := json.Marshal(ManifestResponse{MediaType: ct, Config: ManifestInfo{MediaType: MIME_OCI_CONFIG, Digest: configDgst}, Layers: layers})

			w.Write(b)

		default:
			t.Fatalf("unexpected Accept: %v", ct)
		}
	})

	mux.HandleFunc("GET /v2/{name}/blobs/{digest}", func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Accept")
		dgst := r.PathValue("digest")

		switch ct {
		case MIME_V2_CONFIG, MIME_OCI_CONFIG:
			var configInfo *ConfigInfo

			for _, v := range testTags["1"] {
				if string(v.ConfigDgst) == dgst {
					configInfo = &v.Config
					break
				}
			}
			if configInfo == nil {
				t.Fatalf("unexpected digest: %v", r)
			}
			b, _ := json.Marshal(configInfo)

			w.Write(b)

		default:
		}
	})
	server := httptest.NewServer(mux)

	t.Cleanup(server.Close)
	return server
}

func TestGetRepoTags(t *testing.T) {
	m := NewTestServerRepoTag(t)
	cli := NewMockRegistry(m.URL)
	want := TagsResponse{"luke", HelperTestTags(t)}

	got, _, err := cli.GetTags("luke")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGetRepoTagsCreateDate(t *testing.T) {
	m := NewTestServerRepoTag(t)
	cli := NewMockRegistry(m.URL)

	var want []RepoTagsCreationResponse
	for i := range testTags["1"] {
		want = append(want, RepoTagsCreationResponse{"yahn:1", testTags["1"][i].Arch, testTags["1"][i].Config.Created})
	}

	got, _ := cli.GetRepositoryTagsCreationDate("yahn")
	for _, g := range got {
		assert.Contains(t, want, g)
	}
}
