# webfunction-cli

`wfn` is a command line tool for working with [webfunction](https://webfunction.org) packages.

It's designed to grow a small set of subcommands over time, each focused on
one task.

## Installation

```sh
go install wfn@latest
```

Or clone and build locally:

```sh
git clone https://github.com/robinclart/webfunction-cli.git
cd webfunction-cli
go build -o wfn .
```

## Usage

```sh
wfn                 # show general help
wfn help            # same as above
wfn help <command>  # show detailed help for a specific command
wfn <command> [arguments]
```

## Commands

| Command   | Description                                              |
|-----------|------------------------------------------------------------|
| `help`    | Show general help, or help for a specific command           |
| `codegen` | Generate a typed client for a webfunction package (targets: java, go, php, js, csharp, ruby, python; `--private` to include private endpoints, arguments, and attributes) |
| `convert` | Convert a webfunction package to another document format (targets: openapi, postman; `--private` to include private endpoints, arguments, and attributes) |
| `validate` | Check a webfunction package's structure and flags, writing a report (formats: json, markdown, html) |

More commands will be added over time.

## Project layout

```
.
├── main.go             # entry point: parses args, dispatches to a command
├── cmd/
│   ├── command.go      # Command interface + registry
│   ├── help.go         # the "help" command
│   ├── codegen.go      # the "codegen" command
│   ├── convert.go      # the "convert" command
│   └── validate.go     # the "validate" command
├── csharpgen/          # codegen --target csharp
├── gogen/              # codegen --target go
├── javagen/            # codegen --target java
├── jsgen/              # codegen --target js
├── phpgen/             # codegen --target php
├── pythongen/          # codegen --target python
├── rubygen/            # codegen --target ruby
├── openapiconvert/     # convert --target openapi's package-to-OpenAPI-3.1 conversion
├── postmanconvert/     # convert --target postman's package-to-Postman-collection conversion
├── privatefilter/      # strips private arguments/attributes before convert/codegen
├── validate/           # the checks and report rendering behind the "validate" command
└── testdata/           # fixture packages for validate and convert/codegen checks
```

Each command implements a small `Command` interface (`Name`, `Summary`,
`Usage`, `Run`) and registers itself in an `init()` function. Adding a new
command means adding a new file under `cmd/` — no other files need to change.

## License

MIT — see [LICENSE](LICENSE).