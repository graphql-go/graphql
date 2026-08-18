package parser

import (
	"strings"
	"testing"
)

// A document that nests far deeper than any real query must be rejected with a
// syntax error rather than being allowed to exhaust the goroutine stack.
func TestParseDeeplyNestedInputDoesNotOverflowStack(t *testing.T) {
	cases := map[string]string{
		"list value":    "{f(a:" + strings.Repeat("[", 200000) + ")}",
		"object value":  "{f(a:" + strings.Repeat("{b:", 200000) + "}",
		"selection set": strings.Repeat("{f", 200000),
		"list type":     "query($v:" + strings.Repeat("[", 200000) + "Int){f}",
	}
	for name, src := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Parse(ParseParams{Source: src})
			if err == nil {
				t.Fatalf("expected a syntax error for deeply nested input, got nil")
			}
			if !strings.Contains(err.Error(), "nests deeper") {
				t.Fatalf("expected depth-limit error, got: %v", err)
			}
		})
	}
}

// Reasonably nested documents must still parse successfully.
func TestParseModeratelyNestedInputStillParses(t *testing.T) {
	depth := 100
	src := "{f(a:" + strings.Repeat("[", depth) + strings.Repeat("]", depth) + ")}"
	if _, err := Parse(ParseParams{Source: src}); err != nil {
		t.Fatalf("moderately nested document should parse, got: %v", err)
	}
}
