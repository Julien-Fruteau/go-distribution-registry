package registry

import (
	"context"
	"io"
	"net/http"
	// "strings"
	"testing"
	"time"

  // TODO: using v3 ?!!!
	"github.com/distribution/distribution/v3/configuration"
	reg "github.com/distribution/distribution/v3/registry"
	_ "github.com/distribution/distribution/v3/registry/storage/driver/inmemory"
	"github.com/docker/distribution"
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
)

func setupRegistry(t testing.TB, addr string) *reg.Registry {
	t.Helper()
	config := &configuration.Configuration{}
	config.HTTP.Addr = addr
	config.HTTP.DrainTimeout = time.Duration(10) * time.Second
	config.Storage = map[string]configuration.Parameters{"inmemory": map[string]interface{}{}}
	registry, err := reg.NewRegistry(context.Background(), config)
	assert.NoError(t, err)
	return registry
}

// TODO: use this format in the code base !
type image struct {
	manifest       distribution.Manifest
	manifestDigest digest.Digest
	layers         map[digest.Digest]io.ReadSeeker
}

// func createDummyImage() *image {
// 	// Create a dummy manifest and layers
// 	manifest := distribution.Manifest{}
// 	manifestDigest := digest.FromString("dummy-manifest")
// 	layers := map[digest.Digest]io.ReadSeeker{
// 		digest.FromString("dummy-layer-1"): io.NopCloser(strings.NewReader("dummy-layer-1-content")),
// 		digest.FromString("dummy-layer-2"): io.NopCloser(strings.NewReader("dummy-layer-2-content")),
// 	}
// 	return &image{
// 		manifest:       manifest,
// 		manifestDigest: manifestDigest,
// 		layers:         layers,
// 	}
// }

func TestRegistry(t *testing.T) {
	// setup
	registry := setupRegistry(t, ":5000")
	errchan := make(chan error, 1)
	var err error
	go func() {
		errchan <- registry.ListenAndServe()
	}()
	select {
	case err = <-errchan:
		t.Fatalf("Error listening: %v", err)
	default:
	}

	// defer cleanup
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		registry.Shutdown(ctx)
	}()
	// Wait for some unknown random time for server to start listening
	time.Sleep(3 * time.Second)

	// Make a request to the test registry
	resp, err := http.Get("http://localhost:5000" + "/v2/")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Ensure we get a 200 OK response
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TODO:
// check testutil garbagecollec_test fo upload blobs and construct manifest for test case scenario
// should probably run integration tests against a deployed dkr registry, using a specific version
// func TestDummyImage(t *testing.T) {
// 	// setup
// 	registry := setupRegistry(t, ":5001")
// 	errChan := make(chan error, 1)
// 	go func() {
// 		errChan <- registry.ListenAndServe()
// 	}()
// 	select {
// 	case err := <-errChan:
// 		t.Fatalf("Error listening: %v", err)
// 	default:
// 	}
// 	// defer cleanup
// 	defer func() {
// 		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()
// 		registry.Shutdown(ctx)
// 	}()
// 	// Wait for some unknown random time for server to start listening
// 	time.Sleep(3 * time.Second)
//
// 	// Create a dummy image
// 	dummyImage := createDummyImage()
//
// 	// Post the dummy image manifest
// 	manifestURL := "http://localhost:5001/v2/dummy/manifests/latest"
// 	manifestBody := strings.NewReader(dummyImage.manifest)
// 	resp, err := http.Post(manifestURL, "application/vnd.docker.distribution.manifest.v2+json", manifestBody)
// 	assert.NoError(t, err)
// 	defer resp.Body.Close()
//
// 	// Ensure we get a 201 Created response
// 	assert.Equal(t, http.StatusCreated, resp.StatusCode)
// }
