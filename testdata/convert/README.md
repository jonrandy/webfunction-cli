# testdata/convert/

Fixtures for `wfn convert`. Like `testdata/validate/`, not wired into
`go test`; exercised by a throwaway harness that reads the fixture JSON
directly and calls `openapiconvert.Generate` / `postmanconvert.Generate`
(second argument is `includePrivate`).

- **private-flag.json** - `private` at every level it's allowed: one
  private endpoint (with an argument), a private argument on a public
  endpoint, a private argument inside `object.filter`, a private
  endpoint attribute, and a private attribute inside `object.person`.
  By default none of the private names (`internal-sync`, `sync_secret`,
  `debug_trace`, `hidden_override`, `internal_flag`, `audit_ref`) may
  appear in either output; with `includePrivate` true, all appear in the
  OpenAPI output, and the argument ones in Postman (which only renders
  request bodies, not response attributes).

`private-flag.json` is also used to check `wfn codegen --target js` and
`--target php`, `--target python`, `--target java`, `--target csharp` and
`--target go` (and, as each target gains `--private` support, the others): by default
no private name may appear in the generated code; with `includePrivate`
all of them do.