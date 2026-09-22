// Package validate inspects a parsed webfunction.Package for structural
// issues - dangling references, ambiguous or inconsistent declarations,
// and known wire-format footguns - and reports them as a list of
// Findings, each carrying a severity.
//
// This package only ever reads a *webfunction.Package; it never mutates
// one, and it never talks to the network itself (the caller is expected
// to have already fetched the package via webfunction.FromPackageEndpoint,
// same as codegen/convert).
package validate

import (
	"time"

	"github.com/webfunction-protocol/webfunction-go"
)

// Severity classifies how serious a Finding is.
type Severity string

const (
	// Error findings indicate a real defect in the package - a dangling
	// reference, a name collision that silently masks data, or a flag
	// combination that's internally inconsistent. A report with any
	// Error finding should be treated as "this package has a bug".
	Error Severity = "error"
	// Warning findings indicate something that's valid per the wire
	// grammar but is likely to cause real trouble for consumers (an
	// unconstrained shape, a URL footgun, a choices/values list that
	// doesn't match its own declared type).
	Warning Severity = "warning"
	// Info findings are informational notes - unrecognized-but-maybe-
	// legitimate flags, and known spec ambiguities this validator
	// deliberately does not take a side on. Nothing to necessarily fix.
	Info Severity = "info"
)

// Finding is a single issue (or note) found while validating a package.
type Finding struct {
	// Check is a short, stable, machine-matchable identifier for the
	// check that produced this finding (e.g. "dangling-object-ref").
	// Stable across runs/versions so CI can allowlist or diff on it.
	Check string `json:"check"`
	// Severity is one of Error, Warning, or Info.
	Severity Severity `json:"severity"`
	// Location is a short human-readable path to where the issue was
	// found (e.g. `endpoint "list-people" argument "status"`).
	Location string `json:"location"`
	// Message describes the issue itself.
	Message string `json:"message"`
}

// Report is the full result of validating one package.
type Report struct {
	// PackageName is the package's declared name, or "(unnamed
	// package)" if it didn't declare one - mirrors codegen/convert's
	// own fallback for consistency.
	PackageName string `json:"package_name"`
	// PackageURL is the URL the package was fetched from.
	PackageURL string `json:"package_url"`
	// GeneratedAt is when this report was produced.
	GeneratedAt time.Time `json:"generated_at"`
	// Findings holds every issue/note found, in a stable check-by-check
	// order (not severity-sorted - see Report.BySeverity for that).
	Findings []Finding `json:"findings"`
	// ErrorCount, WarningCount, InfoCount are the number of Findings at
	// each severity - computed once in Run, not recomputed on access.
	ErrorCount   int `json:"error_count"`
	WarningCount int `json:"warning_count"`
	InfoCount    int `json:"info_count"`
}

// Failed reports whether the report contains any Error-severity finding.
// This is what wfn validate's exit code is based on.
func (r *Report) Failed() bool {
	return r.ErrorCount > 0
}

// BySeverity returns a copy of Findings grouped and ordered
// error-then-warning-then-info, stable within each group. Used by the
// report renderers so every format presents findings in the same order
// regardless of which check happened to run first.
func (r *Report) BySeverity() []Finding {
	out := make([]Finding, 0, len(r.Findings))
	for _, sev := range []Severity{Error, Warning, Info} {
		for _, f := range r.Findings {
			if f.Severity == sev {
				out = append(out, f)
			}
		}
	}
	return out
}

// Run validates pkg (fetched from url) and returns the full report.
func Run(pkg *webfunction.Package, url string) *Report {
	name := pkg.Name
	if name == "" {
		name = "(unnamed package)"
	}

	v := &validator{pkg: pkg}
	v.checkBaseURL()
	v.checkDuplicateEndpointNames()
	v.checkDuplicateObjectNames()
	v.checkVersioning()
	v.checkObjectRefs()
	v.checkBareReturns()
	v.checkChoicesAndValues()
	v.checkFlags()
	v.checkSpecAmbiguities()

	r := &Report{
		PackageName: name,
		PackageURL:  url,
		GeneratedAt: time.Now().UTC(),
		Findings:    v.findings,
	}
	for _, f := range r.Findings {
		switch f.Severity {
		case Error:
			r.ErrorCount++
		case Warning:
			r.WarningCount++
		case Info:
			r.InfoCount++
		}
	}
	return r
}

// validator accumulates findings while walking a single package. Not
// safe for concurrent use - Run doesn't need it to be.
type validator struct {
	pkg      *webfunction.Package
	findings []Finding
}

func (v *validator) emit(check string, sev Severity, location, message string) {
	v.findings = append(v.findings, Finding{
		Check:    check,
		Severity: sev,
		Location: location,
		Message:  message,
	})
}