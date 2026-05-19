package registry

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeDockerConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	require.NoError(t, os.WriteFile(path, []byte(content), 0600))
	return path
}

func encodeAuth(user, pass string) string {
	return base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
}

// --- decodeAuthEntry ---

func TestDecodeAuthEntry_Valid(t *testing.T) {
	u, p, err := decodeAuthEntry(encodeAuth("alice", "s3cr3t"))
	require.NoError(t, err)
	assert.Equal(t, "alice", u)
	assert.Equal(t, "s3cr3t", p)
}

func TestDecodeAuthEntry_PasswordWithColon(t *testing.T) {
	u, p, err := decodeAuthEntry(encodeAuth("alice", "pass:with:colons"))
	require.NoError(t, err)
	assert.Equal(t, "alice", u)
	assert.Equal(t, "pass:with:colons", p)
}

func TestDecodeAuthEntry_NotBase64(t *testing.T) {
	_, _, err := decodeAuthEntry("not-valid-base64!!!")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "base64-decoding")
}

func TestDecodeAuthEntry_NoColon(t *testing.T) {
	_, _, err := decodeAuthEntry(base64.StdEncoding.EncodeToString([]byte("nocolon")))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user:pass format")
}

// --- lookupDockerCredentials ---

func TestLookupDockerCredentials_InlineAuthBareHost(t *testing.T) {
	cfg := `{"auths":{"registry.example.com":{"auth":"` + encodeAuth("user1", "pass1") + `"}}}`
	path := writeDockerConfig(t, cfg)

	creds, err := lookupDockerCredentials(path, "registry.example.com")
	require.NoError(t, err)
	assert.Equal(t, "user1", creds.Username)
	assert.Equal(t, "pass1", creds.Password)
}

func TestLookupDockerCredentials_InlineAuthHTTPSPrefixedHost(t *testing.T) {
	cfg := `{"auths":{"https://registry.example.com":{"auth":"` + encodeAuth("user2", "pass2") + `"}}}`
	path := writeDockerConfig(t, cfg)

	creds, err := lookupDockerCredentials(path, "registry.example.com")
	require.NoError(t, err)
	assert.Equal(t, "user2", creds.Username)
	assert.Equal(t, "pass2", creds.Password)
}

func TestLookupDockerCredentials_NoMatchReturnsError(t *testing.T) {
	cfg := `{"auths":{"other.host.com":{"auth":"` + encodeAuth("u", "p") + `"}}}`
	path := writeDockerConfig(t, cfg)

	_, err := lookupDockerCredentials(path, "registry.example.com")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no docker credentials found")
}

func TestLookupDockerCredentials_MissingConfigFile(t *testing.T) {
	_, err := lookupDockerCredentials("/nonexistent/path/config.json", "registry.example.com")
	require.Error(t, err)
	assert.True(t, errors.Is(err, os.ErrNotExist))
}

func TestLookupDockerCredentials_MalformedBase64(t *testing.T) {
	cfg := `{"auths":{"registry.example.com":{"auth":"!!!not-base64"}}}`
	path := writeDockerConfig(t, cfg)

	_, err := lookupDockerCredentials(path, "registry.example.com")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "base64-decoding")
}

func TestLookupDockerCredentials_MalformedJSON(t *testing.T) {
	path := writeDockerConfig(t, `{bad json`)

	_, err := lookupDockerCredentials(path, "registry.example.com")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing docker config")
}

func TestLookupDockerCredentials_CredHelperTakesPrecedence(t *testing.T) {
	// credHelpers entry is present alongside inline auth.
	// The helper binary won't exist, so we verify the error is from the helper
	// attempt (not the inline auth decode), confirming lookup order.
	cfg := `{
		"credHelpers":{"registry.example.com":"nonexistent-helper"},
		"auths":{"registry.example.com":{"auth":"` + encodeAuth("inline", "creds") + `"}}
	}`
	path := writeDockerConfig(t, cfg)

	_, err := lookupDockerCredentials(path, "registry.example.com")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nonexistent-helper")
}

// --- normalizeDockerHost ---

func TestNormalizeDockerHost_BareHost(t *testing.T) {
	assert.Equal(t, "registry.example.com", normalizeDockerHost("registry.example.com"))
}

func TestNormalizeDockerHost_HTTPSPrefix(t *testing.T) {
	assert.Equal(t, "registry.example.com", normalizeDockerHost("https://registry.example.com"))
}

func TestNormalizeDockerHost_HTTPSPrefixWithPath(t *testing.T) {
	assert.Equal(t, "index.docker.io", normalizeDockerHost("https://index.docker.io/v1/"))
}

func TestNormalizeDockerHost_HTTPPrefix(t *testing.T) {
	assert.Equal(t, "localhost:5000", normalizeDockerHost("http://localhost:5000"))
}

// --- resolveHost ---

func TestResolveHost_EnvVarSet(t *testing.T) {
	t.Setenv("REG_HOST", "explicit.host.com")
	host, err := resolveHost("/nonexistent/config.json")
	require.NoError(t, err)
	assert.Equal(t, "explicit.host.com", host)
}

func TestResolveHost_SingleRegistry(t *testing.T) {
	cfg := `{"auths":{"dkr.example.com":{"auth":"` + encodeAuth("u", "p") + `"}}}`
	path := writeDockerConfig(t, cfg)
	os.Unsetenv("REG_HOST")

	host, err := resolveHost(path)
	require.NoError(t, err)
	assert.Equal(t, "dkr.example.com", host)
}

func TestResolveHost_SingleRegistryHTTPSKey(t *testing.T) {
	cfg := `{"auths":{"https://dkr.example.com":{"auth":"` + encodeAuth("u", "p") + `"}}}`
	path := writeDockerConfig(t, cfg)
	os.Unsetenv("REG_HOST")

	host, err := resolveHost(path)
	require.NoError(t, err)
	assert.Equal(t, "dkr.example.com", host)
}

func TestResolveHost_MultipleRegistriesReturnsError(t *testing.T) {
	cfg := `{"auths":{
		"alpha.example.com":{"auth":"` + encodeAuth("u", "p") + `"},
		"beta.example.com":{"auth":"` + encodeAuth("u", "p") + `"}
	}}`
	path := writeDockerConfig(t, cfg)
	os.Unsetenv("REG_HOST")

	_, err := resolveHost(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "multiple registries")
	assert.Contains(t, err.Error(), "alpha.example.com")
	assert.Contains(t, err.Error(), "beta.example.com")
}

func TestResolveHost_NoRegistriesReturnsError(t *testing.T) {
	path := writeDockerConfig(t, `{"auths":{}}`)
	os.Unsetenv("REG_HOST")

	_, err := resolveHost(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no registries found")
}

func TestResolveHost_MissingConfigReturnsError(t *testing.T) {
	os.Unsetenv("REG_HOST")
	_, err := resolveHost("/nonexistent/config.json")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "REG_HOST not set")
}

// --- resolveCredentials ---

func TestResolveCredentials_EnvUserTakesPriority(t *testing.T) {
	cfg := `{"auths":{"myhost":{"auth":"` + encodeAuth("dockeruser", "dockerpass") + `"}}}`
	path := writeDockerConfig(t, cfg)

	t.Setenv("REG_USER", "envuser")
	t.Setenv("REG_PASSWORD", "envpass")

	u, p := resolveCredentials(path, "myhost")
	assert.Equal(t, "envuser", u)
	assert.Equal(t, "envpass", p)
}

func TestResolveCredentials_DockerConfigUsed(t *testing.T) {
	cfg := `{"auths":{"myhost":{"auth":"` + encodeAuth("dockeruser", "dockerpass") + `"}}}`
	path := writeDockerConfig(t, cfg)

	os.Unsetenv("REG_USER")
	os.Unsetenv("REG_PASSWORD")

	u, p := resolveCredentials(path, "myhost")
	assert.Equal(t, "dockeruser", u)
	assert.Equal(t, "dockerpass", p)
}

func TestResolveCredentials_FallsBackToDefault(t *testing.T) {
	os.Unsetenv("REG_USER")
	os.Unsetenv("REG_PASSWORD")

	u, p := resolveCredentials("/nonexistent/config.json", "myhost")
	assert.Equal(t, "admin", u)
	assert.Equal(t, "", p)
}

func TestResolveCredentials_FallsBackWithREGPassword(t *testing.T) {
	os.Unsetenv("REG_USER")
	t.Setenv("REG_PASSWORD", "onlypass")

	u, p := resolveCredentials("/nonexistent/config.json", "myhost")
	assert.Equal(t, "admin", u)
	assert.Equal(t, "onlypass", p)
}
