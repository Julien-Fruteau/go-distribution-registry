package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/julien-fruteau/go/distctl/internal/registry"
)

func main() {
	output := flag.String("output", "json", "output format (json or raw)")
	flag.Parse()

	if *output != "json" && *output != "raw" {
		fmt.Println("invalid output format")
		os.Exit(1)
	}

	repo := flag.Arg(0)
	if repo == "" {
		fmt.Println("repository is required")
		os.Exit(1)
	}

	tag := flag.Arg(1)
	if tag == "" {
		fmt.Println("tag is required")
		os.Exit(1)
	}

	r := registry.NewRegistry()

	inspectInfos, err := r.Inspect(repo, tag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch *output {
	case "json":
		json.NewEncoder(os.Stdout).Encode(inspectInfos)
	case "raw":
		fmt.Println(inspectInfos)
	}
}
