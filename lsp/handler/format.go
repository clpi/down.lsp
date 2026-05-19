package handler

import (
	"regexp"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	reTrailingWS    = regexp.MustCompile(`[ \t]+$`)
	reMultipleBlank = regexp.MustCompile(`\n{3,}`)
	reHeadingSpace  = regexp.MustCompile(`^(#{1,6})([^ #\n])`)
	reListIndent    = regexp.MustCompile(`^(\s*)([*+-]|\d+\.) {2,}`)
)

// Format implements textDocument/formatting.
// It normalizes whitespace, heading spacing, and list formatting.
func (s *State) Format(_ *glsp.Context, p *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	uri := string(p.TextDocument.URI)
	text, ok := s.Documents[uri]
	if !ok {
		return nil, nil
	}

	formatted := formatMarkdown(text, p.Options)
	if formatted == text {
		return nil, nil // no changes
	}

	// Replace the entire document
	lines := strings.Split(text, "\n")
	lastLine := len(lines) - 1
	lastChar := len(lines[lastLine])

	return []protocol.TextEdit{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End: protocol.Position{
					Line:      protocol.UInteger(lastLine),
					Character: protocol.UInteger(lastChar),
				},
			},
			NewText: formatted,
		},
	}, nil
}

func formatMarkdown(text string, opts protocol.FormattingOptions) string {
	lines := strings.Split(text, "\n")
	var result []string

	inCodeBlock := false
	inFrontmatter := false
	prevBlank := false

	for i, line := range lines {
		// Track frontmatter
		if i == 0 && strings.TrimSpace(line) == "---" {
			inFrontmatter = true
			result = append(result, line)
			prevBlank = false
			continue
		}
		if inFrontmatter {
			if strings.TrimSpace(line) == "---" {
				inFrontmatter = false
			}
			result = append(result, line)
			prevBlank = false
			continue
		}

		// Track code blocks — don't format inside them
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			result = append(result, line)
			prevBlank = false
			continue
		}
		if inCodeBlock {
			result = append(result, line)
			prevBlank = false
			continue
		}

		// Remove trailing whitespace. Preserve intentional line breaks (2 trailing spaces).
		stripped := strings.TrimRight(line, " \t")
		if len(line) > len(stripped) && len(line)-len(stripped) >= 2 {
			// Markdown line break: normalize to exactly 2 trailing spaces
			line = stripped + "  "
		} else {
			line = stripped
		}

		// Ensure space after heading marker: ##Foo → ## Foo
		line = reHeadingSpace.ReplaceAllString(line, "$1 $2")

		// Normalize list item spacing: -  item → - item
		if m := reListIndent.FindStringSubmatch(line); m != nil {
			rest := strings.TrimSpace(line[len(m[0]):])
			line = m[1] + m[2] + " " + rest
		}

		// Collapse multiple blank lines
		isBlank := strings.TrimSpace(line) == ""
		if isBlank && prevBlank {
			continue // skip extra blank line
		}

		// Ensure blank line before headings (unless at start)
		if strings.HasPrefix(trimmed, "#") && len(result) > 0 && strings.TrimSpace(result[len(result)-1]) != "" {
			result = append(result, "")
		}

		result = append(result, line)
		prevBlank = isBlank
	}

	// Ensure file ends with a single newline
	out := strings.Join(result, "\n")
	out = strings.TrimRight(out, "\n") + "\n"
	return out
}
