package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/webfunction-protocol/webfunction-go"
	"wfn/generics"
	"wfn/validate"
)

// warnWriter is where generics warnings are written. It is a variable so
// tests can capture them.
var warnWriter io.Writer = os.Stderr

// warnGenerics runs the validator on pkg and writes a line to w for each
// finding that concerns generics, then a closing note. It writes nothing
// when there are none, so packages that don't use generics see no change.
//
// The converters and code generators never fail on these problems:
// generics.Expand falls back to "any" for an unapplied or malformed
// generic. These warnings are how a user finds out that happened.
//
// Every generic-* finding is reported. dangling-object-ref is reported only
// when the package defines a generic, because in a package without one it
// is an ordinary problem the converters have never warned about.
func warnGenerics(w io.Writer, pkg *webfunction.Package, url string) {
	definesGeneric := false
	for i := range pkg.Objects {
		if generics.IsGeneric(&pkg.Objects[i]) {
			definesGeneric = true
			break
		}
	}

	count := 0
	for _, f := range validate.Run(pkg, url).BySeverity() {
		isGenerics := strings.HasPrefix(f.Check, "generic-") || (f.Check == "dangling-object-ref" && definesGeneric)
		if !isGenerics {
			continue
		}
		fmt.Fprintf(w, "warning: generics: %s: %s [%s]\n", f.Location, f.Message, f.Check)
		count++
	}
	if count > 0 {
		fmt.Fprintln(w, `note: generics problems resolve to "any" in the output; run "wfn validate" for the full report`)
	}
}