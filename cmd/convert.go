package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/webfunction-protocol/webfunction-go"
	"wfn/openapiconvert"
)

func init() {
	Register(&ConvertCommand{})
}

// validConvertTargets is the set of formats convert currently knows how
// to convert a webfunction package to. Keep in sync with the generator
// switch in Run.
var validConvertTargets = []string{"openapi"}

// ConvertCommand implements "wfn convert ...".
type ConvertCommand struct{}

func (c *ConvertCommand) Name() string { return "convert" }

func (c *ConvertCommand) Summary() string {
	return "Convert a webfunction package to another document format"
}

func (c *ConvertCommand) Usage() string {
	return `wfn convert --target <format> --url <url> -o <file>

Converts a webfunction package definition into another document format.

Flags (all required):
  --target     Output format. One of: ` + fmt.Sprint(validConvertTargets) + `
  --url        URL of the webfunction package to convert
  -o           Output file to write the converted document to

Example:
  wfn convert --target openapi --url https://api.myservice.com/merchants -o converted.json`
}

func (c *ConvertCommand) Run(args []string) error {
	fs := flag.NewFlagSet(c.Name(), flag.ContinueOnError)

	target := fs.String("target", "", "output format ("+fmt.Sprint(validConvertTargets)+")")
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

	if !isValidConvertTarget(*target) {
		return fmt.Errorf("invalid --target %q, must be one of: %v", *target, validConvertTargets)
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

	var out string
	switch *target {
	case "openapi":
		out, err = openapiconvert.Generate(pkg)
		if err != nil {
			return fmt.Errorf("converting to openapi: %w", err)
		}
	default:
		// isValidConvertTarget already restricts *target to
		// validConvertTargets, so this is unreachable in practice - kept
		// as a safety net.
		return fmt.Errorf("wfn convert: target %q not yet implemented", *target)
	}

	if err := os.WriteFile(*output, []byte(out), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", *output, err)
	}
	fmt.Printf("Wrote %s\n", *output)
	return nil
}

func isValidConvertTarget(target string) bool {
	for _, t := range validConvertTargets {
		if t == target {
			return true
		}
	}
	return false
}