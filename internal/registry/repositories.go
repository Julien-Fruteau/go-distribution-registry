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

// get images config info (inspect) from name and tag.
// returns a list of inspect info, for each manifests found for the tag
func (r *Registry) Inspect(name, tag string) ([]InspectInfo, error) {
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

			blobsResp, _, err := r.GetBlobs(name, digest, mediaType)
			if err != nil {
				resChan <- result{err: fmt.Errorf("error getting blobs for %s %s: %v", name, digest, err)}
				return
			}

			info := NewInspectInfo(name, tagInfo, digest, mediaType, blobsResp, manifestResp, tagsResp)
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

func NewInspectInfo(name string, tagInfo Tag, digest, mediaType string, blobsResp BlobsResponse, manifestResp ManifestResponse, tagsResp TagsResponse) InspectInfo {
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
		Created:      blobsResp.Created,
		Labels:       blobsResp.Config.Labels,
		Architecture: blobsResp.Architecture,
		Os:           blobsResp.Os,
		Layers:       layers,
		LayersData:   layersData,
		Env:          blobsResp.Config.Env,
	}
}
