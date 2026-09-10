# go-distribution-registry

A Go client library and `registry` CLI for browsing a Distribution registry's
repositories, tags, and image configuration. The longer-term goal is to support
cleanup of old tags and image layers.

Registry credentials are read using Docker's standard credential configuration.
For a registry-specific helper, remove its entry from `auths` and configure the
helper in `~/.docker/config.json`:

```json
{
  "credHelpers": {
    "dkr.enercal.nc": "pass"
  }
}
```

Or

```json
{
  "credsStore": "pass"
}
```

This makes the client invoke `docker-credential-pass get` for that registry.
Install the matching `docker-credential-<name>` binary first. A global
`credsStore` is also supported; set `REG_HOST` or pass `--reg <host>` when the
registry cannot be inferred from the config.

Target : clean up repository tag and image layers

Requires Go 1.23 or newer (the module requests toolchain Go 1.23.5).

Build from source:

```sh
git clone https://github.com/Julien-Fruteau/go-distribution-registry.git
cd go-distribution-registry
go build -o registry ./cmd/registry
./registry --help
```

Or install the CLI from a local checkout:

```sh
go install ./cmd/registry
```

The installed executable is named `registry` and goes into `GOBIN`, or
`$(go env GOPATH)/bin` when `GOBIN` is unset. Add that directory to your `PATH`
to use the commands below. For a source build, use `./registry` instead.

## Configuration and authentication

The CLI loads an optional `.env` from the current working directory. Existing
process environment variables take precedence. A missing `.env` produces a
warning but does not stop execution. Copy `.env.tpl` to `.env` if needed.

| Setting                    | Behavior                                                                                                                                 |
| -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| `--reg <host>`             | Global CLI flag; overrides `REG_HOST`. Place it **before** the command.                                                                  |
| `REG_HOST`                 | Registry hostname with optional port, such as `registry.example.com:5000`; omit the scheme and path.                                     |
| `REG_SCHEME`               | Defaults to `http`. Set to `https` for TLS registries; `.env.tpl` uses `https`.                                                          |
| `REG_USER`, `REG_PASSWORD` | If `REG_USER` is set, these override Docker credential lookup. A nonempty username enables HTTP Basic authentication.                    |
| `REG_MIME`                 | Advanced override for the default HTTP `Accept` header (Docker v2 and OCI manifests/indexes). Some operations set their own media types. |

Registry selection uses `--reg`, then `REG_HOST`, then the single registry listed
in `~/.docker/config.json` under `auths` or `credHelpers`. If multiple registries
are listed, select one explicitly. With no listed registries, set `REG_HOST` or
use `--reg`. A global `credsStore` alone does not identify a registry host.

When `REG_USER` is unset, credentials for the selected host are picked up from
`~/.docker/config.json` in this order:

1. The host's `credHelpers` entry.
2. Inline `auths` for the bare host or `https://<host>`.
3. The global `credsStore`.

Configured helpers require the corresponding `docker-credential-<helper>`
executable on `PATH`. If credential lookup fails or finds nothing, the client
tries anonymous access; a protected registry may return an authentication error.
The client sends Basic credentials and does not implement a Bearer-token login
flow. `DOCKER_CONFIG` is not used; the config path is fixed to `~/.docker/config.json`.

For an HTTPS registry with credentials saved by Docker:

```sh
docker login registry.example.com
export REG_SCHEME=https
registry --reg registry.example.com catalog
```

If exactly one registry is saved, both `--reg` and `REG_HOST` can be omitted.
To select a saved registry for subsequent commands:

```sh
export REG_HOST=registry.example.com
export REG_SCHEME=https
registry catalog
```

For an anonymous local HTTP registry:

```sh
REG_SCHEME=http registry --reg localhost:5000 catalog
```

## Usage

These examples describe the available CLI commands.

```text
registry [--reg <host>] <command> [options]
```

| Command    | Arguments                        | Result                                                                                                                                                                                                 |
| ---------- | -------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `catalog`  | None                             | Repository names, following catalog pagination automatically.                                                                                                                                          |
| `tags`     | `<name>`                         | Repository name and its tags.                                                                                                                                                                          |
| `inspect`  | `<name>:<tag>` or `<name> <tag>` | Array of image configurations, including architecture, OS, creation date, and runtime configuration. Supports single-image Docker v2/OCI manifests and multi-platform lists/indexes.                   |
| `tagsDate` | `<name>`                         | Intended to list tag names, architectures, and creation dates; currently fails with an error even after successful collection.                                                                         |
| `matchtag` | `<name>`                         | Finds the highest semantic-version tag, including prereleases such as `x.y.z-qual.n`, with image content matching a floating tag. Its `-r`/`--ref` flag selects the reference tag (default: `stable`). |

All commands accept `-o`/`--output` with `json` (default), `yaml`, or `raw`.
`raw` uses Go's printed value representation, not the original HTTP response.
Use `registry --help` for global help or `registry <command> --help` for command
help. Command help initializes the client first, so it still needs a resolvable
registry host. Although catalog help mentions `--pagination`, that flag is not
registered; do not pass it.

```sh
registry catalog -o yaml
registry tags myteam/backend
registry inspect myteam/backend:latest
registry inspect myteam/backend latest -o yaml
registry --reg registry.example.com inspect registry.example.com/myteam/backend:latest
```

Repository names may include the selected registry's host prefix, which is
stripped before making the request. Including a different host in a repository
name does not switch registries; use `--reg` or `REG_HOST` to select the host.

## Library and development

The library lives at
`github.com/julien-fruteau/go-distribution-registry/external/registry`.
`NewRegistryClient(host)` uses the same environment and Docker configuration
resolution as the CLI; pass an empty host to allow automatic selection. Library
callers must load any `.env` themselves.

Manifest/layer deletion and filesystem gzip-blob discovery are library APIs;
there are no CLI delete, garbage-collection, or filesystem-scan commands.
`WalkFs(root)` returns paths to gzip blob files, not extracted digests.

```sh
# Registry library tests (including local test registries).
go test ./external/registry

# Full suite and compile check.
go test ./...
go build ./...
```

## Features

Implemented features:

- List registry repositories with catalog pagination.
- List repository tags.
- Inspect multi-platform Docker manifest lists and OCI image indexes.
- Inspect single-architecture Docker v2 and OCI images ([#1](https://github.com/Julien-Fruteau/go-distribution-registry/issues/1)).
- Accept both `inspect <name>:<tag>` and `inspect <name> <tag>`.
- Select registries with `--reg` / `REG_HOST` and pick up Docker credentials.
- Allow anonymous access when credentials are unavailable.
- Match floating references to the highest versioned Docker v2, OCI, or legacy Schema 1 tag ([#4](https://github.com/Julien-Fruteau/go-distribution-registry/issues/4)).
- Delete manifests and layers through the library.
- Discover gzip blob file paths in registry filesystem storage.

Planned features and fixes:

- Fix `tagsDate` success handling and output ([#9](https://github.com/Julien-Fruteau/go-distribution-registry/issues/9)).
- Remove or implement the unsupported catalog `--pagination` flag ([#5](https://github.com/Julien-Fruteau/go-distribution-registry/issues/5)).
- Align `.env.tpl` with anonymous registry behavior ([#10](https://github.com/Julien-Fruteau/go-distribution-registry/issues/10)).
- Add a k9s-style TUI with registry contexts and interactive browsing ([#2](https://github.com/Julien-Fruteau/go-distribution-registry/issues/2)).
- Implement retention cleanup while preserving the N most recent tags ([#11](https://github.com/Julien-Fruteau/go-distribution-registry/issues/11)).
- Integrate registry garbage collection into cleanup ([#6](https://github.com/Julien-Fruteau/go-distribution-registry/issues/6)).
- Extract digests from discovered filesystem blob paths ([#8](https://github.com/Julien-Fruteau/go-distribution-registry/issues/8)).
- Implement image upload ([#7](https://github.com/Julien-Fruteau/go-distribution-registry/issues/7)).
