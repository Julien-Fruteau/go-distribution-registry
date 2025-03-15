package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/julien-fruteau/go/distctl/internal/registry"
)

func main() {
	output := flag.String("output", "json", "output format: json or raw")
	flag.Parse()

	r := registry.NewRegistry()
	repositories, err := r.Catalog()
	if err != nil {
		log.Fatal("FATAL error retrieving repositories: ", err)
		os.Exit(1)
	}

	switch *output {
	case "json":
		jsonData, err := json.Marshal(repositories)
		if err != nil {
			log.Fatal("FATAL: ", err)
			return
		}

		fmt.Fprintln(os.Stdout, string(jsonData))
	case "raw":
		fmt.Println(repositories)
	}
}
