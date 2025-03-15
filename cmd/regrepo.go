package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"git.isi.nc/go/dtb-tool/pkg/client"
)

func main() {
	output := flag.String("output", "json", "output format: json or raw")
	flag.Parse()

  cli := client.NewRegistryClient()
  repositories, err := cli.GetCatalog()

	if err != nil {
		println(fmt.Errorf("error retrieving repositories: %v", err))
		os.Exit(1)
	}

	switch *output {
	case "json":
		jsonData, err := json.Marshal(repositories)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Fprintln(os.Stdout, string(jsonData))
	case "raw":
		fmt.Println(repositories)
	}
}
