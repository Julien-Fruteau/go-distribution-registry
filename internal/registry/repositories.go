package registry

import (
	"fmt"
	"net/http"
)

type Repository struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type TagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type ManifestsResponse struct {
	SchemaVersion int            `json:"schemaVersion"`
	MediaType     string         `json:"mediaType"`
	Manifests     []ManifestInfo `json:"manifests"`
}

type ManifestInfo struct {
	MediaType   string            `json:"mediaType"`
	Digest      string            `json:"digest"`
	Size        int               `json:"size"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Platform    Platform          `json:"platform,omitempty"`
}

type Platform struct {
	Architecture string `json:"architecture"`
	Os           string `json:"os"`
}

type ManifestResponse struct {
	SchemaVersion int            `json:"schemaVersion"`
	MediaType     string         `json:"mediaType"`
	Config        ManifestInfo   `json:"config"`
	Layers        []ManifestInfo `json:"layers"`
}

type BlobsResponse struct {
	Architecture string    `json:"architecture"`
	Config       Config    `json:"config"`
	Created      string    `json:"created"`
	History      []History `json:"history"`
	Os           string    `json:"os"`
	Rootfs       Rootfs    `json:"rootfs"`
}

type Config struct {
	User         string            `json:"User"`
	ExposedPorts map[string]any    `json:"ExposedPorts,omitempty"`
	Env          []string          `json:"Env,omitempty"`
	Entrypoint   []string          `json:"Entrypoint,omitempty"`
	Cmd          []string          `json:"Cmd,omitempty"`
	WorkingDir   string            `json:"WorkingDir,omitempty"`
	Labels       map[string]string `json:"Labels,omitempty"`
	ArgsEscaped  bool              `json:"ArgsEscaped,omitempty"`
	Shell        []string          `json:"Shell,omitempty"`
}

type History struct {
	Created    string `json:"created"`
	CreatedBy  string `json:"created_by"`
	Comment    string `json:"comment"`
	EmptyLayer bool   `json:"empty_layer,omitempty"`
}

type Rootfs struct {
	Type    string   `json:"type"`
	DiffIds []string `json:"diff_ids"`
}

type InspectInfo struct {
	Name   string `json:"Name"`
	Digest string `json:"Digest"`
	// RepoTags     []string          `json:"RepoTags"`
	Created      string            `json:"Created"`
	Labels       map[string]string `json:"Labels"`
	Architecture string            `json:"Architecture"`
	Os           string            `json:"Os"`
	Layers       []string          `json:"Layers"`
	LayersData   []LayerData       `json:"LayersData"`
	Env          []string          `json:"Env"`
}

type LayerData struct {
	MIMEType string `json:"MIMEType"`
	Digest   string `json:"Digest"`
	Size     int    `json:"Size"`
}

const (
	tagsPath      = `%s/tags/list`
	manifestsPath = `%s/manifests/%s`
	blobsPath     = `%s/blobs/%s`
)

func (r *Registry) GetTags(repository string) (TagsResponse, http.Header, error) {
	u := fmt.Sprintf(r.baseUrl+tagsPath, repository)
	response, respHeaders, err := HttpDo[TagsResponse](r.httpClient, http.MethodGet, u, r.httpHeaders, nil)
	if err != nil {
		return response, respHeaders, fmt.Errorf("error getting %s tags:\n%v", repository, err)
	}
	return response, respHeaders, nil
}

func (r *Registry) GetManifests(repository, reference string) (ManifestsResponse, http.Header, error) {
	u := fmt.Sprintf(r.baseUrl+manifestsPath, repository, reference)
	h := r.GetCustomHeader(fmt.Sprintf(`%s, %s`, MIME_V1_INDEX, MIME_V2_LIST))
	response, respHeaders, err := HttpDo[ManifestsResponse](r.httpClient, http.MethodGet, u, h, nil)
	if err != nil {
		return response, respHeaders, fmt.Errorf("error getting %s %s manifests:\n%v", repository, reference, err)
	}
	return response, respHeaders, nil
}

func (r *Registry) GetManifest(repository, reference, mediaType string) (ManifestResponse, http.Header, error) {
	u := fmt.Sprintf(r.baseUrl+manifestsPath, repository, reference)
	h := r.GetCustomHeader(mediaType)
	response, respHeaders, err := HttpDo[ManifestResponse](r.httpClient, http.MethodGet, u, h, nil)
	if err != nil {
		return response, respHeaders, fmt.Errorf("error getting %s %s manifest:\n%v", repository, reference, err)
	}
	return response, respHeaders, nil
}

func (r *Registry) GetBlobs(repository, reference, mediaType string) (BlobsResponse, http.Header, error) {
	u := fmt.Sprintf(r.baseUrl+blobsPath, repository, reference)
	h := r.GetCustomHeader(mediaType)
	response, respHeaders, err := HttpDo[BlobsResponse](r.httpClient, http.MethodGet, u, h, nil)
	if err != nil {
		return response, respHeaders, fmt.Errorf("error getting %s %s blobs:\n%v", repository, reference, err)
	}
	return response, respHeaders, nil
}

// WARNING: quetion remains about the 2 firsts 1st mediaType to use,
// list tags : index, and manifest v1 works OR index, and manifests list v2
// manifests tag : index v1
// manifests "platform" digest : use provided mediaType info
// blobs config : use provided mediaType info
// 💡 include all headers when unknown : index v1 + manifests  list v2
// TODO : index  manifest, then image manifest with blobs, then ?!.... to retrieve CreatedAt info
// 1) get manifests/<reference> (ref being a tag)  with application/vnd.oci.image.index.v1+json
// HTTP/1.1 200 OK
// Content-Length: 856
// Content-Type: application/vnd.oci.image.index.v1+json
// Docker-Content-Digest: sha256:66339d5454990689da17ffb99217e333c9d84600344d09ee758cff8f6594cb90
// Docker-Distribution-Api-Version: registry/2.0
// Etag: "sha256:66339d5454990689da17ffb99217e333c9d84600344d09ee758cff8f6594cb90"
// X-Content-Type-Options: nosniff
// Date: Mon, 17 Feb 2025 20:46:53 GMT
//
//	{
//	  "schemaVersion": 2,
//	  "mediaType": "application/vnd.oci.image.index.v1+json",
//	  "manifests": [
//	    {
//	      "mediaType": "application/vnd.oci.image.manifest.v1+json",
//	      "digest": "sha256:de4526735344b87a32698f72069f8da6a11681ff9453f729661eb50d62ca3d17",
//	      "size": 1631,
//	      "platform": {
//	        "architecture": "amd64",
//	        "os": "linux"
//	      }
//	    },
//	    {
//	      "mediaType": "application/vnd.oci.image.manifest.v1+json",
//	      "digest": "sha256:02d363a44f31dd0def7dcccff9fbe3b09b2e806cb44aee3e03baf104d097c5a6",
//	      "size": 566,
//	      "annotations": {
//	        "vnd.docker.reference.digest": "sha256:de4526735344b87a32698f72069f8da6a11681ff9453f729661eb50d62ca3d17",
//	        "vnd.docker.reference.type": "attestation-manifest"
//	      },
//	      "platform": {
//	        "architecture": "unknown",
//	        "os": "unknown"
//	      }
//	    }
//	  ]
//	}%
//
// 2) get manifests sha256:de4526735344b87a32698f72069f8da6a11681ff9453f729661eb50d62ca3d17 with application/vnd.oci.image.manifest.v1+json
// {
// "schemaVersion": 2,
// "mediaType": "application/vnd.oci.image.manifest.v1+json",
//
//	"config": {
//	  "mediaType": "application/vnd.oci.image.config.v1+json",
//	  "digest": "sha256:1649f157365545ac4b8ec167619fb18d2b61f802776e39e46a8156f39762615e",
//	  "size": 11148
//	},
//
// "layers": [
//
//	{
//	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
//	  "digest": "sha256:9d1c7dcd50f5547c998ed553485c4c8ef1bcba72abb1b70c4f7de74572c54278",
//	  "size": 145483495
//	},
//	{
//	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
//	  "digest": "sha256:b9be66bfe7f92b5c42a47c6353d0dfb1f7b9610a9479752228d5f1fe00c100fc",
//	  "size": 2094433
//	},
//	{
//	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
//	  "digest": "sha256:08d8d343d6a4c6fb7033d42667d66a88368e95b0f1ee288621dbaf24149d33ca",
//	  "size": 178
//	},
//	{
//	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
//	  "digest": "sha256:e6c0e3d5828e19ef46e585f10e2af75e11be87f42432301fb92df598d2d2d092",
//	  "size": 477195673
//	},
//	{
//	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
//	  "digest": "sha256:e8ff69f6858575d6e0a8be832b30716e41bce7379acee914fa91da26533a484a",
//	  "size": 477197575
//	},
//	{
//	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
//	  "digest": "sha256:107aba61455803961e2bf3981ee15312ff09177dfa623fe785f4759b63afa9a5",
//	  "size": 7267273
//	},
//	{
//	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
//	  "digest": "sha256:4f4fb700ef54461cfa02571ae0db9a0dc1e0cdb5577484a6d75e68dc38e8acc1",
//	  "size": 32
//	}
//
// ]

// 3) get blobs/sha256:1649f157365545ac4b8ec167619fb18d2b61f802776e39e46a8156f39762615e with application/vnd.oci.image.config.v1+json
// cf blobs.sample.json
// SKOPEO inspect response
// {
//     "Name": "dkr.isi/pole-tech/jasperreports",
//     "Digest": "sha256:66339d5454990689da17ffb99217e333c9d84600344d09ee758cff8f6594cb90",
//     "RepoTags": [
//         "fee3eb3",
//         "7d1e687",
//         "e8a6862",
//         "51f6be1"
//     ],
//     "Created": "2025-02-13T14:01:52.814525715+11:00",
//     "DockerVersion": "",
//     "Labels": {
//         "com.vmware.cp.artifact.flavor": "sha256:c50c90cfd9d12b445b011e6ad529f1ad3daea45c26d20b00732fae3cd71f6a83",
//         "org.opencontainers.image.base.name": "docker.io/bitnami/minideb:bookworm",
//         "org.opencontainers.image.created": "2025-02-11T20:49:58Z",
//         "org.opencontainers.image.description": "Application packaged by Broadcom, Inc.",
//         "org.opencontainers.image.documentation": "https://github.com/bitnami/containers/tree/main/bitnami/tomcat/README.md",
//         "org.opencontainers.image.licenses": "Apache-2.0",
//         "org.opencontainers.image.ref.name": "10.1.34-debian-12-r5",
//         "org.opencontainers.image.source": "https://github.com/bitnami/containers/tree/main/bitnami/tomcat",
//         "org.opencontainers.image.title": "tomcat",
//         "org.opencontainers.image.vendor": "Broadcom, Inc.",
//         "org.opencontainers.image.version": "10.1.34"
//     },
//     "Architecture": "amd64",
//     "Os": "linux",
//     "Layers": [
//         "sha256:9d1c7dcd50f5547c998ed553485c4c8ef1bcba72abb1b70c4f7de74572c54278",
//         "sha256:b9be66bfe7f92b5c42a47c6353d0dfb1f7b9610a9479752228d5f1fe00c100fc",
//         "sha256:08d8d343d6a4c6fb7033d42667d66a88368e95b0f1ee288621dbaf24149d33ca",
//         "sha256:e6c0e3d5828e19ef46e585f10e2af75e11be87f42432301fb92df598d2d2d092",
//         "sha256:e8ff69f6858575d6e0a8be832b30716e41bce7379acee914fa91da26533a484a",
//         "sha256:107aba61455803961e2bf3981ee15312ff09177dfa623fe785f4759b63afa9a5",
//         "sha256:4f4fb700ef54461cfa02571ae0db9a0dc1e0cdb5577484a6d75e68dc38e8acc1"
//     ],
//     "LayersData": [
//         {
//             "MIMEType": "application/vnd.oci.image.layer.v1.tar+gzip",
//             "Digest": "sha256:9d1c7dcd50f5547c998ed553485c4c8ef1bcba72abb1b70c4f7de74572c54278",
//             "Size": 145483495,
//             "Annotations": null
//         },
//         {
//             "MIMEType": "application/vnd.oci.image.layer.v1.tar+gzip",
//             "Digest": "sha256:b9be66bfe7f92b5c42a47c6353d0dfb1f7b9610a9479752228d5f1fe00c100fc",
//             "Size": 2094433,
//             "Annotations": null
//         },
//         {
//             "MIMEType": "application/vnd.oci.image.layer.v1.tar+gzip",
//             "Digest": "sha256:08d8d343d6a4c6fb7033d42667d66a88368e95b0f1ee288621dbaf24149d33ca",
//             "Size": 178,
//             "Annotations": null
//         },
//         {
//             "MIMEType": "application/vnd.oci.image.layer.v1.tar+gzip",
//             "Digest": "sha256:e6c0e3d5828e19ef46e585f10e2af75e11be87f42432301fb92df598d2d2d092",
//             "Size": 477195673,
//             "Annotations": null
//         },
//         {
//             "MIMEType": "application/vnd.oci.image.layer.v1.tar+gzip",
//             "Digest": "sha256:e8ff69f6858575d6e0a8be832b30716e41bce7379acee914fa91da26533a484a",
//             "Size": 477197575,
//             "Annotations": null
//         },
//         {
//             "MIMEType": "application/vnd.oci.image.layer.v1.tar+gzip",
//             "Digest": "sha256:107aba61455803961e2bf3981ee15312ff09177dfa623fe785f4759b63afa9a5",
//             "Size": 7267273,
//             "Annotations": null
//         },
//         {
//             "MIMEType": "application/vnd.oci.image.layer.v1.tar+gzip",
//             "Digest": "sha256:4f4fb700ef54461cfa02571ae0db9a0dc1e0cdb5577484a6d75e68dc38e8acc1",
//             "Size": 32,
//             "Annotations": null
//         }
//     ],
//     "Env": [
//         "PATH=/opt/bitnami/common/bin:/opt/bitnami/java/bin:/opt/bitnami/tomcat/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
//         "HOME=/",
//         "OS_ARCH=amd64",
//         "OS_FLAVOUR=debian-12",
//         "OS_NAME=linux",
//         "APP_VERSION=10.1.34",
//         "BITNAMI_APP_NAME=tomcat",
//         "JAVA_HOME=/opt/bitnami/java",
//         "JASPER_VERSION=8.2.0",
//         "JASPER_HOME=/opt/bitnami/tomcat/webapps/jasperserver",
//         "PGSQL_JDBC_VERSION=42.5.0",
//         "CLICKHOUSE_JDBC_VERSION=0.8.0"
//     ]
// }

// TODO : WIP
// reference is expecting as a tag
// get maniests from tag return a list of manifests for all archs the image has been build for
// will make inspect multi arch, returning a list
func (r *Registry) Inspect(name, reference string) ([]InspectInfo, error) {
	var inspectInfos []InspectInfo
	manifestsResp, respHeaders, err := r.GetManifests(name, reference)
	if err != nil {
		return inspectInfos, fmt.Errorf("error getting manifests for %s %s: %v", name, reference, err)
	}

	digest := respHeaders.Get("Docker-Content-Digest")
	if digest == "" {
		return inspectInfos, fmt.Errorf("no digest found for %s %s", name, reference)
	}

	type result struct {
		info InspectInfo
		err  error
	}
	resChan := make(chan result)

	// Launch goroutines
	for _, manifest := range manifestsResp.Manifests {
		go func(m ManifestInfo) {
			manifestResp, _, err := r.GetManifest(name, m.Digest, m.MediaType)
			if err != nil {
				resChan <- result{err: fmt.Errorf("error getting manifest for %s %s: %v", name, m.Digest, err)}
				return
			}

			mediaType := manifestResp.Config.MediaType
			digest := manifestResp.Config.Digest

			blobsResp, _, err := r.GetBlobs(name, digest, mediaType)
			if err != nil {
				resChan <- result{err: fmt.Errorf("error getting blobs for %s %s: %v", name, digest, err)}
				return
			}

			info := NewInspectInfo(name, digest, mediaType, blobsResp, manifestResp)
			resChan <- result{info: info}
		}(manifest)
	}

	// Collect results
	var errs []error
	for range manifestsResp.Manifests {
		res := <-resChan
		if res.err != nil {
			errs = append(errs, res.err)
			continue
		}
		inspectInfos = append(inspectInfos, res.info)
	}

	// Return error if no successful results
	if len(inspectInfos) == 0 && len(errs) > 0 {
		return nil, fmt.Errorf("all inspections failed: %v", errs)
	}

	return inspectInfos, nil
}

func NewInspectInfo(name, digest, mediaType string, blobsResp BlobsResponse, manifestResp ManifestResponse) InspectInfo {
	// Create layers slice from manifest response
	layers := make([]string, len(manifestResp.Layers))
	layersData := make([]LayerData, len(manifestResp.Layers))

	for i, layer := range manifestResp.Layers {
		layers[i] = layer.Digest
		layersData[i] = LayerData{
			MIMEType: layer.MediaType,
			Digest:   layer.Digest,
			Size:     layer.Size,
		}
	}

	return InspectInfo{
		Name:         name,
		Digest:       digest,
		Created:      blobsResp.Created,
		Labels:       blobsResp.Config.Labels,
		Architecture: blobsResp.Architecture,
		Os:           blobsResp.Os,
		Layers:       layers,
		LayersData:   layersData,
		Env:          blobsResp.Config.Env,
	}
}
