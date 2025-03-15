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

// TODO: use distribution.Manifest and digest.Digest
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

// TODO: replace by distribution/v3/manifest/schema2.Manifest
// Correspond to a specific image manifest (to push or get,pull)
// The manifest info in this case will not have Annotations, nor Platform
// check https://distribution.github.io/distribution/spec/manifest-v2-2/
type ManifestResponse struct {
	SchemaVersion int            `json:"schemaVersion"`
	MediaType     string         `json:"mediaType"`
	Config        ManifestInfo   `json:"config"`
	Layers        []ManifestInfo `json:"layers"`
}

type Rootfs struct {
	Type    string   `json:"type"`
	DiffIds []string `json:"diff_ids"`
}

type InspectInfo struct {
	Name         string            `json:"Name"`
	Digest       string            `json:"Digest"`
	Tag          Tag               `json:"Tag"`
	RepoTags     []string          `json:"RepoTags"`
	Created      string            `json:"Created"`
	Labels       map[string]string `json:"Labels"`
	Architecture string            `json:"Architecture"`
	Os           string            `json:"Os"`
	Layers       []string          `json:"Layers"`
	LayersData   []LayerData       `json:"LayersData"`
	Env          []string          `json:"Env"`
}

type Tag struct {
	Digest string `json:"Digest"`
	Tag    string `json:"Tag"`
}

type LayerData struct {
	MIMEType string `json:"MIMEType"`
	Digest   string `json:"Digest"`
	Size     int    `json:"Size"`
}

// when doing `inspect`, get blobs/<>/digest content-type config.v1+json (oci),
// or container.image.v1+json for img manifest v2
type ConfigInfo struct {
	Architecture string `json:"architecture"`
	Author       string `json:"author,omitempty"`
	Config       struct {
		User         string            `json:"User,omitempty"`
		ExposedPorts map[string]any    `json:"ExposedPorts,omitempty"`
		Env          []string          `json:"Env,omitempty"`
		Entrypoint   []string          `json:"Entrypoint,omitempty"`
		Cmd          []string          `json:"Cmd,omitempty"`
		WorkingDir   string            `json:"WorkingDir,omitempty"`
		Labels       map[string]string `json:"Labels,omitempty"`
		ArgsEscaped  bool              `json:"ArgsEscaped,omitempty"`
		Shell        []string          `json:"Shell,omitempty"`
	} `json:"config"`
	Created string `json:"created"`
	History []struct {
		Created    string `json:"created"`
		CreatedBy  string `json:"created_by,omitempty"`
		Comment    string `json:"comment,omitempty"`
		EmptyLayer bool   `json:"empty_layer,omitempty"`
	} `json:"history"`
	OS     string `json:"os"`
	RootFS struct {
		Type    string   `json:"type"`
		DiffIDs []string `json:"diff_ids"`
	} `json:"rootfs"`
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
	h := r.GetCustomHeader(fmt.Sprintf(`%s, %s`, MIME_OCI_INDEX, MIME_V2_LIST))
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

// ⚠️  might only be ok with content-type config, use with care or prefer configinfo
// func (r *Registry) GetBlobs(repository, reference, mediaType string) (BlobsResponse, http.Header, error) {
// 	u := fmt.Sprintf(r.baseUrl+blobsPath, repository, reference)
// 	h := r.GetCustomHeader(mediaType)
// 	response, respHeaders, err := HttpDo[BlobsResponse](r.httpClient, http.MethodGet, u, h, nil)
// 	if err != nil {
// 		return response, respHeaders, fmt.Errorf("error getting %s %s blobs:\n%v", repository, reference, err)
// 	}
// 	return response, respHeaders, nil
// }

func (r *Registry) ConfigInfo(repository, reference, mediaType string) (ConfigInfo, http.Header, error) {
	if mediaType != MIME_V2_CONFIG && mediaType != MIME_OCI_CONFIG {
		return ConfigInfo{}, nil, fmt.Errorf("unexpected media type %s, wants %s or %s", mediaType, MIME_V2_CONFIG, MIME_OCI_CONFIG)
	}
	u := fmt.Sprintf(r.baseUrl+blobsPath, repository, reference)
	h := r.GetCustomHeader(mediaType)
	response, respHeaders, err := HttpDo[ConfigInfo](r.httpClient, http.MethodGet, u, h, nil)
	if err != nil {
		return response, respHeaders, fmt.Errorf("error getting %s %s blobs:\n%v", repository, reference, err)
	}
	return response, respHeaders, nil
}

// get images config info (inspect) from name and tag.
// returns a list of inspect info, for each manifests found for the tag
// inspect info is based on config info + additional image info
func (r *Registry) InspectCustom(name, tag string) ([]InspectInfo, error) {
	var inspectInfos []InspectInfo
	tagsResp, _, err := r.GetTags(name)
	if err != nil {
		return inspectInfos, fmt.Errorf("error getting tags for %s: %v", name, err)
	}

	manifestsResp, respHeaders, err := r.GetManifests(name, tag)
	if err != nil {
		return inspectInfos, fmt.Errorf("error getting manifests for %s %s: %v", name, tag, err)
	}

	tagDigest := respHeaders.Get("Docker-Content-Digest")
	if tagDigest == "" {
		return inspectInfos, fmt.Errorf("no digest found for %s %s", name, tag)
	}

	tagInfo := Tag{
		Digest: tagDigest,
		Tag:    tag,
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

			configResp, _, err := r.ConfigInfo(name, digest, mediaType)
			if err != nil {
				resChan <- result{err: fmt.Errorf("error getting blobs for %s %s: %v", name, digest, err)}
				return
			}

			info := NewInspectInfo(name, tagInfo, digest, mediaType, configResp, manifestResp, tagsResp)
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

	if len(errs) > 0 {
		return inspectInfos, fmt.Errorf("some inspections failed: %v", errs)
	}

	return inspectInfos, nil
}

func NewInspectInfo(name string, tagInfo Tag, digest, mediaType string, configResp ConfigInfo, manifestResp ManifestResponse, tagsResp TagsResponse) InspectInfo {
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
		Tag:          tagInfo,
		RepoTags:     tagsResp.Tags,
		Created:      configResp.Created,
		Labels:       configResp.Config.Labels,
		Architecture: configResp.Architecture,
		Os:           configResp.OS,
		Layers:       layers,
		LayersData:   layersData,
		Env:          configResp.Config.Env,
	}
}

// returns the raw config info (blobs/<>/digest content-type MIME_V2_CONFIG or MIME_OCI_CONFIG)
func (r *Registry) Inspect(name, tag string) ([]ConfigInfo, error) {
	var inspectInfos []ConfigInfo

	manifestsResp, respHeaders, err := r.GetManifests(name, tag)
	if err != nil {
		return inspectInfos, fmt.Errorf("error getting manifests for %s %s: %v", name, tag, err)
	}

	tagDigest := respHeaders.Get("Docker-Content-Digest")
	if tagDigest == "" {
		return inspectInfos, fmt.Errorf("no digest found for %s %s", name, tag)
	}

	type result struct {
		info ConfigInfo
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

			info, _, err := r.ConfigInfo(name, digest, mediaType)
			if err != nil {
				resChan <- result{err: fmt.Errorf("error getting configInfo for %s %s: %v", name, digest, err)}
				return
			}

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

	if len(errs) > 0 {
		return inspectInfos, fmt.Errorf("some inspections failed: %v", errs)
	}

	return inspectInfos, nil
}

//	func (r *Registry) BuildImageManifest() {
//	  m := &schema2.Manifest{}
//	  m.SchemaVersion = 2
//	  m.Config = distribution.Descriptor{
//	  	MediaType:   "application/vnd.oci.image.config.v1+json",
//	  	Size:        os.Stat(file),
//	  	Digest:      digest.FromBytes([]byte(full_manifest)),
//	  	URLs:        []string{},
//	  	Annotations: map[string]string{},
//	  	Platform:    &v1.Platform{},
//	  }
//
// }
func (r *Registry) UploadImage() {
}

// TODO: Upload stuff
// 1) build the image manifest : ❓ how to : needs to build the digest algo
//  this is schema2.Manifest or ManifestResponse
// 2) push individual layers
// 2.a) POST /v2/<name>/blobs/uploads/
//  check layer existence with : HEAD /v2/<name>/blobs/<digest> : 200OK means exists (no body)
// 2.b) is POST ok, get 202 accepted + header Location : /v2/<name>/blobs/uploads/<uuid>
// 2.c) upload monolitic
// PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
// Content-Length: <size of layer>
// Content-Type: application/octet-stream
// <Layer Binary Data>
// 2.c bis) upload chunck
//  PATCH /v2/<name>/blobs/uploads/<uuid>
//  Content-Length: <size of chunck>
//  Content-Range: <start>-<end>
//  Content-Type: application/octet-stream
//  < Layer Chunck Binary Data>
//  When chunck accepted, get 202 with header Range: bytes=0-<offset>
// 2.c bis part 2) upload the signed manifest
// PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
// Content-Length: <size of chunck>
// Content-Range: <start of range>-<end of range>
// Content-Type: application/octet-stream
// < Last Layer Chunch Binary Data >
// Get 201 Created if OK
// 4) when all layers are uploaded, upload image manifest
// PUT /v2/<name>/manifests/<reference>
// Content-Type: <manifest media type>
// {
//    "name": <name>,
//    "tag": <tag>,
//    "fsLayers": [
//        {
//            "blobSum": <digest>
//        },
//        ...
//    ],
//    "history": <v1 images>,
//    "signature": <JWS>,
//    ...
// }

// Cancel upload
// DELETE /v2/<name>/blobs/uploads/<uuid>
//

// TODO: Finish deleting stuff
func (r *Registry) DeleteTag(repository, tag, mediaType string) (bool, http.Header, error) {
	u := fmt.Sprintf(r.baseUrl+tagsPath, repository)
	h := r.GetCustomHeader(mediaType)
	response, respHeaders, err := HttpDo[bool](r.httpClient, http.MethodDelete, u, h, nil)
	if err != nil {
		return response, respHeaders, fmt.Errorf("error deleting %s tag:\n%v", tag, err)
	}
	return response, respHeaders, nil
}

// 🔥 If a layer is deleted which is referenced by a manifest in the registry, then the complete images will not be resolvable.
func (r *Registry) DeleteLayer(name, digest, mediaType string) (bool, http.Header, error) {
	u := fmt.Sprintf(r.baseUrl+blobsPath, digest)
	h := r.GetCustomHeader(mediaType)
	response, respHeaders, err := HttpDo[bool](r.httpClient, http.MethodDelete, u, h, nil)
	if err != nil {
		return response, respHeaders, fmt.Errorf("error deleting %s blob:\n%v", digest, err)
	}
	return response, respHeaders, nil
}
