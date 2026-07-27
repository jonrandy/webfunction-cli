package cmd

import "fmt"

func init() {
	Register(&CodegenCommand{})
}

// CodegenCommand implements "wfn codegen ...".
//
// TODO: this is a stub. Real behaviour (fetching a webfunction Package
// manifest and generating a client, presumably) to be filled in once the
// spec for this command is decided.
type CodegenCommand struct{}

func (c *CodegenCommand) Name() string { return "codegen" }

func (c *CodegenCommand) Summary() string {
	return "Generate code from a webfunction package (not yet implemented)"
}

func (c *CodegenCommand) Usage() string {
	return `wfn codegen [arguments]

Generates code from a webfunction package.

This command is not yet implemented.`
}

func (c *CodegenCommand) Run(args []string) error {
	fmt.Println("wfn codegen: not yet implemented")
	return nil
}
