# testdata/convert/

Fixtures for `wfn convert`. Every fixture here, plus the non-generic ones in
`testdata/validate/` (the `generics-*` ones are skipped on purpose), is run
through `openapiconvert.Generate` / `postmanconvert.Generate` by
`golden_test.go` in each package, with and without `includePrivate`, and
compared byte-for-byte with `golden/{openapi,postman}/<dir>-<fixture>.<public|private>.json`.
After an intended output change, regenerate with `go test ./openapiconvert
./postmanconvert -run TestGoldenOutput -update` and review the diff.

- **private-flag.json** - `private` at every level it's allowed: one
  private endpoint (with an argument), a private argument on a public
  endpoint, a private argument inside `object.filter`, a private
  endpoint attribute, and a private attribute inside `object.person`.
  By default none of the private names (`internal-sync`, `sync_secret`,
  `debug_trace`, `hidden_override`, `internal_flag`, `audit_ref`) may
  appear in either output; with `includePrivate` true, all appear in the
  OpenAPI output, and the argument ones in Postman (which only renders
  request bodies, not response attributes).

- **types.json** - every type form: each string and number refinement, plain
  `number`/`boolean`/`any`, bare `array`/`object`, array-of, nullable and
  non-null unions, `choices` and `values`, a paginated endpoint, an object
  used in both argument and attribute context (`address`), and a
  self-referential object (`node`). Exists so the golden files reach the
  type-mapping code in both converters.

`private-flag.json` is also used to check every codegen target (`js`, `php`,
`python`, `java`, `csharp`, `go`, `ruby`): by default no private name may
appear in the generated code; with `includePrivate` all of them do (except
where a target never renders a given member, e.g. ruby's object-argument
hashes).