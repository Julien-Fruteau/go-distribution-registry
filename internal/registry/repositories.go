package registry

import (
	"encoding/json"
	"fmt"
	"io"
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

const (
	tagsPath     = "%s/tags/list"
	manifestPath = "%s/manifests/%s"
)

// Response Get Tags
// -----------------
//
// 200 OK
// Content-Type: application/json
//
//	{
//	    "name": <name>,
//	    "tags": [
//	        <tag>,
//	        ...
//	    ]
//	}
func (r *Registry) GetTags(repository string) ([]string, error) {
	var tags []string
	u := fmt.Sprintf(r.baseUrl+tagsPath, repository)
	response, _, err := HttpDo[TagsResponse](r.httpClient, http.MethodGet, u, r.httpHeaders, nil)
	if err != nil {
		return tags, fmt.Errorf("error gettings tags: %v", err)
	}
	tags = response.Tags
	return tags, nil
}

// TODO : index  manifest, then image manifest with blobs, then ?!.... to retrieve CreatedAt info
// GetManifest get an image manifest.
// The name and reference parameter identify the image and are required. The reference may include a tag or digest.
func (r *Registry) GetManifest(name, reference string) (string, error) {
	u := fmt.Sprintf(r.baseUrl+manifestPath, name, reference)

	req, err := GetNewRequest(http.MethodGet, u, r.httpHeaders, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error getting manifest for %s %s: %v", name, reference, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		var respErr RegistryError
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return "", fmt.Errorf("%d, error getting manifest: %v", resp.StatusCode, string(body))
		}
		return "", fmt.Errorf("%d, error getting manifest: %v", resp.StatusCode, respErr)
	}

	return "", nil
}
