package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	httpUtils "git.isi.nc/go/dtb-tool/pkg/http"
)

const (
	catalogPath = "_catalog"
)

type CatalogResponse struct {
	Repositories []string `json:"repositories"`
}

// ℹ️l LISTING REPOSITORIES ℹ️
//
// Base request: GET /v2/_catalog
//
// Starting paginated flow request: GET /v2/_catalog?n=<integer>
//
// The response from request looks like:
//
// 200 OK
// Content-Type: application/json
// Link: <<url>?n=<n from the request>&last=<last repository in response>>; rel="next"
//
//	{
//	    "repositories": [
//	        <name>,
//	        ...
//	    ]
//	}
//
// the Link header:
//   - if NOT provided: all results received
//   - if provided: last must be used to get the next pagination
//
// Next : GET /v2/_catalog?n=<n from the request>&last=<last repository value from previous response>
func (r *Registry) Catalog(httpClient *http.Client) ([]string, error) {
	repositories := make([]string, 0)
	// number of repositories to get per request
	n := "100"
	// used to retrieve last repository name from Link response header
	last := ""
	// regex to parse Link response header
	re := regexp.MustCompile(`<([^>]+)>`)

	// paginate flow
	for {
		req, err := httpUtils.GetNewRequest(http.MethodGet, r.BaseUrl+catalogPath, map[string]string{"n": n, "last": last})
		if err != nil {
			return repositories, fmt.Errorf("error creating request: %v", err)
		}

		req.Header.Set("Accept", r.Conf.Mime)
		// req.Header.Set("Accept", "application/vnd.oci.image.index.v1+json")
		req.Header.Set("Authorization", httpUtils.GetBasicAuthHeader(r.Conf.Username, r.Conf.Password))

		resp, err := httpClient.Do(req)
		if err != nil {
			return repositories, fmt.Errorf("error getting catalog: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return repositories, fmt.Errorf("error reading response: %v", err)
		}

		switch resp.StatusCode {
		case http.StatusOK:
			var data CatalogResponse
			err = json.Unmarshal(body, &data)
			if err != nil {
				return repositories, fmt.Errorf("error unmarshal response: %v", err)
			}

			repositories = append(repositories, data.Repositories...)

			// 📢 if link header is not provided, we reached end of pagination, exit
			respLink := resp.Header.Get("Link")
			if respLink == "" {
				return repositories, nil
			}

			// else continue
			decoded, err := url.QueryUnescape(respLink)
			if err != nil {
				return repositories, fmt.Errorf("error decoding url: %v", err)
			}

			// Find all matches in the input string
			matches := re.FindAllStringSubmatch(decoded, -1)
			lastUrl := matches[0][1]

			parsedURL, err := url.ParseRequestURI(lastUrl)
			if err != nil {
				return repositories, fmt.Errorf("error parsing url: %v", err)
			}

			// Extract query parameters
			queryParams := parsedURL.Query()

			// Access individual parameters
			last = queryParams.Get("last")

		default:
			var respErr RegistryError
			err = json.Unmarshal(body, &respErr)
			if err != nil {
				return repositories, fmt.Errorf("%d, error getting catalog: %v", resp.StatusCode, body)
			}
			return repositories, fmt.Errorf("%d, error getting catalog: %v", resp.StatusCode, respErr)
		}

	}
}
