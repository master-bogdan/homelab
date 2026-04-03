package main

import (
	"flag"
	"log"

	"github.com/swaggo/swag/gen"
)

func main() {
	searchDir := flag.String("search", ".", "directories to parse (comma-separated)")
	mainFile := flag.String("main", "cmd/server/main.go", "main API file")
	outputDir := flag.String("output", "docs", "output directory")
	flag.Parse()

	generator := gen.New()
	err := generator.Build(&gen.Config{
		SearchDir:     *searchDir,
		MainAPIFile:   *mainFile,
		OutputDir:     *outputDir,
		OutputTypes:   []string{"go", "json"},
		ParseInternal: true,
		Excludes:      "docs",
		OverridesFile: gen.DefaultOverridesFile,
	})
	if err != nil {
		log.Fatal(err)
	}
}
