package registry

// application/vnd.docker.distribution.manifest.v2+json: New image manifest format (schemaVersion = 2)
// application/vnd.docker.distribution.manifest.list.v2+json: Manifest list, aka “fat manifest”
// application/vnd.docker.container.image.v1+json: Container config JSON
// application/vnd.docker.image.rootfs.diff.tar.gzip: “Layer”, as a gzipped tar
// application/vnd.docker.image.rootfs.foreign.diff.tar.gzip: “Layer”, as a gzipped tar that should never be pushed
// application/vnd.docker.plugin.v1+json: Plugin config JSON

// Media Type	Description
// application/vnd.docker.distribution.manifest.v2+json	Docker Image Manifest v2
// application/vnd.docker.distribution.manifest.list.v2+json	Docker Image Index (multi-architecture images)
// application/vnd.oci.image.manifest.v1+json	OCI Image Manifest
// application/vnd.oci.image.index.v1+json	OCI Image Index (similar to manifest list)
// application/vnd.docker.container.image.v1+json	Docker Image Configuration
// application/vnd.oci.image.config.v1+json	OCI Image Configuration
// application/vnd.docker.plugin.v1+json	Docker Plugin Manifest

// curl -H "Accept: application/vnd.docker.distribution.manifest.v2+json" \
// -H "Accept: application/vnd.oci.image.manifest.v1+json" \
// -H "Accept: application/vnd.docker.distribution.manifest.list.v2+json" \
// -H "Accept: application/vnd.oci.image.index.v1+json" \
// -s -D - \
// "https://<registry>/v2/<repository>/manifests/<tag>"

// https://github.com/opencontainers/image-spec/blob/main/media-types.md
// application/vnd.oci.descriptor.v1+json: Content Descriptor
// application/vnd.oci.layout.header.v1+json: OCI Layout
// application/vnd.oci.image.index.v1+json: Image Index
// application/vnd.oci.image.manifest.v1+json: Image manifest
// application/vnd.oci.image.config.v1+json: Image config
// application/vnd.oci.image.layer.v1.tar: "Layer", as a tar archive
// application/vnd.oci.image.layer.v1.tar+gzip: "Layer", as a tar archive compressed with gzip
// application/vnd.oci.image.layer.v1.tar+zstd: "Layer", as a tar archive compressed with zstd
// application/vnd.oci.empty.v1+json: Empty for unused descriptors
const (
	MIME_V2_INDEX                 = "application/vnd.docker.distribution.index.v2+json"
	MIME_V2_MANIFEST              = "application/vnd.docker.distribution.manifest.v2+json"
	MIME_V2_LIST                  = "application/vnd.docker.distribution.manifest.list.v2+json"
	MIME_V2_CONTAINER_CONFIG_JSON = "application/vnd.docker.container.image.v1+json"
	MIME_V2_LAYER_GZIP            = "application/vnd.docker.image.rootfs.diff.tar.gzip"
	MIME_V2_PLUGIN_JSON           = "application/vnd.docker.plugin.v1+json"
	MIME_V1_INDEX                 = "application/vnd.oci.image.index.v1+json"
	MIME_V1_MANIFEST              = "application/vnd.oci.image.manifest.v1+json"
	MIME_V1_CONFIG                = "application/vnd.oci.image.config.v1+json"
)
