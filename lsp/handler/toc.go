package handler

import (
	"fmt"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// GenerateTOC builds a markdown Table of Contents from the document headings.
// Returns a workspace edit that inserts the TOC at the specified position.
func (s *State) GenerateTOC(uri string, insertLine int) *protocol.WorkspaceEdit {
	text, ok := s.Documents[uri]
	if !ok {
		return nil
	}

	lines := strings.Split(text, "\n")
	var tocLines []string
	tocLines = append(tocLines, "## Table of Contents", "")

	inCode := false
	inFrontmatter := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip frontmatter
		if i == 0 && trimmed == "---" {
			inFrontmatter = true
			continue
		}
		if inFrontmatter {
			if trimmed == "---" {
				inFrontmatter = false
			}
			continue
		}

		// Skip code blocks
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inCode = !inCode
			continue
		}
		if inCode {
			continue
		}

		// Process headings
		if strings.HasPrefix(trimmed, "#") {
			level := 0
			for _, ch := range trimmed {
				if ch == '#' {
					level++
				} else {
					break
				}
			}
			if level >= 2 && level <= 6 { // Skip H1 (usually the doc title)
				headingText := strings.TrimSpace(trimmed[level:])
				if headingText == "" || strings.ToLower(headingText) == "table of contents" {
					continue
				}
				indent := strings.Repeat("  ", level-2)
				slug := tocSlug(headingText)
				tocLines = append(tocLines, fmt.Sprintf("%s- [%s](#%s)", indent, headingText, slug))
			}
		}
	}

	if len(tocLines) <= 2 {
		return nil // No headings found
	}

	tocLines = append(tocLines, "") // trailing newline
	tocText := strings.Join(tocLines, "\n")

	// Build workspace edit
	changes := map[protocol.DocumentUri][]protocol.TextEdit{
		protocol.DocumentUri(uri): {
			{
				Range: protocol.Range{
					Start: protocol.Position{Line: protocol.UInteger(insertLine), Character: 0},
					End:   protocol.Position{Line: protocol.UInteger(insertLine), Character: 0},
				},
				NewText: tocText + "\n",
			},
		},
	}

	return &protocol.WorkspaceEdit{Changes: changes}
}

// tocSlug creates a GitHub-flavored anchor slug from heading text.
func tocSlug(text string) string {
	text = strings.ToLower(text)
	var result []byte
	prevDash := false
	for _, b := range []byte(text) {
		if (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') {
			result = append(result, b)
			prevDash = false
		} else if b == ' ' || b == '-' || b == '_' {
			if !prevDash && len(result) > 0 {
				result = append(result, '-')
				prevDash = true
			}
		}
		// Skip all other characters (punctuation, etc.)
	}
	// Trim trailing dash
	for len(result) > 0 && result[len(result)-1] == '-' {
		result = result[:len(result)-1]
	}
	return string(result)
}

// TOCCodeAction creates a code action that inserts a TOC when triggered.
func (s *State) TOCCodeAction(_ *glsp.Context, p *protocol.CodeActionParams) *protocol.CodeAction {
	uri := string(p.TextDocument.URI)
	text, ok := s.Documents[uri]
	if !ok {
		return nil
	}

	// Check if document has enough headings to warrant a TOC
	headingCount := 0
	lines := strings.Split(text, "\n")
	hasTOC := false
	inCode := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCode = !inCode
			continue
		}
		if inCode {
			continue
		}
		if strings.HasPrefix(trimmed, "##") {
			headingCount++
		}
		if strings.ToLower(trimmed) == "## table of contents" {
			hasTOC = true
		}
	}

	if headingCount < 3 || hasTOC {
		return nil
	}

	// Find insertion point (after frontmatter + first heading, or at top)
	insertLine := 0
	inFM := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if i == 0 && trimmed == "---" {
			inFM = true
			continue
		}
		if inFM {
			if trimmed == "---" {
				inFM = false
				insertLine = i + 1
			}
			continue
		}
		if strings.HasPrefix(trimmed, "# ") {
			insertLine = i + 1
			break
		}
	}

	// Skip blank lines after heading
	for insertLine < len(lines) && strings.TrimSpace(lines[insertLine]) == "" {
		insertLine++
	}

	edit := s.GenerateTOC(uri, insertLine)
	if edit == nil {
		return nil
	}

	kind := protocol.CodeActionKindSource
	title := "Generate Table of Contents"
	return &protocol.CodeAction{
		Title: title,
		Kind:  &kind,
		Edit:  edit,
	}
}
