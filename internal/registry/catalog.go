package registry

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
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
func (r *RegistryClient) Catalog() ([]string, error) {
	repositories := make([]string, 0)
	// number of repositories to get per request
	n := "100"
	// used to retrieve last repository name from Link response header
	last := ""
	// regex to parse Link response header
	re := regexp.MustCompile(`<([^>]+)>`)
	u := r.baseUrl + catalogPath

	// paginate flow
	for {
		params := map[string]string{"n": n, "last": last}
		response, h, err := HttpDo[CatalogResponse](r.httpClient, http.MethodGet, u, r.httpHeaders, params)
		if err != nil {
			return repositories, fmt.Errorf("error getting catalog: %v", err)
		}
		repositories = append(repositories, response.Repositories...)

		respLink := h.Get("Link")

		// 📢 if link header is not provided, we reached end of pagination, exit 🚀
		if respLink == "" {
			return repositories, nil
		}

		// 📢 else continue
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

	}
}
