package gogen

import "strings"

// docLines splits a docs string (markdown, per the wire spec) into lines
// safe to place inside a Go "//" comment block. Unlike jsgen/phpgen's
// docLines, no escaping is needed here - "//" line comments have no
// closing token to accidentally trigger early, unlike JSDoc/PHPDoc's
// "*/".
func docLines(docs string) []string {
	docs = strings.TrimSpace(docs)
	if docs == "" {
		return nil
	}
	rawLines := strings.Split(docs, "\n")
	lines := make([]string, 0, len(rawLines))
	for _, l := range rawLines {
		lines = append(lines, strings.TrimRight(l, " \t\r"))
	}
	return lines
}

// writeComment writes lines as a Go "//" comment block, one line per
// call to b.WriteString, prefixed with indent.
func writeComment(b *strings.Builder, indent string, lines []string) {
	for _, l := range lines {
		if l == "" {
			b.WriteString(indent + "//\n")
			continue
		}
		b.WriteString(indent + "// " + l + "\n")
	}
}
