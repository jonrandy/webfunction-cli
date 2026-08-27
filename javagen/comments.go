package javagen

import "strings"

// docLines splits a docs string (markdown, per the wire spec) into lines
// safe to place inside a Javadoc "/** ... */" block. Unlike gogen's
// docLines (plain "//" comments have no closing token to worry about),
// this escapes any literal "*/" sequence a docs string might contain -
// mirrors jsgen's/phpgen's docLines, which face the identical hazard for
// their own "/** */" and "/* */" comment styles.
func docLines(docs string) []string {
	docs = strings.TrimSpace(docs)
	if docs == "" {
		return nil
	}
	rawLines := strings.Split(docs, "\n")
	lines := make([]string, 0, len(rawLines))
	for _, l := range rawLines {
		l = strings.ReplaceAll(l, "*/", "*\\/")
		lines = append(lines, strings.TrimRight(l, " \t\r"))
	}
	return lines
}

// writeJavadoc writes lines as a full "/** ... */" Javadoc block,
// indented by indent. Writes nothing if lines is empty.
func writeJavadoc(b *strings.Builder, indent string, lines []string) {
	if len(lines) == 0 {
		return
	}
	b.WriteString(indent + "/**\n")
	for _, l := range lines {
		if l == "" {
			b.WriteString(indent + " *\n")
			continue
		}
		b.WriteString(indent + " * " + l + "\n")
	}
	b.WriteString(indent + " */\n")
}
