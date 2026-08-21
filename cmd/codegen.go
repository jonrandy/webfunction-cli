package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/webfunction-protocol/webfunction-go"
	"wfn/jsgen"
	"wfn/phpgen"
)

func init() {
	Register(&CodegenCommand{})
}

// validTargets is the set of languages codegen currently knows how to
// generate. Keep this in sync with whatever the generator actually
// implements.
var validTargets = []string{"java", "go", "php", "js", "csharp"}

// defaultNamespace is --namespace's default, used by any target that
// needs one (currently just php; expected to be shared by other
// class-based targets later, e.g. csharp/java).
const defaultNamespace = "WebFunctionClient"

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
  --target     Target language. One of: ` + fmt.Sprint(validTargets) + `
  --url        URL of the webfunction package to generate a client for
  -o           Output file to write the generated code to

Flags (optional):
  --namespace  Namespace for the generated class (currently: php). Default: ` + defaultNamespace + `

Example:
  wfn codegen --target java --url https://example.com/some-package -o client.java`
}

func (c *CodegenCommand) Run(args []string) error {
	fs := flag.NewFlagSet(c.Name(), flag.ContinueOnError)

	target := fs.String("target", "", "target language ("+fmt.Sprint(validTargets)+")")
	url := fs.String("url", "", "URL of the webfunction package")
	output := fs.String("o", "", "output file name")
	namespace := fs.String("namespace", defaultNamespace, "namespace for the generated class (currently: php)")

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

	client, err := webfunction.FromPackageEndpoint(*url, webfunction.Options{})
	if err != nil {
		return fmt.Errorf("fetching package: %w", err)
	}
	pkg := client.Package()

	name := pkg.Name
	if name == "" {
		name = "(unnamed package)"
	}
	fmt.Printf("Fetched %s (%d endpoint(s)) from %s\n", name, len(pkg.Endpoints), *url)

	var source string
	switch *target {
	case "js":
		source, err = jsgen.Generate(pkg, *url)
		if err != nil {
			return fmt.Errorf("generating js: %w", err)
		}
	case "php":
		source, err = phpgen.Generate(pkg, *url, *namespace)
		if err != nil {
			return fmt.Errorf("generating php: %w", err)
		}
	default:
		// TODO: implement the remaining targets (java, go, csharp).
		fmt.Printf("wfn codegen: target %q not yet implemented\n", *target)
		return nil
	}

	if err := os.WriteFile(*output, []byte(source), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", *output, err)
	}
	fmt.Printf("Wrote %s\n", *output)
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