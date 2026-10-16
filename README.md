# webfunction-cli

`wfn` is a command line tool for working with [webfunction](https://webfunction.org) packages.

It's designed to grow a small set of subcommands over time, each focused on
one task.

## Installation

Clone and build locally:

```sh
git clone https://github.com/webfunction-protocol/webfunction-cli.git
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
| `codegen` | Generate a typed client for a webfunction package (targets: java, go, php, js, csharp, ruby, python; `--private` to include private endpoints, arguments, and attributes; `--file` to read the package from a local JSON file, though `--url` is still required) |
| `convert` | Convert a webfunction package to another document format (targets: openapi, postman; `--url` or `--file` for the package; `--private` to include private endpoints, arguments, and attributes; generics are expanded into concrete types) |
| `validate` | Check a webfunction package's structure, flags, and generics usage, writing a report (formats: json, markdown, html; `--url` or `--file` for the package) |

More commands will be added over time.

## Where the package comes from

`convert` and `validate` take exactly one of:

- `--url <url>`: fetch the package from its package endpoint.
- `--file <path>`: read a local JSON file holding the same document the
  endpoint returns. Handy for testing and for packages that aren't deployed yet.

Giving both, or neither, is an error.

`codegen` always needs `--url`, because the generated client fetches the
package from that address at runtime. There, `--file` is optional and only
replaces the fetch at generation time, so you can generate a client for a live
URL from a local (possibly edited) package.

With `--file`, the status line reads `Loaded ... from file:///...` instead of
`Fetched ... from <url>`, and `validate` records that `file://` URL as the
report's `package_url`.

```sh
wfn convert --target openapi --file ./package.json -o openapi.json
wfn validate --file ./package.json -o report.html --format html
wfn codegen --target go --url https://example.com/pkg --file ./package.json -o client.go
```

## Generics

A package can define single-parameter generic objects: an object definition
that uses the bare type `T`, applied with a single-key object such as
`{ "object.list": "object.user" }`. Support in this tool:

- `validate` checks how generics are used: unapplied generics, malformed
  applications, applications inside object definitions, name collisions, and
  so on.
- `convert` expands each application into a concrete type before converting.
  For OpenAPI, `{ "object.list": "object.user" }` becomes a `ListOfUser`
  component. Output for packages that don't use generics is unchanged.
- `convert` also prints a `warning: generics: ...` line on stderr for each
  generics problem the validator finds, then carries on, resolving the
  problem to `any`. It never fails because of them; run `wfn validate` for the
  full report.
- `codegen` does not expand generics yet.

## Project layout

```
.
├── main.go             # entry point: parses args, dispatches to a command
├── cmd/
│   ├── command.go      # Command interface + registry
│   ├── help.go         # the "help" command
│   ├── codegen.go      # the "codegen" command
│   ├── convert.go      # the "convert" command
│   ├── generics_warning.go # generics warnings printed by convert
│   ├── package_source.go   # loads a package from --url or --file, for all commands
│   └── validate.go     # the "validate" command
├── csharpgen/          # codegen --target csharp
├── gogen/              # codegen --target go
├── javagen/            # codegen --target java
├── jsgen/              # codegen --target js
├── phpgen/             # codegen --target php
├── pythongen/          # codegen --target python
├── rubygen/            # codegen --target ruby
├── generics/           # expands generic objects into concrete objects (used by validate and convert)
├── openapiconvert/     # convert --target openapi's package-to-OpenAPI-3.1 conversion
├── postmanconvert/     # convert --target postman's package-to-Postman-collection conversion
├── privatefilter/      # strips private endpoints/arguments/attributes before convert/codegen
├── validate/           # the checks and report rendering behind the "validate" command
└── testdata/           # fixture packages for the validate and convert tests, plus golden converter output
```

Each command implements a small `Command` interface (`Name`, `Summary`,
`Usage`, `Run`) and registers itself in an `init()` function. Adding a new
command means adding a new file under `cmd/` — no other files need to change.

## Testing

```sh
go test ./...
```

Each converter's output for every fixture package is pinned byte for byte by
golden files in `testdata/convert/golden/`. After an intended change to
converter output, regenerate them and review the diff:

```sh
go test ./openapiconvert ./postmanconvert -run TestGoldenOutput -update
```

See `testdata/convert/README.md` for details.

## License

MIT — see [LICENSE](LICENSE).