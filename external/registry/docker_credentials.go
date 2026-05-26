package registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	dockerclient "github.com/docker/docker-credential-helpers/client"
)

// dockerConfigPath is the path to ~/.docker/config.json.
// Package-level var so tests can substitute a fixture path without changing NewRegistryClient's signature.
var dockerConfigPath = filepath.Join(os.Getenv("HOME"), ".docker", "config.json")

type dockerConfigAuth struct {
	Auth string `json:"auth"` // base64-encoded "user:pass"
}

type dockerConfig struct {
	Auths       map[string]dockerConfigAuth `json:"auths"`
	CredHelpers map[string]string           `json:"credHelpers"` // host -> helper name suffix
	CredsStore  string                      `json:"credsStore"`  // global fallback helper
}

type dockerCredentials struct {
	Username string
	Password string
}

func loadDockerConfig(configPath string) (dockerConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return dockerConfig{}, fmt.Errorf("reading docker config %s: %w", configPath, err)
	}
	var cfg dockerConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return dockerConfig{}, fmt.Errorf("parsing docker config %s: %w", configPath, err)
	}
	return cfg, nil
}

// decodeAuthEntry base64-decodes an inline auth value and splits on the first colon.
// Passwords may contain colons, so SplitN with n=2 matches Docker's own behaviour.
func decodeAuthEntry(encoded string) (username, password string, err error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", fmt.Errorf("base64-decoding auth entry: %w", err)
	}
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("auth entry %q is not in user:pass format", string(decoded))
	}
	return parts[0], parts[1], nil
}

func callCredHelper(helper, serverHost string) (dockerCredentials, error) {
	programFunc := dockerclient.NewShellProgramFunc("docker-credential-" + helper)
	creds, err := dockerclient.Get(programFunc, serverHost)
	if err != nil {
		return dockerCredentials{}, fmt.Errorf("credential helper %q for host %q: %w", helper, serverHost, err)
	}
	return dockerCredentials{Username: creds.Username, Password: creds.Secret}, nil
}

// lookupDockerCredentials resolves credentials for serverHost from ~/.docker/config.json.
// Resolution order: per-host credHelper → inline auths (bare and https:// prefixed) → global credsStore.
func lookupDockerCredentials(configPath, serverHost string) (dockerCredentials, error) {
	cfg, err := loadDockerConfig(configPath)
	if err != nil {
		return dockerCredentials{}, err
	}

	if helper, ok := cfg.CredHelpers[serverHost]; ok {
		return callCredHelper(helper, serverHost)
	}

	for _, key := range []string{serverHost, "https://" + serverHost} {
		if entry, ok := cfg.Auths[key]; ok && entry.Auth != "" {
			u, p, err := decodeAuthEntry(entry.Auth)
			if err != nil {
				return dockerCredentials{}, fmt.Errorf("host %s: %w", serverHost, err)
			}
			return dockerCredentials{Username: u, Password: p}, nil
		}
	}

	if cfg.CredsStore != "" {
		return callCredHelper(cfg.CredsStore, serverHost)
	}

	return dockerCredentials{}, fmt.Errorf("no docker credentials found for host %q", serverHost)
}

// normalizeDockerHost strips the scheme and any path from a docker config key,
// returning a bare hostname (with optional port).
func normalizeDockerHost(key string) string {
	key = strings.TrimPrefix(key, "https://")
	key = strings.TrimPrefix(key, "http://")
	if idx := strings.IndexByte(key, '/'); idx != -1 {
		key = key[:idx]
	}
	return key
}

// resolveHost returns the registry hostname to connect to.
//
// Priority:
//  1. optHost provided via -reg flag → use it.
//  2. REG_HOST explicitly set in environment → use it.
//  3. Exactly one registry found in docker config → use that host automatically.
//  4. Multiple registries found → error asking to use -reg or set REG_HOST.
//  5. No registries found → error.
func resolveHost(configPath, optHost string) (string, error) {
	if optHost != "" {
		return optHost, nil
	}

	if host, ok := os.LookupEnv("REG_HOST"); ok {
		return host, nil
	}

	cfg, err := loadDockerConfig(configPath)
	if err != nil {
		return "", fmt.Errorf("REG_HOST not set and cannot read docker config: %w", err)
	}

	seen := make(map[string]struct{})
	for key := range cfg.Auths {
		if h := normalizeDockerHost(key); h != "" {
			seen[h] = struct{}{}
		}
	}
	for h := range cfg.CredHelpers {
		if h != "" {
			seen[h] = struct{}{}
		}
	}

	switch len(seen) {
	case 0:
		return "", fmt.Errorf("REG_HOST not set and no registries found in docker config")
	case 1:
		for h := range seen {
			return h, nil
		}
	default:
		hosts := make([]string, 0, len(seen))
		for h := range seen {
			hosts = append(hosts, h)
		}
		sort.Strings(hosts)
		return "", fmt.Errorf("multiple registries found in docker config; use -reg <host> or set REG_HOST to one of: %s", strings.Join(hosts, ", "))
	}
	return "", nil
}

// resolveCredentials determines the username and password for a registry host.
//
// Priority:
//  1. REG_USER explicitly set in environment → use REG_USER + REG_PASSWORD (backward compat).
//  2. Docker credential lookup from config.json / credstore for host.
//  3. No credentials found → return empty strings (no error); the registry will
//     return 401 if authentication is actually required.
func resolveCredentials(configPath, host string) (username, password string, err error) {
	if envUser, ok := os.LookupEnv("REG_USER"); ok {
		envPass, _ := os.LookupEnv("REG_PASSWORD")
		return envUser, envPass, nil
	}

	if creds, lookupErr := lookupDockerCredentials(configPath, host); lookupErr == nil {
		return creds.Username, creds.Password, nil
	}

	return "", "", nil
}
