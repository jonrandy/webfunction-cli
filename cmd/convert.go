package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/webfunction-protocol/webfunction-go"
	"wfn/openapiconvert"
	"wfn/postmanconvert"
)

func init() {
	Register(&ConvertCommand{})
}

// validConvertTargets is the set of formats convert currently knows how
// to convert a webfunction package to. Keep in sync with the generator
// switch in Run.
var validConvertTargets = []string{"openapi", "postman"}

// ConvertCommand implements "wfn convert ...".
type ConvertCommand struct{}

func (c *ConvertCommand) Name() string { return "convert" }

func (c *ConvertCommand) Summary() string {
	return "Convert a webfunction package to another document format"
}

func (c *ConvertCommand) Usage() string {
	return `wfn convert --target <format> --url <url> -o <file>

Converts a webfunction package definition into another document format.

Flags:
  --target     Output format. One of: ` + fmt.Sprint(validConvertTargets) + ` (required)
  --url        URL of the webfunction package to convert (required)
  -o           Output file to write the converted document to (required)
  --private    Include endpoints flagged private in the output (default: excluded)

Example:
  wfn convert --target openapi --url https://api.reservepay.com/merchants -o converted.json
  wfn convert --target postman --url https://api.reservepay.com/merchants -o merchants.postman_collection.json --private`
}

func (c *ConvertCommand) Run(args []string) error {
	fs := flag.NewFlagSet(c.Name(), flag.ContinueOnError)

	target := fs.String("target", "", "output format ("+fmt.Sprint(validConvertTargets)+")")
	url := fs.String("url", "", "URL of the webfunction package")
	output := fs.String("o", "", "output file name")
	private := fs.Bool("private", false, "include endpoints flagged private in the output")

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
		out, err = openapiconvert.Generate(pkg, *private)
		if err != nil {
			return fmt.Errorf("converting to openapi: %w", err)
		}
	case "postman":
		out, err = postmanconvert.Generate(pkg, *private)
		if err != nil {
			return fmt.Errorf("converting to postman: %w", err)
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