package csharpgen

import "strings"

// docLines splits a docs string (markdown, per the wire spec) into lines
// safe to place inside a series of "/// " XML doc comment lines. Unlike
// Javadoc/JSDoc/PHPDoc block comments, C#'s XML doc comments have no
// closing token to worry about (each line is its own "///" comment), but
// the content is technically XML, so literal "<", ">", and "&" need
// escaping or a later doc line could be misread as a tag/entity.
func docLines(docs string) []string {
	docs = strings.TrimSpace(docs)
	if docs == "" {
		return nil
	}
	rawLines := strings.Split(docs, "\n")
	lines := make([]string, 0, len(rawLines))
	for _, l := range rawLines {
		lines = append(lines, xmlEscape(strings.TrimRight(l, " \t\r")))
	}
	return lines
}

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// writeSummary writes lines as a "/// <summary>...</summary>" XML doc
// block, indented by indent. Writes nothing if lines is empty. Additional
// pre-rendered XML doc lines (e.g. "/// <param ...>", "/// <exception ...>")
// can be passed via extra, appended after the summary block - mirrors how
// every other target's doc-writer appends @param/@throws-equivalent lines
// after the main description.
func writeSummary(b *strings.Builder, indent string, lines []string, extra ...string) {
	if len(lines) == 0 && len(extra) == 0 {
		return
	}
	if len(lines) > 0 {
		b.WriteString(indent + "/// <summary>\n")
		for _, l := range lines {
			if l == "" {
				b.WriteString(indent + "///\n")
				continue
			}
			b.WriteString(indent + "/// " + l + "\n")
		}
		b.WriteString(indent + "/// </summary>\n")
	}
	for _, e := range extra {
		b.WriteString(indent + "/// " + e + "\n")
	}
}
