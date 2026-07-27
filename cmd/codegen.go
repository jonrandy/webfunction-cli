package cmd

import (
	"flag"
	"fmt"

	"wfn/webfunction"
)

func init() {
	Register(&CodegenCommand{})
}

// validTargets is the set of languages codegen currently knows how to
// generate. Keep this in sync with whatever the generator actually
// implements.
var validTargets = []string{"java", "go", "php", "js", "csharp"}

// CodegenCommand implements "wfn codegen ...".
type CodegenCommand struct{}

func (c *CodegenCommand) Name() string { return "codegen" }

func (c *CodegenCommand) Summary() string {
	return "Generate code from a webfunction package"
}

func (c *CodegenCommand) Usage() string {
	return `wfn codegen --target <language> --url <url> -o <file>

Generates a client for a webfunction package in the given target language.

Flags (all required):
  --target   Target language. One of: ` + fmt.Sprint(validTargets) + `
  --url      URL of the webfunction package to generate a client for
  -o         Output file to write the generated code to

Example:
  wfn codegen --target java --url https://example.com/some-package -o client.java`
}

func (c *CodegenCommand) Run(args []string) error {
	fs := flag.NewFlagSet(c.Name(), flag.ContinueOnError)

	target := fs.String("target", "", "target language ("+fmt.Sprint(validTargets)+")")
	url := fs.String("url", "", "URL of the webfunction package")
	output := fs.String("o", "", "output file name")

	if err := fs.Parse(args); err != nil {
		return err
	}

	var missing []string
	if *target == "" {
		missing = append(missing, "--target")
	}
	if *url == "" {
		missing = append(missing, "--url")
	}
	if *output == "" {
		missing = append(missing, "-o")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required flag(s): %v\n\n%s", missing, c.Usage())
	}

	if !isValidTarget(*target) {
		return fmt.Errorf("invalid --target %q, must be one of: %v", *target, validTargets)
	}

	pkg, err := webfunction.FetchPackage(*url)
	if err != nil {
		return fmt.Errorf("fetching package: %w", err)
	}

	name := pkg.Name
	if name == "" {
		name = "(unnamed package)"
	}
	fmt.Printf("Fetched %s (%d endpoint(s)) from %s\n", name, len(pkg.Endpoints), *url)

	// TODO: generate code for *target from pkg, writing the result to
	// *output. Not yet implemented - only the JS target is planned so far.
	fmt.Printf("wfn codegen: code generation not yet implemented (target=%s output=%s)\n", *target, *output)
	return nil
}

func isValidTarget(target string) bool {
	for _, t := range validTargets {
		if t == target {
			return true
		}
	}
	return false
}