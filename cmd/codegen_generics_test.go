package cmd

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/webfunction-protocol/webfunction-go"

	"wfn/csharpgen"
	"wfn/generics"
	"wfn/gogen"
	"wfn/javagen"
	"wfn/jsgen"
	"wfn/phpgen"
	"wfn/pythongen"
	"wfn/rubygen"
)

// monomorphTargets lists the codegen targets that expand generics so far.
// Each step of the rollout adds one; every test below runs for all of them.
var monomorphTargets = []string{"js", "go", "php", "python", "java", "csharp", "ruby"}

func genericsFixture(name string) string {
	return filepath.Join("..", "testdata", "codegen", "generics", name)
}

// genericPairs are generic packages and the hand-monomorphized twin that
// says what their generated code must be.
var genericPairs = []struct{ name, generic, twin string }{
	{"generic", "generic.json", "twin.json"},
	{"collision", "collision.json", "collision-twin.json"},
}

// describe renders a package's endpoints and objects as text, so two
// packages can be compared member by member.
func describe(pkg *webfunction.Package) string {
	var b strings.Builder
	members := func(indent string, args []webfunction.Argument, attrs []webfunction.Attribute) {
		for _, a := range args {
			fmt.Fprintf(&b, "%sarg %s: %s %v\n", indent, a.Name, a.Type.String(), a.Flags)
		}
		for _, a := range attrs {
			fmt.Fprintf(&b, "%sattr %s: %s %v\n", indent, a.Name, a.Type.String(), a.Flags)
		}
	}
	for _, e := range pkg.Endpoints {
		fmt.Fprintf(&b, "endpoint %s returns %s %v\n", e.Name, e.Returns.String(), e.Flags)
		members("  ", e.Arguments, e.Attributes)
	}
	for _, o := range pkg.Objects {
		fmt.Fprintf(&b, "object %s\n", o.Name)
		members("  ", o.Arguments, o.Attributes)
	}
	return b.String()
}

// generateDirect calls a target's Generate itself, returning everything it
// produced as one string.
func generateDirect(t *testing.T, target string, pkg *webfunction.Package, private bool) string {
	t.Helper()
	var out, rbs string
	var err error
	switch target {
	case "js":
		out, err = jsgen.Generate(pkg, goldenSourceURL, private)
	case "go":
		out, err = gogen.Generate(pkg, goldenSourceURL, defaultNamespace, private)
	case "php":
		out, err = phpgen.Generate(pkg, goldenSourceURL, defaultNamespace, private)
	case "python":
		out, err = pythongen.Generate(pkg, goldenSourceURL, private)
	case "java":
		out, err = javagen.Generate(pkg, goldenSourceURL, defaultNamespace, private)
	case "csharp":
		out, err = csharpgen.Generate(pkg, goldenSourceURL, defaultNamespace, private)
	case "ruby":
		out, rbs, err = rubygen.Generate(pkg, goldenSourceURL, defaultNamespace, private)
	default:
		t.Fatalf("unknown target %q", target)
	}
	if err != nil {
		t.Fatalf("%s: %v", target, err)
	}
	return out + "\n--- rbs ---\n" + rbs
}

// allFiles joins every file a target wrote into one string.
func allFiles(files map[string]string) string {
	var b strings.Builder
	for _, ext := range []string{".js", ".go", ".php", ".py", ".java", ".cs", ".rb", ".rbs"} {
		b.WriteString(files[ext])
	}
	return b.String()
}

var nonAlnum = regexp.MustCompile(`[^a-z0-9]`)

// squash lowercases text and drops everything but letters and digits, so a
// member can be found whatever case convention a target gives it.
func squash(s string) string { return nonAlnum.ReplaceAllString(strings.ToLower(s), "") }

// firstDifference describes where two generated files first differ.
func firstDifference(got, want string) string {
	g, w := strings.Split(got, "\n"), strings.Split(want, "\n")
	for i := 0; i < len(g) && i < len(w); i++ {
		if g[i] != w[i] {
			return fmt.Sprintf("line %d:\n   got: %s\n  want: %s", i+1, g[i], w[i])
		}
	}
	return fmt.Sprintf("lengths differ: %d lines, want %d", len(g), len(w))
}

// The twins must say what Expand says, or the tests below prove nothing.
func TestGenericTwinsMatchExpansion(t *testing.T) {
	for _, p := range genericPairs {
		t.Run(p.name, func(t *testing.T) {
			expanded := describe(generics.Expand(parsePackageJSON(t, readFile(t, genericsFixture(p.generic)))))
			twin := describe(parsePackageJSON(t, readFile(t, genericsFixture(p.twin))))
			if expanded != twin {
				t.Errorf("the twin differs from the expansion.\nexpanded:\n%s\ntwin:\n%s", expanded, twin)
			}
		})
	}
}

func TestGenericCodegenFixturesAreClean(t *testing.T) {
	for _, p := range genericPairs {
		for _, name := range []string{p.generic, p.twin} {
			if lines := warnings(t, parsePackageJSON(t, readFile(t, genericsFixture(name)))); lines != nil {
				t.Errorf("%s should have no generics problems, got:\n%s", name, strings.Join(lines, "\n"))
			}
		}
	}
}

// A generic package must generate exactly the code its hand-monomorphized
// twin generates, public and --private.
func TestCodegenGenericsMatchTheTwin(t *testing.T) {
	for _, target := range monomorphTargets {
		for _, p := range genericPairs {
			for _, private := range []bool{false, true} {
				name := fmt.Sprintf("%s/%s/private=%v", target, p.name, private)
				t.Run(name, func(t *testing.T) {
					got := generateCode(t, target, genericsFixture(p.generic), private)
					want := generateCode(t, target, genericsFixture(p.twin), private)
					for ext, text := range want {
						if got[ext] != text {
							t.Errorf("%s differs from the twin's, %s", ext, firstDifference(got[ext], text))
						}
					}
				})
			}
		}
	}
}

func TestCodegenNamesTheGeneratedObjects(t *testing.T) {
	templateName := regexp.MustCompile(`\b(List|Page|Filter|Pair|Both|Wrapper)(Attributes|Args)\b`)
	for _, target := range monomorphTargets {
		t.Run(target, func(t *testing.T) {
			out := allFiles(generateCode(t, target, genericsFixture("generic.json"), false))
			names := []string{"ListOfUser", "ListOfUserOrNull", "PageOfUser", "PageOfPageOfUser", "PageOfArrayOfUser", "PairOfNumber", "PairOfString", "BothOfNumber", "WrapperOfString"}
			if target != "ruby" {
				// Ruby forwards arguments as plain keyword arguments and gives an
				// arguments-only object no class, in the twin as well.
				names = append(names, "FilterOfString")
			}
			for _, name := range names {
				if !strings.Contains(out, name) {
					t.Errorf("the generated code should contain %s", name)
				}
			}
			if m := templateName.FindString(out); m != "" {
				t.Errorf("a generic template must not be emitted as a type, found %s", m)
			}
		})
	}
}

// Private flags travel onto instantiations, and an instantiation only a
// private endpoint reaches is emitted only with --private.
func TestCodegenGenericInstantiationsHonorPrivate(t *testing.T) {
	for _, target := range monomorphTargets {
		t.Run(target, func(t *testing.T) {
			public := squash(allFiles(generateCode(t, target, genericsFixture("generic.json"), false)))
			private := squash(allFiles(generateCode(t, target, genericsFixture("generic.json"), true)))
			cases := []struct{ what, needle string }{
				{"a private attribute of a generic", "cursor"},
				{"a private attribute of a plain object", "recoverycode"},
				{"a private argument", "debugtrace"},
				{"a private endpoint", "listaudit"},
				{"an instantiation only a private endpoint reaches", "pageofaudit"},
			}
			for _, c := range cases {
				if strings.Contains(public, c.needle) {
					t.Errorf("%s must not appear without --private", c.what)
				}
				if !strings.Contains(private, c.needle) {
					t.Errorf("%s must appear with --private", c.what)
				}
			}
		})
	}
}

// Generate must not mutate the package it is given, and must give the same
// answer every time.
func TestCodegenGenericsIsPureAndDeterministic(t *testing.T) {
	for _, target := range monomorphTargets {
		t.Run(target, func(t *testing.T) {
			pkg := parsePackageJSON(t, readFile(t, genericsFixture("generic.json")))
			before := describe(pkg)
			first := generateDirect(t, target, pkg, false)
			second := generateDirect(t, target, pkg, false)
			if first != second {
				t.Errorf("two runs differ, %s", firstDifference(second, first))
			}
			if describe(pkg) != before {
				t.Error("Generate mutated the package it was given")
			}
		})
	}
}

// A user object and a generated one that become the same identifier are
// both kept (the target numbers the second), never merged or dropped.
func TestCodegenKeepsBothObjectsWhenNamesCollide(t *testing.T) {
	for _, target := range monomorphTargets {
		t.Run(target, func(t *testing.T) {
			out := squash(allFiles(generateCode(t, target, genericsFixture("collision.json"), false)))
			for _, needle := range []string{"oldfield", "items"} {
				if !strings.Contains(out, needle) {
					t.Errorf("%s is missing from the output; one of the two objects was lost", needle)
				}
			}
		})
	}
}