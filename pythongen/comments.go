package pythongen

import "strings"

// docLines splits a docs string (markdown, per the wire spec) into lines
// safe to place inside a triple-quoted docstring. Python docstrings need
// no per-line escaping the way C#'s XML doc comments do (no "<"/"&"
// concerns) - the only real hazard is a literal `"""` sequence inside the
// docs text itself, which would prematurely close the docstring; that's
// neutralized by escaping any embedded `"""` to `\"\"\"`.
func docLines(docs string) []string {
	docs = strings.TrimSpace(docs)
	if docs == "" {
		return nil
	}
	rawLines := strings.Split(docs, "\n")
	lines := make([]string, 0, len(rawLines))
	for _, l := range rawLines {
		lines = append(lines, strings.ReplaceAll(strings.TrimRight(l, " \t\r"), `"""`, `\"\"\"`))
	}
	return lines
}

// writeDocstring writes lines as a triple-quoted Python docstring,
// indented by indent, with optional extra pre-rendered lines (e.g.
// "Raises:" entries) appended after a blank line - mirrors how every
// other target's doc-writer appends @throws/<exception>-equivalent lines
// after the main description. Writes nothing if both are empty.
func writeDocstring(b *strings.Builder, indent string, lines []string, extra ...string) {
	if len(lines) == 0 && len(extra) == 0 {
		return
	}
	b.WriteString(indent + `"""`)
	switch {
	case len(lines) > 0:
		b.WriteString(lines[0] + "\n")
		for _, l := range lines[1:] {
			if l == "" {
				b.WriteString("\n")
				continue
			}
			b.WriteString(indent + l + "\n")
		}
		if len(extra) > 0 {
			b.WriteString("\n")
			for _, e := range extra {
				b.WriteString(indent + e + "\n")
			}
		}
	case len(extra) > 0:
		// No main description, only extra (Raises:) lines - PEP 257 style,
		// first content line right after the opening quotes, no leading
		// blank line.
		for i, e := range extra {
			if i == 0 {
				b.WriteString(e + "\n")
			} else {
				b.WriteString(indent + e + "\n")
			}
		}
	}
	b.WriteString(indent + `"""` + "\n")
}