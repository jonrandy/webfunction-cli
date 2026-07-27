package cmd

import (
	"fmt"
	"os"
)

// Command is the interface every subcommand must implement.
type Command interface {
	// Name is the word the user types to invoke this command, e.g. "codegen".
	Name() string

	// Summary is a single short line shown in the general help listing.
	Summary() string

	// Usage is the longer, detailed help shown for "wfn help <command>".
	Usage() string

	// Run executes the command with the remaining args (after the command
	// name has been stripped off).
	Run(args []string) error
}

// registry holds every known command, keyed by name.
var registry = map[string]Command{}

// Register adds a command to the registry. Commands call this from an
// init() in their own file, so adding a new command is just: create the
// file, implement Command, and call Register in init().
func Register(c Command) {
	if _, exists := registry[c.Name()]; exists {
		panic(fmt.Sprintf("wfn: command %q registered twice", c.Name()))
	}
	registry[c.Name()] = c
}

// Lookup returns the command with the given name, if any.
func Lookup(name string) (Command, bool) {
	c, ok := registry[name]
	return c, ok
}

// All returns every registered command.
func All() map[string]Command {
	return registry
}

// Fail prints an error to stderr and exits with a non-zero status. Commands
// can use this for unrecoverable input errors instead of returning an error
// they'd have to format themselves.
func Fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "wfn: "+format+"\n", args...)
	os.Exit(1)
}
