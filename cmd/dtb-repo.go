package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"git.isi.nc/go/dtb-tool/pkg/registry"
)

func main() {
	httpClient := &http.Client{}
	registry := registry.NewRegistry()
	repositories, err := registry.Catalog(httpClient)
	if err != nil {
		println(fmt.Errorf("error retrieving repositories: %v", err))
		os.Exit(1)
	}

  // TODO: choose format :
  // raw string
	// fmt.Println(repositories)

  // json
	jsonData, err := json.Marshal(repositories)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the JSON data as a string
	fmt.Fprintln(os.Stdout, string(jsonData))
}
