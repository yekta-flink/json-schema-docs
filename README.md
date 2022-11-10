# json-schema-docs

A simple JSON Schema to Markdown generator, forked from `marcusolsson/json-schema-docs` on v0.2.1.

This generator doesn't attempt to support the full JSON Schema specification. Instead, it's designed with the rationale that most people are only using a subset of the spec.

## Install

```bash
go install github.com/goflink/json-schema-docs
```

## Run

```bash
json-schema-docs -schema ./user.schema.json > user.md
```

To create the Markdown with a higher base level:
```bash
json-schema-docs -schema ./user.schema.json -markdownlevel 3 > user.md
```

To use a template when generating the Markdown:

```bash
json-schema-docs -schema ./user.schema.json -template user.md.tpl > user.md
```

`template` is the path to a file containing a Go template, such as this one:

```bash
+++
title = "{{ .Title }}"
description = "{{ .Description }}"
+++

# API reference

This is the reference documentation for an API.

{{ .Markdown 2 4}}
```

The argument to `.Markdown` is the heading level you want the docs to start at, and the maximum heading depth for generation.

## License

[Apache 2.0 License](LICENSE)
