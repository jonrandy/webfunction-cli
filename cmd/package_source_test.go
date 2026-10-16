package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/webfunction-protocol/webfunction-go"
)

// fixtureFiles lists every package fixture, validate and convert alike.
func fixtureFiles(t *testing.T) []string {
	t.Helper()
	var files []string
	for _, pattern := range []string{"validate/*.json", "convert/*.json"} {
		matches, err := filepath.Glob(filepath.Join("..", "testdata", pattern))
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, matches...)
	}
	if len(files) == 0 {
		t.Fatal("no fixtures found")
	}
	return files
}

// serveFile answers every request with the file's contents, standing in
// for a Web Function package endpoint.
func serveFile(t *testing.T, path string) string {
	t.Helper()
	body := readFile(t, path)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv.URL
}

// runCommand runs a command, returning what it printed to stdout and its error.
func runCommand(t *testing.T, run func([]string) error, args ...string) (string, error) {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()
	runErr := run(args)
	w.Close()
	os.Stdout = old
	return <-done, runErr
}

func convertRun(args []string) error  { return (&ConvertCommand{}).Run(args) }
func validateRun(args []string) error { return (&ValidateCommand{}).Run(args) }
func codegenRun(args []string) error  { return (&CodegenCommand{}).Run(args) }

func tempOut(t *testing.T, name string) string { return filepath.Join(t.TempDir(), name) }

func absPath(t *testing.T, path string) string {
	t.Helper()
	abs, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.ToSlash(abs)
}

func TestFileURL(t *testing.T) {
	cases := []struct{ in, want string }{
		{"/home/me/pkg.json", "file:///home/me/pkg.json"},
		{"/tmp/my docs/pkg.json", "file:///tmp/my%20docs/pkg.json"},
		{"C:/Users/me/pkg.json", "file:///C:/Users/me/pkg.json"},
	}
	for _, c := range cases {
		if got := fileURL(c.in); got != c.want {
			t.Errorf("fileURL(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestLoadPackageFile(t *testing.T) {
	rel := fixturePath("generics-clean.json")
	pkg, from, err := loadPackageFile(rel)
	if err != nil {
		t.Fatal(err)
	}
	if pkg == nil || len(pkg.Endpoints) != 4 {
		t.Fatalf("expected the 4-endpoint fixture, got %+v", pkg)
	}
	if want := "file://" + filepath.ToSlash(absPath(t, rel)); !strings.HasSuffix(from, "/testdata/validate/generics-clean.json") || !strings.HasPrefix(from, "file:///") || strings.Contains(from, "..") {
		t.Errorf("location should be an absolute file:// URL like %q, got %q", want, from)
	}

	bad := tempOut(t, "bad.json")
	if err := os.WriteFile(bad, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := loadPackageFile(bad); err == nil || !strings.Contains(err.Error(), "parsing") || !strings.Contains(err.Error(), "bad.json") {
		t.Errorf("invalid JSON should give a parsing error naming the file, got %v", err)
	}

	missing := filepath.Join(t.TempDir(), "nope.json")
	if _, _, err := loadPackageFile(missing); err == nil || !strings.Contains(err.Error(), "nope.json") {
		t.Errorf("a missing file should give an error naming it, got %v", err)
	}
}

// Reading a package from a file must give exactly the package fetching it
// from an endpoint does, for every fixture, so the two routes can never
// disagree about what a package means.
func TestFilePackageEqualsEndpointPackage(t *testing.T) {
	for _, path := range fixtureFiles(t) {
		t.Run(filepath.Base(filepath.Dir(path))+"-"+filepath.Base(path), func(t *testing.T) {
			fromFile, _, err := loadPackageFile(path)
			if err != nil {
				t.Fatal(err)
			}
			client, err := webfunction.FromPackageEndpoint(serveFile(t, path), webfunction.Options{})
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(fromFile, client.Package()) {
				a, _ := json.Marshal(fromFile)
				b, _ := json.Marshal(client.Package())
				t.Errorf("packages differ\n file: %.300s\n http: %.300s", a, b)
			}
		})
	}
}

func TestConvertFromFileEqualsConvertFromURL(t *testing.T) {
	for _, fixture := range []string{"convert/types.json", "convert/generics-types.json", "convert/private-flag.json"} {
		path := filepath.Join("..", "testdata", fixture)
		for _, target := range []string{"openapi", "postman"} {
			for _, extra := range [][]string{nil, {"--private"}} {
				name := fixture + "-" + target + strings.Join(extra, "")
				t.Run(name, func(t *testing.T) {
					viaURL, viaFile := tempOut(t, "url.json"), tempOut(t, "file.json")
					urlStdout, err := runCommand(t, convertRun, append([]string{"--target", target, "--url", serveFile(t, path), "-o", viaURL}, extra...)...)
					if err != nil {
						t.Fatal(err)
					}
					fileStdout, err := runCommand(t, convertRun, append([]string{"--target", target, "--file", path, "-o", viaFile}, extra...)...)
					if err != nil {
						t.Fatal(err)
					}
					if readFile(t, viaURL) != readFile(t, viaFile) {
						t.Error("converting from a file must produce the same document as converting from a URL")
					}
					if !strings.Contains(urlStdout, "Fetched ") || strings.Contains(urlStdout, "Loaded ") {
						t.Errorf("URL route should still say Fetched, got %q", urlStdout)
					}
					if !strings.Contains(fileStdout, "Loaded ") || !strings.Contains(fileStdout, "from file:///") || strings.Contains(fileStdout, "Fetched ") {
						t.Errorf("file route should say Loaded ... from file:///..., got %q", fileStdout)
					}
				})
			}
		}
	}
}

func TestConvertFromFileStillWarnsAboutGenerics(t *testing.T) {
	var warned bytes.Buffer
	old := warnWriter
	warnWriter = &warned
	t.Cleanup(func() { warnWriter = old })

	_, err := runCommand(t, convertRun, "--target", "openapi", "--file", fixturePath("generics-flawed.json"), "-o", tempOut(t, "out.json"))
	if err != nil {
		t.Fatalf("convert failed: %v", err)
	}
	if !strings.Contains(warned.String(), "[generic-unapplied]") {
		t.Errorf("expected generics warnings, got %q", warned.String())
	}
}

func TestPackageFlagErrors(t *testing.T) {
	file := fixturePath("generics-clean.json")
	out := func() string { return tempOut(t, "out") }
	cases := []struct {
		name string
		run  func([]string) error
		args []string
		want []string // every one must appear in the error
	}{
		{"convert: both flags", convertRun, []string{"--target", "openapi", "--url", "http://x.test/p", "--file", file, "-o", out()}, []string{"--url", "--file", "cannot be used together"}},
		{"convert: neither flag", convertRun, []string{"--target", "openapi", "-o", out()}, []string{"missing required flag", "--url or --file"}},
		{"convert: missing file", convertRun, []string{"--target", "openapi", "--file", filepath.Join(t.TempDir(), "nope.json"), "-o", out()}, []string{"nope.json"}},
		{"validate: both flags", validateRun, []string{"--url", "http://x.test/p", "--file", file, "-o", out()}, []string{"--url", "--file", "cannot be used together"}},
		{"validate: neither flag", validateRun, []string{"-o", out()}, []string{"missing required flag", "--url or --file"}},
		{"validate: missing file", validateRun, []string{"--file", filepath.Join(t.TempDir(), "nope.json"), "-o", out()}, []string{"nope.json"}},
		{"codegen: file without url", codegenRun, []string{"--target", "js", "--file", file, "-o", out()}, []string{"missing required flag", "--url"}},
		{"codegen: missing file", codegenRun, []string{"--target", "js", "--url", "http://127.0.0.1:1/p", "--file", filepath.Join(t.TempDir(), "nope.json"), "-o", out()}, []string{"nope.json"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := runCommand(t, c.run, c.args...)
			if err == nil {
				t.Fatal("expected an error")
			}
			for _, w := range c.want {
				if !strings.Contains(err.Error(), w) {
					t.Errorf("error should mention %q, got: %v", w, err)
				}
			}
		})
	}
}

// A URL that cannot be fetched must still fail the way it always has, so
// adding --file did not change the existing error.
func TestFetchFailureIsReportedAsFetching(t *testing.T) {
	const dead = "http://127.0.0.1:1/package" // nothing listens here
	cases := []struct {
		name string
		run  func([]string) error
		args []string
	}{
		{"convert", convertRun, []string{"--target", "openapi", "--url", dead, "-o", tempOut(t, "out")}},
		{"validate", validateRun, []string{"--url", dead, "-o", tempOut(t, "out")}},
		{"codegen", codegenRun, []string{"--target", "js", "--url", dead, "-o", tempOut(t, "out")}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := runCommand(t, c.run, c.args...)
			if err == nil || !strings.Contains(err.Error(), "fetching package") {
				t.Errorf("want a fetching package error, got %v", err)
			}
		})
	}
}

func readReport(t *testing.T, path string) map[string]any {
	t.Helper()
	var report map[string]any
	if err := json.Unmarshal([]byte(readFile(t, path)), &report); err != nil {
		t.Fatalf("report is not JSON: %v", err)
	}
	return report
}

func TestValidateFromFile(t *testing.T) {
	for _, name := range []string{"generics-clean.json", "generics-flawed.json"} {
		t.Run(name, func(t *testing.T) {
			path := fixturePath(name)
			viaURL, viaFile := tempOut(t, "url.json"), tempOut(t, "file.json")
			_, urlErr := runCommand(t, validateRun, "--url", serveFile(t, path), "-o", viaURL)
			_, fileErr := runCommand(t, validateRun, "--file", path, "-o", viaFile)
			if (urlErr == nil) != (fileErr == nil) {
				t.Errorf("both routes should succeed or fail together: url=%v file=%v", urlErr, fileErr)
			}

			fileReport, urlReport := readReport(t, viaFile), readReport(t, viaURL)
			got, _ := fileReport["package_url"].(string)
			if !strings.HasPrefix(got, "file:///") || !strings.HasSuffix(got, "/testdata/validate/"+name) {
				t.Errorf("package_url should be the file:// URL of the file, got %q", got)
			}
			// The URL differs by design and the timestamp by the clock; everything else must match.
			for _, key := range []string{"package_url", "generated_at"} {
				delete(fileReport, key)
				delete(urlReport, key)
			}
			if !reflect.DeepEqual(fileReport, urlReport) {
				t.Error("validating a file must give the same findings as validating the same package from a URL")
			}
		})
	}
}

func TestCodegenFromFile(t *testing.T) {
	path := filepath.Join("..", "testdata", "convert", "types.json")
	const runtimeURL = "http://127.0.0.1:1/package" // nothing listens here: any fetch would fail
	for _, target := range validTargets {
		t.Run(target, func(t *testing.T) {
			outName := "client.out"
			if target == "ruby" {
				outName = "client.rb"
			}
			viaFile, viaURL := tempOut(t, outName), tempOut(t, outName)

			fileStdout, err := runCommand(t, codegenRun, "--target", target, "--url", runtimeURL, "--file", path, "-o", viaFile)
			if err != nil {
				t.Fatalf("codegen from a file failed: %v", err)
			}
			srvURL := serveFile(t, path)
			urlStdout, err := runCommand(t, codegenRun, "--target", target, "--url", srvURL, "-o", viaURL)
			if err != nil {
				t.Fatalf("codegen from a URL failed: %v", err)
			}
			if !strings.Contains(fileStdout, "Loaded ") || !strings.Contains(fileStdout, "from file:///") || strings.Contains(fileStdout, "Fetched ") {
				t.Errorf("file route should say Loaded ... from file:///..., got %q", fileStdout)
			}
			if !strings.Contains(urlStdout, "Fetched ") || strings.Contains(urlStdout, "Loaded ") {
				t.Errorf("URL route should still say Fetched, got %q", urlStdout)
			}

			got := readFile(t, viaFile)
			if !strings.Contains(got, runtimeURL) {
				t.Error("the generated client must embed the --url, which is its runtime address")
			}
			if strings.Contains(got, "file://") {
				t.Error("the file's location must not leak into the generated client")
			}
			if want := strings.ReplaceAll(readFile(t, viaURL), srvURL, runtimeURL); got != want {
				t.Error("generating from a file must match generating from the same package fetched from a URL")
			}
		})
	}
}