package phpgen

import "strings"

// docLines splits a docs string (markdown, per the wire spec) into lines
// safe to place inside a PHP docblock comment - each returned line has no
// leading/trailing whitespace and can't accidentally close the comment
// early. Mirrors jsgen's docLines exactly (same underlying problem: PHP
// docblocks close on "*/" the same way JSDoc ones do).
func docLines(docs string) []string {
	docs = strings.TrimSpace(docs)
	if docs == "" {
		return nil
	}

	docs = strings.ReplaceAll(docs, "*/", "*\\/")

	rawLines := strings.Split(docs, "\n")
	lines := make([]string, 0, len(rawLines))
	for _, l := range rawLines {
		lines = append(lines, strings.TrimRight(l, " \t\r"))
	}
	return lines
}