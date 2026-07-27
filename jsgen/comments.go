package jsgen

import "strings"

// docLines splits a docs string (markdown, per the wire spec) into lines
// safe to place inside a JSDoc block comment - each returned line has no
// leading/trailing whitespace and can't accidentally close the comment
// early.
func docLines(docs string) []string {
	docs = strings.TrimSpace(docs)
	if docs == "" {
		return nil
	}

	// "*/" would close the enclosing JSDoc comment early; escape it.
	docs = strings.ReplaceAll(docs, "*/", "*\\/")

	rawLines := strings.Split(docs, "\n")
	lines := make([]string, 0, len(rawLines))
	for _, l := range rawLines {
		lines = append(lines, strings.TrimRight(l, " \t\r"))
	}
	return lines
}