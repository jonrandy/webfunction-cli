package validate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
)

// ValidFormats is the set of --format values wfn validate accepts.
var ValidFormats = []string{"json", "markdown", "html"}

// Render produces the report in the given format ("json", "markdown", or
// "html"). Returns an error for any other format string.
func (r *Report) Render(format string) (string, error) {
	switch format {
	case "json":
		return r.renderJSON()
	case "markdown":
		return r.renderMarkdown(), nil
	case "html":
		return r.renderHTML()
	default:
		return "", fmt.Errorf("unknown format %q, must be one of: %v", format, ValidFormats)
	}
}

func (r *Report) renderJSON() (string, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling report: %w", err)
	}
	return string(b) + "\n", nil
}

// renderMarkdown is the canonical human-readable format - both the
// standalone "markdown" output and the source that "html" is rendered
// from via goldmark, so the two formats never drift apart in content.
//
// Within each severity band, findings are sub-grouped by subject
// (endpoint/object) so everything about one endpoint reads together,
// rather than being scattered in whatever order the checks happened to
// run. Package-scope findings (there's only ever one "subject" for
// those - the package itself) are listed directly under the severity
// heading with no sub-heading of their own.
func (r *Report) renderMarkdown() string {
	var b strings.Builder

	fmt.Fprintf(&b, "# Validation report: %s\n\n", r.PackageName)
	fmt.Fprintf(&b, "- **Source:** %s\n", r.PackageURL)
	fmt.Fprintf(&b, "- **Generated:** %s\n", r.GeneratedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(&b, "- **Errors:** %d  **Warnings:** %d  **Info:** %d\n\n", r.ErrorCount, r.WarningCount, r.InfoCount)

	if len(r.Findings) == 0 {
		b.WriteString("No findings. Package looks structurally clean by every check this validator currently runs.\n")
		return b.String()
	}

	for _, sev := range []Severity{Error, Warning, Info} {
		var inSeverity []Finding
		for _, f := range r.Findings {
			if f.Severity == sev {
				inSeverity = append(inSeverity, f)
			}
		}
		if len(inSeverity) == 0 {
			continue
		}
		fmt.Fprintf(&b, "## %s (%d)\n\n", severityHeading(sev), len(inSeverity))

		for _, group := range groupByLocation(inSeverity) {
			if group.heading != "" {
				fmt.Fprintf(&b, "### %s\n\n", group.heading)
			}
			for _, f := range group.findings {
				if detail := detailString(f); detail != "" {
					fmt.Fprintf(&b, "- **[%s]** `%s` — %s\n", detail, f.Check, f.Message)
				} else {
					fmt.Fprintf(&b, "- `%s` — %s\n", f.Check, f.Message)
				}
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

// locGroup is one subject's findings within a single severity band.
type locGroup struct {
	heading  string // "" for the package group (no sub-heading rendered)
	findings []Finding
}

// groupByLocation groups findings (already filtered to one severity) by
// subject: package findings first (ungrouped), then endpoints
// alphabetically, then objects alphabetically. Order within a group is
// the order findings were passed in.
func groupByLocation(findings []Finding) []locGroup {
	var pkgFindings []Finding
	endpointOrder := []string{}
	endpointFindings := map[string][]Finding{}
	objectOrder := []string{}
	objectFindings := map[string][]Finding{}

	for _, f := range findings {
		switch f.Scope {
		case "endpoint":
			if _, ok := endpointFindings[f.Subject]; !ok {
				endpointOrder = append(endpointOrder, f.Subject)
			}
			endpointFindings[f.Subject] = append(endpointFindings[f.Subject], f)
		case "object":
			if _, ok := objectFindings[f.Subject]; !ok {
				objectOrder = append(objectOrder, f.Subject)
			}
			objectFindings[f.Subject] = append(objectFindings[f.Subject], f)
		default:
			pkgFindings = append(pkgFindings, f)
		}
	}

	sort.Strings(endpointOrder)
	sort.Strings(objectOrder)

	var groups []locGroup
	if len(pkgFindings) > 0 {
		groups = append(groups, locGroup{findings: pkgFindings})
	}
	for _, name := range endpointOrder {
		groups = append(groups, locGroup{heading: fmt.Sprintf("Endpoint: %s", name), findings: endpointFindings[name]})
	}
	for _, name := range objectOrder {
		groups = append(groups, locGroup{heading: fmt.Sprintf("Object: %s", name), findings: objectFindings[name]})
	}
	return groups
}

// detailString renders just the Part/Field-specific fragment of a
// finding's location (e.g. `argument "status"`), for use inside a
// sub-heading that already names the subject - avoids repeating the
// full `endpoint "x" argument "y"` sentence redundantly.
func detailString(f Finding) string {
	switch f.Part {
	case "returns":
		return "returns"
	case "argument":
		return fmt.Sprintf("argument %q", f.Field)
	case "attribute":
		return fmt.Sprintf("attribute %q", f.Field)
	default:
		return ""
	}
}

func severityHeading(sev Severity) string {
	switch sev {
	case Error:
		return "Errors"
	case Warning:
		return "Warnings"
	default:
		return "Info"
	}
}

// renderHTML converts the markdown report to a self-contained HTML
// document via goldmark, with a small amount of surrounding CSS for
// readability (severity-colored headings). The report content itself is
// authored once as markdown (renderMarkdown) - this only wraps it.
func (r *Report) renderHTML() (string, error) {
	md := r.renderMarkdown()

	var body bytes.Buffer
	if err := goldmark.Convert([]byte(md), &body); err != nil {
		return "", fmt.Errorf("rendering markdown to html: %w", err)
	}

	var out strings.Builder
	out.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<meta charset=\"utf-8\">\n")
	fmt.Fprintf(&out, "<title>Validation report: %s</title>\n", html.EscapeString(r.PackageName))
	out.WriteString(`<style>
body { font-family: -apple-system, sans-serif; max-width: 860px; margin: 2rem auto; padding: 0 1rem; line-height: 1.5; color: #1a1a1a; }
h1 { border-bottom: 2px solid #ddd; padding-bottom: 0.3rem; }
h2 { margin-top: 2rem; }
h3 { margin-top: 1.2rem; margin-bottom: 0.4rem; font-size: 1.05em; color: #444; }
ul { padding-left: 1.2rem; }
li { margin: 0.4rem 0; }
code { background: #f3f4f6; padding: 0.1rem 0.35rem; border-radius: 3px; font-size: 0.9em; }
</style>
</head>
<body>
`)
	out.Write(body.Bytes())
	out.WriteString("\n</body>\n</html>\n")
	return out.String(), nil
}