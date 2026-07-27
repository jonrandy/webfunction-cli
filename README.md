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
| `codegen` | Generate code from a webfunction package (not yet implemented) |

More commands will be added over time.

## Project layout

```
.
├── main.go       # entry point: parses args, dispatches to a command
└── cmd/
    ├── command.go  # Command interface + registry
    ├── help.go     # the "help" command
    └── codegen.go  # the "codegen" command
```

Each command implements a small `Command` interface (`Name`, `Summary`,
`Usage`, `Run`) and registers itself in an `init()` function. Adding a new
command means adding a new file under `cmd/` — no other files need to change.

## License

MIT — see [LICENSE](LICENSE).
