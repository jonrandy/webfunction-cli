# testdata/validate/

Fixtures for `wfn validate`. Like `testdata/openapi/` and
`testdata/postman/`, these aren't wired into `go test` (no checked-in
mock HTTP server) - they were exercised manually via a throwaway
`cmd/genvalidate/main.go` harness that reads a fixture JSON directly
(skips the network fetch entirely, since a local mock server hangs in
this sandbox) and calls `validate.Run` + `Report.Render` for each
format. That harness is not committed; regenerate it if you need to
re-verify by hand:

```go
// cmd/genvalidate/main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/webfunction-protocol/webfunction-go"
	"wfn/validate"
)

func main() {
	data, _ := os.ReadFile(os.Args[1])
	var pkg webfunction.Package
	if err := json.Unmarshal(data, &pkg); err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		os.Exit(1)
	}
	report := validate.Run(&pkg, "file://"+os.Args[1])
	format := "json"
	if len(os.Args) > 2 {
		format = os.Args[2]
	}
	out, _ := report.Render(format)
	fmt.Println(out)
	fmt.Fprintf(os.Stderr, "%d error(s), %d warning(s), %d info note(s); Failed()=%v\n",
		report.ErrorCount, report.WarningCount, report.InfoCount, report.Failed())
}
```

## Fixtures

- **clean.json** - a well-formed package. Verified to produce **zero
  findings** across all three formats. (Building this fixture correctly
  actually required fixing a first-draft mistake - see "found via
  fixtures" below.)
- **flawed.json** - deliberately exercises nearly every check in one
  package: 5 errors / 5 warnings / 4 info notes, `Failed()=true`. Covers:
  - `base-url-missing-trailing-slash` (`base_url` is
    `https://api.example.com/v1`, no trailing slash - see the note
    below on why this is error-severity, not a style nit)
  - `duplicate-endpoint-name` (two `list-people` endpoints)
  - `duplicate-object-name` (two `person` objects, one argument-context
    only, one attribute-context only - the second is masked, not just
    redundant)
  - `versioned-empty-versions` (`"versioned"` flag, no `versions` list)
  - `dangling-object-ref` (`list-people` returns `object.ghost`, which
    doesn't exist)
  - `bare-object-return` / `bare-array-return` (two endpoints, one of
    each, neither with endpoint-level `attributes`)
  - `duplicate-choice-value` and `choice-value-type-mismatch` (the
    `status` argument's choices include a repeated `"active"` and a
    numeric `1` against a `string` type)
  - `unrecognized-flag` at both package level (`made_up_package_flag`)
    and endpoint level (`made_up_flag`)
  - `ambiguous-package-error-inheritance` (package declares `errors`)
  - `ambiguous-events-concept` (an endpoint sets `event_source`)
- **versioned-mismatch.json** - isolated case for `version-not-in-
  versions`: `"versioned"` flag set, a real non-empty `versions` list,
  but `version` isn't a member of it. Kept separate from flawed.json
  because it's mutually exclusive with `versioned-empty-versions` in the
  same package (checkVersioning returns after the empty-list case).
- **array-element-choices.json** - regression fixture for a real bug Jon
  caught after the first delivery: for an `array<string>`-typed argument,
  `choices` entries are individual allowed *elements*, not whole allowed
  arrays, so they must be checked against the array's element type, not
  the array type itself (`Type.Valid` on an array alt expects the value
  being checked to itself be a `[]any`). Has a `tags` argument
  (`array<string>`, choices `["red","green","blue"]`) that must produce
  **zero** findings, and a `bad-tags` argument (same type, choices
  `["red", 42, "red"]`) that must produce exactly a
  `choice-value-type-mismatch` on `42` and a `duplicate-choice-value` on
  the repeated `"red"` - confirming the element type (not the array
  type) is what's actually being checked, and that a valid string
  choice against an array<string> field is no longer a false positive.

## base_url check reversed (Jon's correction)

The original check flagged a *trailing* slash on `base_url` as a
warning, on the assumption of naive `base_url + "/" + name`
concatenation. Jon pointed out `base_url` should actually *have* a
trailing slash, which prompted checking how the real reference clients
join it - `webfunction-go`'s `client.go` (mirroring Ruby's
`URI.join`) uses RFC 3986 relative reference resolution, not string
concatenation. Verified directly: a `base_url` with a path segment and
no trailing slash has the endpoint name silently *replace* that last
path segment (`https://api.example.com/v1` + `list-people` ->
`https://api.example.com/list-people`, `v1` dropped) rather than
appending it. A bare-domain `base_url` (no path at all) is unaffected
either way. Reclassified as `base-url-missing-trailing-slash`, moved
from checks_warning.go to checks_error.go, and bumped to **error**
severity (confirmed with Jon) since it's a silent wrong-request bug,
not a style nit. The finding message names the exact path segment that
would be dropped.

## Found via fixtures (not by reading code)

- **clean.json's own first draft had a bug**: `"returns": ["array",
  ["object.person"]]` was meant to mean "array of `object.person`", but
  per the wire grammar that's actually a *union* of a bare untyped array
  AND a nested array-of-object.person - two separate alternatives, not
  one. Running it through validate correctly flagged the bare-array
  alternative as `bare-array-return`, which is what caught the mistake.
  Fixed to the correct single-nested-array form: `"returns":
  [["object.person"]]`. Left in this note deliberately - it's a real
  example of the "verify by running" convention doing its job on the
  fixture itself, not just on the code under test.
