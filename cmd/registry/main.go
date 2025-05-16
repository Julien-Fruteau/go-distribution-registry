package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/julien-fruteau/go-distribution-registry/external/registry"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"github.com/joho/godotenv"
)

func main() {
	// Define subcommands
	inspectCmd := pflag.NewFlagSet("inspect", pflag.ContinueOnError)
	catalogCmd := pflag.NewFlagSet("catalog", pflag.ContinueOnError)
	tagsCmd := pflag.NewFlagSet("tags", pflag.ContinueOnError)
	tagsDateCmd := pflag.NewFlagSet("tagsDate", pflag.ContinueOnError)

	var output string
	// Inspect command flags
	inspectCmd.StringVarP(&output, "output", "o", "json", "output format: json, yaml or raw")
	inspectCmd.Usage = printInspectHelp

	// Catalog command flags
	// catalogPagination := catalogCmd.Bool("pagination", false, "enable pagination")
	catalogCmd.StringVarP(&output, "output", "o", "json", "output format: json, yaml or raw")
	catalogCmd.Usage = printCatalogHelp

	// Tag command flags
	tagsCmd.StringVarP(&output, "output", "o", "json", "output format: json, yaml or raw")
	tagsCmd.Usage = printTagHelp

	// tagsDate
	tagsDateCmd.StringVarP(&output, "output", "o", "json", "output format: json, yaml or raw")
	tagsDateCmd.Usage = printTagDateHelp

	if len(os.Args) < 2 {
		fmt.Println("expected 'inspect', 'catalog' or 'tag' subcommands")
		printMainUsage()
		os.Exit(1)
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		printMainUsage()
		os.Exit(0)
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	r := registry.NewRegistryClient()

	switch os.Args[1] {
	case "inspect":
		if err := inspectCmd.Parse(os.Args[2:]); err != nil {
			if err == pflag.ErrHelp {
				os.Exit(0)
			}
			inspectCmd.Usage()
			os.Exit(1)
		}
		if inspectCmd.NArg() != 2 {
			fmt.Println("inspect command requires exactly 2 arguments: name and tag")
			printInspectHelp()
			os.Exit(1)
		}
		name := inspectCmd.Arg(0)
		tag := inspectCmd.Arg(1)

		manifest, err := r.Inspect(name, tag)
		if err != nil {
			log.Fatal("FATAL error inspecting repository: ", err)
		}
		outputResult(manifest, output)

	case "catalog":
		if err := catalogCmd.Parse(os.Args[2:]); err != nil {
			if err == pflag.ErrHelp {
				os.Exit(0)
			}
			catalogCmd.Usage()
			os.Exit(1)
		}
		repositories, err := r.Catalog()
		if err != nil {
			log.Fatal("FATAL error retrieving repositories: ", err)
		}
		// Note: pagination flag is set but not used in this example
		// You'll need to implement pagination in your registry package
		outputResult(repositories, output)

	case "tags":
		if err := tagsCmd.Parse(os.Args[2:]); err != nil {
			if err == pflag.ErrHelp {
				os.Exit(0)
			}
			tagsCmd.Usage()
			os.Exit(1)
		}
		if tagsCmd.NArg() != 1 {
			fmt.Println("tags command requires exactly 1 argument: name")
			os.Exit(1)
		}
		name := tagsCmd.Arg(0)

		tags, _, err := r.GetTags(name)
		if err != nil {
			log.Fatal("FATAL error retrieving tags: ", err)
		}
		outputResult(tags, output)

	case "tagsDate":
		if err := tagsDateCmd.Parse(os.Args[2:]); err != nil {
			if err == pflag.ErrHelp {
				os.Exit(0)
			}
			tagsDateCmd.Usage()
			os.Exit(1)
		}
		if tagsDateCmd.NArg() != 1 {
			fmt.Println("tagsDate command requires exactly 1 argument: name")
			os.Exit(1)
		}
		name := tagsDateCmd.Arg(0)

		repoTagsCreateDate, err := r.GetRepositoryTagsCreationDate(name)
		if err != nil {
			log.Fatal("FATAL error retrieving tags creation date: ", err)
		}
		outputResult(repoTagsCreateDate, output)

	default:
		printMainUsage()
		os.Exit(1)
	}
}

func printMainUsage() {
	fmt.Fprintf(os.Stdout, `Usage: %s <command> [options]

Commands:
  inspect     Inspect a repository tag
  catalog     List all repositories
  tags        List all tags for a repository
  tagsDate    List all tags creation date for a repository

Use "%s <command> --help" for more information about a command.
`, os.Args[0], os.Args[0])
}

func printInspectHelp() {
	fmt.Fprintf(os.Stdout, `Usage: %s inspect [options] <name> <tag>

Inspect a repository tag manifest.

Arguments:
  name        Repository name
  tag         Repository tag

Options:
  -o, --output string   Output format: json, yaml or raw (default "json")
  -h, --help            Help for inspect command
`, os.Args[0])
}

func printCatalogHelp() {
	fmt.Fprintf(os.Stdout, `Usage: %s catalog [options]

List all repositories in the registry.

Options:
  --pagination          Enable pagination
  -o, --output string   Output format: json, yaml or raw (default "json")
  -h, --help            Help for catalog command
`, os.Args[0])
}

func printTagHelp() {
	fmt.Fprintf(os.Stdout, `Usage: %s tags [options] <name>

List all tags for a repository.

Arguments:
  name        Repository name

Options:
  -o, --output string   Output format: json, yaml or raw (default "json")
  -h, --help            Help for tag command
`, os.Args[0])
}

func printTagDateHelp() {
	fmt.Fprintf(os.Stdout, `Usage: %s tagsDate [options] <name>

List all tags creation date for a repository.

Arguments:
  name        Repository name

Options:
  -o, --output string   Output format: json, yaml or raw (default "json")
  -h, --help            Help for tag command
`, os.Args[0])
}

func outputResult(data any, format string) {
	switch format {
	case "json":
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatal("FATAL: ", err)
			return
		}
		fmt.Fprintln(os.Stdout, string(jsonData))
	case "yaml":
		yamlData, err := yaml.Marshal(data)
		if err != nil {
			log.Fatal("FATAL: ", err)
			return
		}
		fmt.Fprintln(os.Stdout, string(yamlData))
	case "raw":
		fmt.Println(data)
	}
}
