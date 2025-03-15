package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	httpUtils "github.com/julien-fruteau/go/distctl/internal/http"
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
	tagsPath = "%s/tags/list"
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
func (r *Registry) GetTags(httpClient *http.Client, repository string) ([]string, error) {
	var tags []string

	u := fmt.Sprintf(r.BaseUrl+tagsPath, repository)

	req, err := httpUtils.GetNewRequest(http.MethodGet, u, r.HttpHeaders, nil)
	if err != nil {
		return tags, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return tags, fmt.Errorf("error getting tags for %s: %v", repository, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tags, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		var respErr RegistryError
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return tags, fmt.Errorf("%d, error getting tags: %v", resp.StatusCode, string(body))
		}
		return tags, fmt.Errorf("%d, error getting tags: %v", resp.StatusCode, respErr)
	}

	var tr TagsResponse
	err = json.Unmarshal(body, &tr)
	if err != nil {
		return tags, fmt.Errorf("error unmarshal response: %v", err)
	}

	tags = tr.Tags

	return tags, nil
}
