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
	repo := flag.String("repo", "", "repository name")
	format := flag.String("format", "json", "output format: json or raw")

	flag.Parse()

	rs := registry.NewRegistrySvc()

	tags, err := rs.GetTags(*repo)
	if err != nil {
		log.Fatal("error retrieving tags: ", err)
		os.Exit(1)
	}

	switch *format {
	case "json":
		jsonData, err := json.Marshal(tags)
		if err != nil {
			log.Fatal("FATAL: ", err)
			return
		}

		fmt.Fprintln(os.Stdout, string(jsonData))
	case "raw":
		fmt.Println(tags)
	}

}
