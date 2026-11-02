package cmd

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The generated Ruby client must be usable without repeating the --url it
// was generated from, as the Go and Python clients already are: Client
// carries the URL as SOURCE_URL and Client.new_client builds a ready-to-use
// client from it. from_package_endpoint(url, ...) stays for callers who
// want a different address.

func generateRubyFor(t *testing.T, url string) (rb, rbs string) {
	t.Helper()
	oldWarn := warnWriter
	warnWriter = io.Discard
	defer func() { warnWriter = oldWarn }()
	out := filepath.Join(t.TempDir(), "client.rb")
	args := []string{"--target", "ruby", "--url", url, "--file", fixturePath("clean.json"), "-o", out}
	if _, err := runCommand(t, codegenRun, args...); err != nil {
		t.Fatalf("codegen --target ruby failed: %v", err)
	}
	return readFile(t, out), readFile(t, strings.TrimSuffix(out, ".rb")+".rbs")
}

func TestRubyClientBakesSourceURL(t *testing.T) {
	rb, rbs := generateRubyFor(t, goldenSourceURL)
	for _, want := range []string{
		`SOURCE_URL = "` + goldenSourceURL + `"`,
		"def self.new_client(bearer_auth: nil, version: nil, pipelined: false)",
		"from_package_endpoint(SOURCE_URL, bearer_auth: bearer_auth, version: version, pipelined: pipelined)",
		"def self.from_package_endpoint(url, bearer_auth: nil, version: nil, pipelined: false)",
	} {
		if !strings.Contains(rb, want) {
			t.Errorf("generated .rb is missing %q", want)
		}
	}
	for _, want := range []string{
		"SOURCE_URL: String",
		"def self.new_client: (?bearer_auth: String?, ?version: String?, ?pipelined: bool) -> Client",
		"def self.from_package_endpoint: (String url, ?bearer_auth: String?, ?version: String?, ?pipelined: bool) -> Client",
	} {
		if !strings.Contains(rbs, want) {
			t.Errorf("generated .rbs is missing %q", want)
		}
	}
}

// A URL is data, not code: Ruby interpolates "#{...}" inside double quotes,
// so the baked-in literal must escape it.
func TestRubySourceURLIsNotInterpolated(t *testing.T) {
	rb, _ := generateRubyFor(t, `https://example.test/p?q=#{1+1}&r="x"`)
	want := `SOURCE_URL = "https://example.test/p?q=\#{1+1}&r=\"x\""`
	if !strings.Contains(rb, want) {
		t.Errorf("generated .rb is missing the escaped literal %q", want)
	}
}

// Loads the generated file with real Ruby and checks the constant and the
// constructor exist, and that the URL comes back byte-for-byte. Runs only
// when WFN_RUBY_LIB points at the webfunction-ruby gem's lib directory.
func TestRubyGeneratedClientLoadsWithSourceURL(t *testing.T) {
	lib := os.Getenv("WFN_RUBY_LIB")
	if lib == "" {
		t.Skip("set WFN_RUBY_LIB to the webfunction-ruby lib directory to run this test")
	}
	if _, err := exec.LookPath("ruby"); err != nil {
		t.Skip("ruby not installed")
	}
	abs, err := filepath.Abs(lib)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(abs, "web_function.rb")); err != nil {
		t.Fatalf("WFN_RUBY_LIB=%q: no web_function.rb in %s (use an absolute path; a leading ~ is not expanded here): %v", lib, abs, err)
	}
	url := `https://example.test/p?q=#{1+1}&r="x"`
	rb, _ := generateRubyFor(t, url)
	path := filepath.Join(t.TempDir(), "client.rb")
	if err := os.WriteFile(path, []byte(rb), 0o644); err != nil {
		t.Fatal(err)
	}
	script := `require ARGV[0]; c = WebFunctionClient::Client; puts c::SOURCE_URL; puts c.respond_to?(:new_client)`
	got, err := exec.Command("ruby", "-I", abs, "-e", script, path).CombinedOutput()
	if err != nil {
		t.Fatalf("ruby failed: %v\n%s", err, got)
	}
	if want := url + "\ntrue\n"; string(got) != want {
		t.Errorf("ruby printed %q, want %q", got, want)
	}
}