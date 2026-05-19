package semantic

import (
	"regexp"
	"sort"
	"strings"
)

var (
	reHeading    = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	reTag        = regexp.MustCompile(`(?:^|\s)#([A-Za-z][A-Za-z0-9_/-]*)`)
	reMention    = regexp.MustCompile(`(?:^|\s)@([A-Za-z][A-Za-z0-9_.-]*)`)
	reWikiLink   = regexp.MustCompile(`\[\[([^\]]+)\]\]`)
	reMdLink     = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	reTask       = regexp.MustCompile(`^(\s*[-*+] \[[ xX]\])`)
	reCodeSpan   = regexp.MustCompile("`([^`]+)`")
	reDate       = regexp.MustCompile(`\b(\d{4}-\d{2}-\d{2})\b`)
	reFMKey      = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*):\s`)
	reBold       = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	reItalic     = regexp.MustCompile(`(?:^|[^*])\*([^*]+)\*(?:[^*]|$)`)
	reBlockquote = regexp.MustCompile(`^(>\s?.*)$`)
	reCodeFence  = regexp.MustCompile("^```")
)

// Tokenize scans a markdown document and returns sorted semantic tokens.
func Tokenize(text string) []Token {
	lines := strings.Split(text, "\n")
	var tokens []Token

	inFrontmatter := false
	inCodeBlock := false
	fmStarted := false

	for lineIdx, line := range lines {
		// Track frontmatter boundaries
		if lineIdx == 0 && strings.TrimSpace(line) == "---" {
			inFrontmatter = true
			fmStarted = true
			continue
		}
		if inFrontmatter && strings.TrimSpace(line) == "---" {
			inFrontmatter = false
			continue
		}
		if inFrontmatter {
			if m := reFMKey.FindStringIndex(line); m != nil {
				keyEnd := strings.Index(line, ":")
				tokens = append(tokens, Token{
					Line:      lineIdx,
					StartChar: 0,
					Length:    keyEnd,
					Type:      TokFrontmatter,
				})
			}
			continue
		}

		// Track code fences
		if reCodeFence.MatchString(strings.TrimSpace(line)) {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			continue
		}

		// Headings
		if m := reHeading.FindStringSubmatchIndex(line); m != nil {
			tokens = append(tokens, Token{
				Line:      lineIdx,
				StartChar: 0,
				Length:    len(line),
				Type:      TokHeading,
				Modifiers: ModDeclaration,
			})
			// Still scan the heading text for inline elements below
		}

		// Blockquotes
		if m := reBlockquote.FindStringIndex(line); m != nil {
			tokens = append(tokens, Token{
				Line:      lineIdx,
				StartChar: m[0],
				Length:    m[1] - m[0],
				Type:      TokBlockquote,
			})
		}

		// Tasks
		if m := reTask.FindStringIndex(line); m != nil {
			mod := 0
			if strings.Contains(line[m[0]:m[1]], "[x]") || strings.Contains(line[m[0]:m[1]], "[X]") {
				mod = ModDeprecated // completed task
			}
			tokens = append(tokens, Token{
				Line:      lineIdx,
				StartChar: m[0],
				Length:    m[1] - m[0],
				Type:      TokTask,
				Modifiers: mod,
			})
		}

		// Tags (#tag)
		for _, m := range reTag.FindAllStringIndex(line, -1) {
			// Adjust start to skip leading whitespace — find the '#'
			start := m[0]
			for start < m[1] && line[start] != '#' {
				start++
			}
			tokens = append(tokens, Token{
				Line:      lineIdx,
				StartChar: start,
				Length:    m[1] - start,
				Type:      TokTag,
			})
		}

		// Mentions (@person)
		for _, m := range reMention.FindAllStringIndex(line, -1) {
			start := m[0]
			for start < m[1] && line[start] != '@' {
				start++
			}
			tokens = append(tokens, Token{
				Line:      lineIdx,
				StartChar: start,
				Length:    m[1] - start,
				Type:      TokMention,
			})
		}

		// Wiki links [[target]]
		for _, m := range reWikiLink.FindAllStringIndex(line, -1) {
			tokens = append(tokens, Token{
				Line:      lineIdx,
				StartChar: m[0],
				Length:    m[1] - m[0],
				Type:      TokWikiLink,
				Modifiers: ModLink,
			})
		}

		// Markdown links [text](url)
		for _, m := range reMdLink.FindAllStringIndex(line, -1) {
			tokens = append(tokens, Token{
				Line:      lineIdx,
				StartChar: m[0],
				Length:    m[1] - m[0],
				Type:      TokMdLink,
				Modifiers: ModLink,
			})
		}

		// Code spans `code`
		for _, m := range reCodeSpan.FindAllStringIndex(line, -1) {
			tokens = append(tokens, Token{
				Line:      lineIdx,
				StartChar: m[0],
				Length:    m[1] - m[0],
				Type:      TokCodeSpan,
			})
		}

		// Dates (YYYY-MM-DD)
		for _, m := range reDate.FindAllStringIndex(line, -1) {
			tokens = append(tokens, Token{
				Line:      lineIdx,
				StartChar: m[0],
				Length:    m[1] - m[0],
				Type:      TokDate,
			})
		}

		// Bold **text**
		for _, m := range reBold.FindAllStringIndex(line, -1) {
			tokens = append(tokens, Token{
				Line:      lineIdx,
				StartChar: m[0],
				Length:    m[1] - m[0],
				Type:      TokBold,
			})
		}

		// Italic *text* (not inside bold)
		for _, m := range reItalic.FindAllStringIndex(line, -1) {
			// Find the actual * position
			start := m[0]
			for start < m[1] && line[start] != '*' {
				start++
			}
			end := m[1]
			for end > start && line[end-1] != '*' {
				end--
			}
			if start < end {
				tokens = append(tokens, Token{
					Line:      lineIdx,
					StartChar: start,
					Length:    end - start,
					Type:      TokItalic,
				})
			}
		}
	}

	// Suppress the unused variable
	_ = fmStarted

	// Sort by line, then by start character
	sort.Slice(tokens, func(i, j int) bool {
		if tokens[i].Line != tokens[j].Line {
			return tokens[i].Line < tokens[j].Line
		}
		return tokens[i].StartChar < tokens[j].StartChar
	})

	return tokens
}
