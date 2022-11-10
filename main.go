package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
)

func main() {
	var (
		schemaPath        = flag.String("schema", "", "Path to the JSON Schema")
		templatePath      = flag.String("template", "", "Path to a template")
		baseMarkdownLevel = flag.Int("markdownlevel", 1, "Base Markdown level for headings")
	)

	flag.Parse()

	if (*baseMarkdownLevel > 3) || (*baseMarkdownLevel < 1) {
		fmt.Println("Markdown heading level should be minimum 1, maximum 3.")
		os.Exit(1)
	}

	if *schemaPath == "" {
		fmt.Println("no path to schema")
		os.Exit(1)
	}

	f, err := os.Open(*schemaPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	schema, err := newSchema(f, ".")
	if err != nil {
		log.Fatal(err)
	}

	tpl, err := getOrDefaultTemplate(*templatePath, baseMarkdownLevel)
	if err != nil {
		log.Fatal(err)
	}

	if err := tpl.Execute(os.Stdout, schema); err != nil {
		log.Fatal(err)
	}
}

func getOrDefaultTemplate(path string, baseMarkdownLevel *int) (*template.Template, error) {
	if path == "" {
		return template.New("docs").Parse(fmt.Sprintf(`{{ .Markdown %d %d }}`, *baseMarkdownLevel, *baseMarkdownLevel+2))
	}
	return template.ParseFiles(path)
}
