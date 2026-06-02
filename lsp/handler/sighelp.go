package handler

import (
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// MarkdownSignatureHelp provides contextual signature help for markdown constructs.
// Triggers on [, (, !, >, |, $ to show parameter hints for links, images, tables, etc.
func (s *State) MarkdownSignatureHelp(_ *glsp.Context, p *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	uri := string(p.TextDocument.URI)
	text, ok := s.Documents[uri]
	if !ok {
		return nil, nil
	}

	lines := strings.Split(text, "\n")
	lineIdx := int(p.Position.Line)
	if lineIdx >= len(lines) {
		return nil, nil
	}
	line := lines[lineIdx]
	col := int(p.Position.Character)
	if col > len(line) {
		col = len(line)
	}
	prefix := line[:col]

	// Detect which construct we're inside
	sig := detectSignatureContext(prefix, line)
	if sig == nil {
		return nil, nil
	}

	return sig, nil
}

func detectSignatureContext(prefix, fullLine string) *protocol.SignatureHelp {
	// Image: ![alt](url "title")
	if idx := strings.LastIndex(prefix, "!["); idx >= 0 {
		after := prefix[idx:]
		if !strings.Contains(after, ")") {
			return imageSignature(after)
		}
	}

	// Markdown link: [text](url "title")
	if idx := lastUnmatchedBracket(prefix); idx >= 0 {
		after := prefix[idx:]
		if !strings.Contains(after, ")") {
			return linkSignature(after)
		}
	}

	// Table: | cell | cell |
	if strings.Contains(prefix, "|") && !strings.HasPrefix(strings.TrimSpace(prefix), "```") {
		return tableSignature(prefix)
	}

	// Callout: > [!type]
	trimmed := strings.TrimSpace(prefix)
	if strings.HasPrefix(trimmed, "> [!") && !strings.Contains(trimmed, "]") {
		return calloutSignature()
	}

	// Math: $...$
	dollarCount := strings.Count(prefix, "$") - strings.Count(prefix, "\\$")
	if dollarCount%2 == 1 {
		return mathSignature()
	}

	// Frontmatter key: at start of document in --- block
	if strings.Contains(prefix, ":") && !strings.HasPrefix(trimmed, "#") {
		// Could be frontmatter — check if we're between --- markers
		return nil // frontmatter help handled elsewhere
	}

	return nil
}

func imageSignature(context string) *protocol.SignatureHelp {
	activeParam := protocol.UInteger(0)
	activeSig := protocol.UInteger(0)

	// Determine which parameter we're in
	if strings.Contains(context, "](") {
		activeParam = 1 // URL parameter
		if strings.Contains(context[strings.Index(context, "]("):], " \"") {
			activeParam = 2 // title parameter
		}
	}

	altDoc := "Alternative text displayed when image can't load"
	urlDoc := "URL or relative path to the image file"
	titleDoc := "Optional tooltip shown on hover (in quotes)"

	return &protocol.SignatureHelp{
		ActiveSignature: &activeSig,
		ActiveParameter: &activeParam,
		Signatures: []protocol.SignatureInformation{
			{
				Label:         "![alt text](url \"title\")",
				Documentation: "Insert an image with alt text, URL, and optional title",
				Parameters: []protocol.ParameterInformation{
					{Label: "alt text", Documentation: altDoc},
					{Label: "url", Documentation: urlDoc},
					{Label: "\"title\"", Documentation: titleDoc},
				},
			},
		},
	}
}

func linkSignature(context string) *protocol.SignatureHelp {
	activeParam := protocol.UInteger(0)
	activeSig := protocol.UInteger(0)

	if strings.Contains(context, "](") {
		activeParam = 1
		if strings.Contains(context[strings.Index(context, "]("):], " \"") {
			activeParam = 2
		}
	}

	textDoc := "Visible link text"
	urlDoc := "URL, relative path, or #heading-anchor"
	titleDoc := "Optional tooltip shown on hover (in quotes)"

	return &protocol.SignatureHelp{
		ActiveSignature: &activeSig,
		ActiveParameter: &activeParam,
		Signatures: []protocol.SignatureInformation{
			{
				Label:         "[link text](url \"title\")",
				Documentation: "Insert a hyperlink with display text, target URL, and optional title",
				Parameters: []protocol.ParameterInformation{
					{Label: "link text", Documentation: textDoc},
					{Label: "url", Documentation: urlDoc},
					{Label: "\"title\"", Documentation: titleDoc},
				},
			},
			{
				Label:         "[[wiki link|display text]]",
				Documentation: "Insert a wiki-style link to another document in the workspace",
				Parameters: []protocol.ParameterInformation{
					{Label: "wiki link", Documentation: "Target document name (without extension)"},
					{Label: "display text", Documentation: "Optional display text (after |)"},
				},
			},
		},
	}
}

func tableSignature(prefix string) *protocol.SignatureHelp {
	activeSig := protocol.UInteger(0)
	// Count pipe characters to determine active column
	pipes := strings.Count(prefix, "|")
	activeParam := protocol.UInteger(0)
	if pipes > 0 {
		activeParam = protocol.UInteger(pipes - 1)
	}

	return &protocol.SignatureHelp{
		ActiveSignature: &activeSig,
		ActiveParameter: &activeParam,
		Signatures: []protocol.SignatureInformation{
			{
				Label:         "| Column 1 | Column 2 | Column 3 |",
				Documentation: "Markdown table row. Use `| --- |` for the header separator. Align with `:---`, `:---:`, `---:`",
				Parameters: []protocol.ParameterInformation{
					{Label: "Column 1", Documentation: "First column cell content"},
					{Label: "Column 2", Documentation: "Second column cell content"},
					{Label: "Column 3", Documentation: "Additional columns..."},
				},
			},
		},
	}
}

func calloutSignature() *protocol.SignatureHelp {
	activeSig := protocol.UInteger(0)
	activeParam := protocol.UInteger(0)

	return &protocol.SignatureHelp{
		ActiveSignature: &activeSig,
		ActiveParameter: &activeParam,
		Signatures: []protocol.SignatureInformation{
			{
				Label:         "> [!type] Title",
				Documentation: "Callout/admonition block. Types: note, tip, warning, danger, info, example, quote, abstract, bug, success, failure, question",
				Parameters: []protocol.ParameterInformation{
					{Label: "type", Documentation: "Callout type: note | tip | warning | danger | info | example | quote | abstract | bug | success | failure | question"},
					{Label: "Title", Documentation: "Optional custom title (defaults to type name)"},
				},
			},
		},
	}
}

func mathSignature() *protocol.SignatureHelp {
	activeSig := protocol.UInteger(0)
	activeParam := protocol.UInteger(0)

	return &protocol.SignatureHelp{
		ActiveSignature: &activeSig,
		ActiveParameter: &activeParam,
		Signatures: []protocol.SignatureInformation{
			{
				Label:         "$LaTeX formula$",
				Documentation: "Inline math using LaTeX syntax. Use $$ for display/block math.",
				Parameters: []protocol.ParameterInformation{
					{Label: "formula", Documentation: "LaTeX math expression (e.g., \\frac{a}{b}, x^2, \\sum_{i=1}^n)"},
				},
			},
		},
	}
}

// lastUnmatchedBracket finds the last '[' that hasn't been closed with ']('.
func lastUnmatchedBracket(s string) int {
	depth := 0
	for i := len(s) - 1; i >= 0; i-- {
		switch s[i] {
		case ')':
			depth++
		case '(':
			depth--
		case '[':
			if depth <= 0 && i > 0 && s[i-1] != '!' {
				return i
			}
			if depth <= 0 && i == 0 {
				return i
			}
		}
	}
	return -1
}
