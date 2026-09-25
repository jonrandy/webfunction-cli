package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/webfunction-protocol/webfunction-go"
	"wfn/validate"
)

func init() {
	Register(&ValidateCommand{})
}

// ValidateCommand implements "wfn validate ...".
type ValidateCommand struct{}

func (c *ValidateCommand) Name() string { return "validate" }

func (c *ValidateCommand) Summary() string {
	return "Validate the structure of a webfunction package"
}

func (c *ValidateCommand) Usage() string {
	return `wfn validate --url <url> -o <file> [--format json|markdown|html]

Fetches a webfunction package and checks its structure for dangling
references, inconsistent flags/versions, and other known wire-format
footguns, writing a report to -o.

Flags:
  --url        URL of the webfunction package to validate (required)
  -o           Output file to write the report to (required)
  --format     Report format: json, markdown, or html (default: json)

Exit status: non-zero if the report contains any error-severity finding
(useful for CI), even though the report file is still written either way.
Fetch/write failures also exit non-zero, same as codegen/convert.

Example:
  wfn validate --url https://api.reservepay.com/merchants -o report.json
  wfn validate --url https://api.reservepay.com/merchants -o report.html --format html`
}

func (c *ValidateCommand) Run(args []string) error {
	fs := flag.NewFlagSet(c.Name(), flag.ContinueOnError)

	url := fs.String("url", "", "URL of the webfunction package")
	output := fs.String("o", "", "output file name")
	format := fs.String("format", "json", "report format ("+fmt.Sprint(validate.ValidFormats)+")")

	if err := fs.Parse(args); err != nil {
		return err
	}

	var missing []string
	if *url == "" {
		missing = append(missing, "--url")
	}
	if *output == "" {
		missing = append(missing, "-o")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required flag(s): %v\n\n%s", missing, c.Usage())
	}

	if !isValidFormat(*format) {
		return fmt.Errorf("invalid --format %q, must be one of: %v", *format, validate.ValidFormats)
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

	report := validate.Run(pkg, *url)

	rendered, err := report.Render(*format)
	if err != nil {
		return fmt.Errorf("rendering report: %w", err)
	}

	if err := os.WriteFile(*output, []byte(rendered), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", *output, err)
	}
	fmt.Printf("Wrote %s\n", *output)
	fmt.Printf("%d error(s), %d warning(s), %d info note(s)\n", report.ErrorCount, report.WarningCount, report.InfoCount)

	if report.Failed() {
		return fmt.Errorf("validation found %d error(s) - see %s", report.ErrorCount, *output)
	}
	return nil
}

func isValidFormat(format string) bool {
	for _, f := range validate.ValidFormats {
		if f == format {
			return true
		}
	}
	return false
}