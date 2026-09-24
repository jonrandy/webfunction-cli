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
  package: 5 errors / 5 warnings / 3 info notes, `Failed()=true`. Covers:
  - `event-source-invalid-return-type` (`list-people` sets `event_source`
    but returns `object.ghost`, not `["string"]`)
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
  - `unrecognized-flag` at package level (`made_up_package_flag`),
    endpoint level (`made_up_flag`), and attribute level
    (`made_up_attribute_flag` on `person`'s `id` attribute)
- **event-source-valid.json** - a correctly-typed `event_source`
  endpoint (`returns: ["string"]`) - must produce **zero** findings,
  confirming `event-source-invalid-return-type` has no false positive
  on the valid case.
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

## Spec pages resolved several open questions (Jon asked about /error, /package)

Reading https://webfunction.org/package and https://webfunction.org/error
directly resolved multiple things this project had left as open questions:

- **Package-level vs endpoint-level `errors` isn't inheritance at all.**
  `/error` states they're two separate, parallel, purely advisory lists -
  package-level for codes shared across endpoints (auth, rate limits),
  endpoint-level for that endpoint's own business logic. A server may
  return codes in neither list; clients must handle unlisted codes
  regardless. There's nothing to be ambiguous about, so the
  `ambiguous-package-error-inheritance` info note (and
  `checkSpecAmbiguities` entirely) was removed rather than kept "resolved
  as fine." This also confirms codegen's own convention (scoping
  `@throws`-style output to an endpoint's own `errors` only) was already
  spec-correct.
- **`events`/`event_source_url`/`pipeline_url` are real, spec-defined
  package keys** (`/package` documents both plus a full `Event`
  schema) - not an early over-read as `webfunction-go`'s `Package`
  doc comment had guessed. Comment corrected in `webfunction-go/
  package.go` to reflect this; the Ruby reference client apparently just
  hasn't implemented them yet. Actually adding `EventSourceURL`/`Events`
  fields to `Package` (a real feature, needs an `Event` type) is left as
  a flagged follow-up, out of scope for this session.
- **New check unlocked by this**: `/package`'s flag table states an
  endpoint with `event_source` MUST declare `returns` of exactly
  `["string"]` - a real, checkable MUST. Added as
  `event-source-invalid-return-type` (error severity,
  `isBareStringReturn` helper checks for a single unrefined `string`
  alternative and nothing else).
- **Attribute-level flags were a gap**: `/package`'s flag table lists
  `nullable` as an Attribute-scoped flag, but `checkFlags` never checked
  attribute flags at all (endpoint or object attributes). Added
  `knownAttributeFlags = {"nullable"}` and wired attribute-flag checking
  into `checkFlags` for both endpoint attributes and object attributes.
  While tightening this, also removed `"versioned"` from
  `knownEndpointFlags` - the spec states each flag MUST only be used at
  its designated level, and `versioned` is Package-only.
- **The base_url trailing-slash check is gone, not just re-severitized.**
  `/package`'s "URL composition" section states the actual join rule
  plainly: append directly if `base_url` ends in `/`, otherwise insert a
  single `/`. That's simple string normalization, not RFC 3986 relative
  reference resolution - which is what `webfunction-go`'s `client.go`
  was actually using (mirroring what was believed to be Ruby's
  `URI.join` behavior). That mismatch was a real, separate bug in
  `webfunction-go` itself (not a package-authoring pitfall `validate`
  should be flagging): a `base_url` with a path and no trailing slash
  had its last path segment silently replaced instead of the endpoint
  name being appended under it. Fixed `joinEndpointURL` in
  `webfunction-go/client.go` to do the spec's plain normalization
  instead, added `TestJoinEndpointURL` there as a permanent regression
  test, and removed the `base-url-missing-trailing-slash` check from
  `validate` entirely - with the client fixed, `base_url` no longer
  needs a trailing slash under any circumstance, so there's nothing left
  to check. (A different, real base_url check the spec's "Validation
  requirements" section does call for - `base_url` MUST be valid per RFC
  3986 - was flagged as a good next candidate but not built this
  session, to stay scoped to what was actually asked.)

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