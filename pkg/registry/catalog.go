package registry

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

  httpUtils "git.isi.nc/go/dtb-tool/pkg/http"
)

const (
	catalogPath = "_catalog"
)

type CatalogResponse struct {
	Repositories []string `json:"repositories"`
	// Code    string            `json:"code"`
	// Message string            `json:"message"`
	// Detail  map[string]string `json:"detail"`
}

func (r *Registry) GetCatalog(httpCli *http.Client) (CatalogResponse, error) {
	var catalog CatalogResponse

	repositories := make([]string, 0)
	n := "100"
	last := ""

	// paginate
	for {
		// Create a url.Values map to store query parameters
		params := url.Values{}
		params.Add("n", n)
		params.Add("last", url.QueryEscape(last))

		// Encode the parameters into a query string
		queryString := params.Encode()

		url := r.BaseUrl + catalogPath + "?" + queryString

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return catalog, errors.New(fmt.Sprintf("Error creating request: %v", err))
		}

    req.Header.Set("Accept", "application/vnd.oci.image.index.v1+json")
    req.Header.Add("Authorization", httpUtils.GetBasicAuthHeader(r.Conf.Username, r.Conf.Password))

    resp, err := httpCli.Do(req)
    if err != nil {
      return catalog, errors.New(fmt.Sprintf("Error getting catalog: %v", err))
    }
    defer resp.Body.Close()
	}

	return catalog, nil
}
