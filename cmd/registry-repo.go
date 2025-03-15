package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"git.isi.nc/go/dtb-tool/internal/client"
)

type RegistryResponse struct {
	Repositories []string `json:"repositories"`
	// Code    string            `json:"code"`
	// Message string            `json:"message"`
	// Detail  map[string]string `json:"detail"`
}

func main() {
	client := client.NewClient()

	// reg := registry.NewRegistry()
	repositories := make([]string, 0)
	n := "100"
	last := ""
	// complete := false

	for {
		// Create a url.Values map to store query parameters
		params := url.Values{}
		params.Add("n", n)
		params.Add("last", url.QueryEscape(last))

		// Encode the parameters into a query string
		queryString := params.Encode()

		// Append the query string to the default URL
		fullURL := reg.BaseUrl + "_catalog" + "?" + queryString

		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Accept", "application/vnd.oci.image.index.v1+json")
		req.Header.Add("Authorization", reg.AuthHeader)

		resp, err := reg.Client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		// fmt.Println("Status Code:", resp.StatusCode)

		if resp.StatusCode == 200 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("error reading response:", err)
				return
			}

			// fmt.Println(body)

			var data RegistryResponse
			err = json.Unmarshal(body, &data)
			if err != nil {
				fmt.Println("error unmarshal response:", err)
				return
			}

			repositories = append(repositories, data.Repositories...)

			// fmt.Println(data)

			respLink := resp.Header.Get("Link")

			// fmt.Println(respLink)

			if respLink == "" {
				// fmt.Println("Stop pagination")
				break
			}

			decoded, err := url.QueryUnescape(respLink)
			if err != nil {
				fmt.Println("Error decoding URL:", err)
				return
			}

			// fmt.Println("Decoded String:", decoded)

			re := regexp.MustCompile(`<([^>]+)>`)

			// Find all matches in the input string
			matches := re.FindAllStringSubmatch(decoded, -1)
			lastUrl := matches[0][1]

			parsedURL, err := url.ParseRequestURI(lastUrl)
			if err != nil {
				fmt.Println("Error parsing URL:", err)
				return
			}

			// Extract query parameters
			queryParams := parsedURL.Query()

			// Access individual parameters
			last = queryParams.Get("last")
			// fmt.Println("Last value:", last)

		}
	}

	// fmt.Printf("repositories: %v\n", repositories)
	jsonData, err := json.Marshal(repositories)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the JSON data as a string
	fmt.Fprintln(os.Stdout, string(jsonData))
}
