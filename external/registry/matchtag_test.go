package registry

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func newMatchTagRegistry(t *testing.T, tags []string, manifests map[string]string, mediaTypes map[string]string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v2/test/image/tags/list", func(w http.ResponseWriter, _ *http.Request) {
		require.NoError(t, json.NewEncoder(w).Encode(TagsResponse{Name: "test/image", Tags: tags}))
	})
	mux.HandleFunc("GET /v2/test/image/manifests/{reference}", func(w http.ResponseWriter, r *http.Request) {
		reference := r.PathValue("reference")
		manifest, ok := manifests[reference]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", mediaTypes[reference])
		_, err := w.Write([]byte(manifest))
		require.NoError(t, err)
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server
}

func TestMatchTagReturnsHighestDockerV2VersionWithMatchingConfig(t *testing.T) {
	manifests := map[string]string{
		"stable":     `{"schemaVersion":2,"mediaType":"` + MIME_V2_MANIFEST + `","config":{"digest":"sha256:matching"}}`,
		"1.9.0":      `{"schemaVersion":2,"mediaType":"` + MIME_V2_MANIFEST + `","config":{"digest":"sha256:matching"}}`,
		"1.10.0":     `{"schemaVersion":2,"mediaType":"` + MIME_V2_MANIFEST + `","config":{"digest":"sha256:matching"}}`,
		"2.0.0":      `{"schemaVersion":2,"mediaType":"` + MIME_V2_MANIFEST + `","config":{"digest":"sha256:matching"}}`,
		"10.0.0":     `{"schemaVersion":2,"mediaType":"` + MIME_V2_MANIFEST + `","config":{"digest":"sha256:different"}}`,
		"not-semver": `{"schemaVersion":2,"mediaType":"` + MIME_V2_MANIFEST + `","config":{"digest":"sha256:matching"}}`,
	}
	mediaTypes := make(map[string]string, len(manifests))
	for tag := range manifests {
		mediaTypes[tag] = MIME_V2_MANIFEST
	}
	server := newMatchTagRegistry(t, []string{"stable", "1.9.0", "1.10.0", "2.0.0", "10.0.0", "not-semver"}, manifests, mediaTypes)

	client := NewMockRegistry(server.URL)
	match, err := client.MatchTag("test/image", "stable")

	require.NoError(t, err)
	require.Equal(t, "2.0.0", match)
}

func TestMatchTagMatchesOCIManifestConfig(t *testing.T) {
	manifests := map[string]string{
		"stable": `{"schemaVersion":2,"mediaType":"` + MIME_OCI_MANIFEST + `","config":{"digest":"sha256:matching"}}`,
		"3.1.4":  `{"schemaVersion":2,"mediaType":"` + MIME_OCI_MANIFEST + `","config":{"digest":"sha256:matching"}}`,
	}
	mediaTypes := map[string]string{"stable": MIME_OCI_MANIFEST, "3.1.4": MIME_OCI_MANIFEST}
	server := newMatchTagRegistry(t, []string{"stable", "3.1.4"}, manifests, mediaTypes)
	client := NewMockRegistry(server.URL)

	match, err := client.MatchTag("test/image", "stable")

	require.NoError(t, err)
	require.Equal(t, "3.1.4", match)
}

func TestMatchTagMatchesSchema1LayersIndependentOfTagPayload(t *testing.T) {
	manifests := map[string]string{
		"stable": `{"schemaVersion":1,"name":"test/image","tag":"stable","fsLayers":[{"blobSum":"sha256:b"},{"blobSum":"sha256:a"}]}`,
		"4.2.0":  `{"schemaVersion":1,"name":"test/image","tag":"4.2.0","fsLayers":[{"blobSum":"sha256:a"},{"blobSum":"sha256:b"}]}`,
		"5.0.0":  `{"schemaVersion":1,"name":"test/image","tag":"5.0.0","fsLayers":[{"blobSum":"sha256:different"}]}`,
	}
	mediaTypes := map[string]string{
		"stable": MIME_V1_PRETTYJWS,
		"4.2.0":  MIME_V1_MANIFEST,
		"5.0.0":  MIME_V1_MANIFEST,
	}
	server := newMatchTagRegistry(t, []string{"stable", "4.2.0", "5.0.0"}, manifests, mediaTypes)
	client := NewMockRegistry(server.URL)

	match, err := client.MatchTag("test/image", "stable")

	require.NoError(t, err)
	require.Equal(t, "4.2.0", match)
}

func TestMatchTagReturnsErrorWhenNoVersionedTagMatches(t *testing.T) {
	manifests := map[string]string{
		"stable": `{"schemaVersion":2,"mediaType":"` + MIME_V2_MANIFEST + `","config":{"digest":"sha256:reference"}}`,
		"1.0.0":  `{"schemaVersion":2,"mediaType":"` + MIME_V2_MANIFEST + `","config":{"digest":"sha256:different"}}`,
	}
	mediaTypes := map[string]string{"stable": MIME_V2_MANIFEST, "1.0.0": MIME_V2_MANIFEST}
	server := newMatchTagRegistry(t, []string{"stable", "1.0.0"}, manifests, mediaTypes)
	client := NewMockRegistry(server.URL)

	match, err := client.MatchTag("test/image", "stable")

	require.Empty(t, match)
	require.EqualError(t, err, "no versioned tag matches test/image stable")
}

func TestMatchTagMatchesAndOrdersSemanticVersionPrereleases(t *testing.T) {
	manifests := map[string]string{
		"staging":         `{"schemaVersion":2,"mediaType":"` + MIME_V2_MANIFEST + `","config":{"digest":"sha256:matching"}}`,
		"1.2.0-qual.1531": `{"schemaVersion":2,"mediaType":"` + MIME_V2_MANIFEST + `","config":{"digest":"sha256:matching"}}`,
		"1.2.0-qual.1532": `{"schemaVersion":2,"mediaType":"` + MIME_V2_MANIFEST + `","config":{"digest":"sha256:matching"}}`,
		"1.2.0-dev.1533":  `{"schemaVersion":2,"mediaType":"` + MIME_V2_MANIFEST + `","config":{"digest":"sha256:different"}}`,
	}
	mediaTypes := make(map[string]string, len(manifests))
	for tag := range manifests {
		mediaTypes[tag] = MIME_V2_MANIFEST
	}
	server := newMatchTagRegistry(t, []string{"staging", "1.2.0-qual.1531", "1.2.0-qual.1532", "1.2.0-dev.1533"}, manifests, mediaTypes)
	client := NewMockRegistry(server.URL)

	match, err := client.MatchTag("test/image", "staging")

	require.NoError(t, err)
	require.Equal(t, "1.2.0-qual.1532", match)
}

func TestMatchTagPrefersReleaseOverPrerelease(t *testing.T) {
	manifests := map[string]string{
		"staging":    `{"schemaVersion":2,"mediaType":"` + MIME_V2_MANIFEST + `","config":{"digest":"sha256:matching"}}`,
		"1.2.0-rc.9": `{"schemaVersion":2,"mediaType":"` + MIME_V2_MANIFEST + `","config":{"digest":"sha256:matching"}}`,
		"1.2.0":      `{"schemaVersion":2,"mediaType":"` + MIME_V2_MANIFEST + `","config":{"digest":"sha256:matching"}}`,
	}
	mediaTypes := make(map[string]string, len(manifests))
	for tag := range manifests {
		mediaTypes[tag] = MIME_V2_MANIFEST
	}
	server := newMatchTagRegistry(t, []string{"staging", "1.2.0-rc.9", "1.2.0"}, manifests, mediaTypes)
	client := NewMockRegistry(server.URL)

	match, err := client.MatchTag("test/image", "staging")

	require.NoError(t, err)
	require.Equal(t, "1.2.0", match)
}
