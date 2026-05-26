package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

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

	// Parse global flags before the subcommand. SetInterspersed(false) stops
	// parsing at the first non-flag argument (the subcommand).
	globalFlags := pflag.NewFlagSet("global", pflag.ContinueOnError)
	globalFlags.SetInterspersed(false)
	var regHost string
	globalFlags.StringVar(&regHost, "reg", "", "registry host (overrides REG_HOST env var)")
	if err := globalFlags.Parse(os.Args[1:]); err != nil && err != pflag.ErrHelp {
		fmt.Fprintf(os.Stderr, "error parsing global flags: %v\n", err)
		printMainUsage()
		os.Exit(1)
	}
	remainingArgs := globalFlags.Args()

	if len(remainingArgs) == 0 {
		fmt.Println("expected 'inspect', 'catalog' or 'tag' subcommands")
		printMainUsage()
		os.Exit(1)
	}

	if err := godotenv.Load(); err != nil {
		log.Printf("warning: .env not loaded (%v), falling back to docker credentials", err)
	}

	r, err := registry.NewRegistryClient(regHost)
	if err != nil {
		log.Fatal(err)
	}

	switch remainingArgs[0] {
	case "inspect":
		if err := inspectCmd.Parse(remainingArgs[1:]); err != nil {
			if err == pflag.ErrHelp {
				os.Exit(0)
			}
			inspectCmd.Usage()
			os.Exit(1)
		}
		var name, tag string
		switch inspectCmd.NArg() {
		case 1:
			arg := inspectCmd.Arg(0)
			idx := strings.LastIndex(arg, ":")
			if idx < 0 {
				fmt.Println("inspect: single argument must be in the format <name>:<tag>")
				printInspectHelp()
				os.Exit(1)
			}
			name, tag = arg[:idx], arg[idx+1:]
		case 2:
			name = inspectCmd.Arg(0)
			tag = inspectCmd.Arg(1)
		default:
			fmt.Println("inspect command requires 1 argument <name>:<tag> or 2 arguments <name> <tag>")
			printInspectHelp()
			os.Exit(1)
		}
		name = r.NormalizeName(name)

		manifest, err := r.Inspect(name, tag)
		if err != nil {
			log.Fatal("FATAL error inspecting repository: ", err)
		}
		outputResult(manifest, output)

	case "catalog":
		if err := catalogCmd.Parse(remainingArgs[1:]); err != nil {
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
		if err := tagsCmd.Parse(remainingArgs[1:]); err != nil {
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
		name := r.NormalizeName(tagsCmd.Arg(0))

		tags, _, err := r.GetTags(name)
		if err != nil {
			log.Fatal("FATAL error retrieving tags: ", err)
		}
		outputResult(tags, output)

	case "tagsDate":
		if err := tagsDateCmd.Parse(remainingArgs[1:]); err != nil {
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
		name := r.NormalizeName(tagsDateCmd.Arg(0))

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
	fmt.Fprintf(os.Stdout, `Usage: %s [--reg <host>] <command> [options]

Global Options:
  --reg string   Registry host (overrides REG_HOST env var; required when multiple registries exist in docker config)

Commands:
  inspect     Inspect a repository tag
  catalog     List all repositories
  tags        List all tags for a repository
  tagsDate    List all tags creation date for a repository

Use "%s <command> --help" for more information about a command.
`, os.Args[0], os.Args[0])
}

func printInspectHelp() {
	fmt.Fprintf(os.Stdout, `Usage: %s inspect [options] <name>:<tag>
       %s inspect [options] <name> <tag>

Inspect a repository tag manifest.

Arguments:
  name        Repository name (e.g. myrepo/backend)
  tag         Repository tag  (e.g. latest)

  A single argument in the form <name>:<tag> is also accepted.

Options:
  -o, --output string   Output format: json, yaml or raw (default "json")
  -h, --help            Help for inspect command
`, os.Args[0], os.Args[0])
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
