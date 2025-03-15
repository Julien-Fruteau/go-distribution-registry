package registry

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

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

type image struct {
	manifest       distribution.Manifest
	manifestDigest digest.Digest
	layers         map[digest.Digest]io.ReadSeeker
}

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
func TestDummyImage(t *testing.T) {
	// setup
	registry := setupRegistry(t, ":5001")
	// defer cleanup
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		registry.Shutdown(ctx)
	}()
	// Wait for some unknown random time for server to start listening
	time.Sleep(3 * time.Second)
}
