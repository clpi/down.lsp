package handler

import (
	"context"
	"strings"

	"github.com/clpi/down.lsp/lsp/handler/completion/entries"
	"github.com/clpi/down.lsp/lsp/knowledge"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *State) Completion(
	c *glsp.Context,
	p *protocol.CompletionParams,
) (interface{}, error) {
	items := []protocol.CompletionItem{}
	items = entries.SnippetCompletions(items)
	items = entries.EmojiCompletions(items)
	items = entries.FileCompletions(items)
	items = entries.HtmlTagCompletions(items)
	items = entries.WorkspaceCompletions(items)
	items = s.knowledgeCompletions(items, p)
	items = s.aiCompletions(items, p)
	return items, nil
}

func (s *State) knowledgeCompletions(items []protocol.CompletionItem, p *protocol.CompletionParams) []protocol.CompletionItem {
	if s.Graph == nil {
		return items
	}

	uri := string(p.TextDocument.URI)
	doc, ok := s.Documents[uri]
	if !ok {
		return items
	}

	lines := strings.Split(doc, "\n")
	if int(p.Position.Line) >= len(lines) {
		return items
	}
	line := lines[p.Position.Line]
	col := int(p.Position.Character)
	if col > len(line) {
		col = len(line)
	}
	prefix := line[:col]

	var query string
	if idx := strings.LastIndexAny(prefix, " \t@#["); idx >= 0 {
		query = prefix[idx+1:]
	} else {
		query = prefix
	}

	if len(query) < 2 {
		return items
	}

	results := s.Graph.Search(query)
	kindAI := protocol.CompletionItemKindReference
	for _, ent := range results {
		detail := string(ent.Kind)
		if len(ent.Sources) > 1 {
			detail += " (referenced in multiple docs)"
		}
		items = append(items, protocol.CompletionItem{
			Label:  ent.Name,
			Kind:   &kindAI,
			Detail: &detail,
			Documentation: &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: entityDoc(ent),
			},
		})
	}
	return items
}

func entityDoc(ent *knowledge.Entity) string {
	var sb strings.Builder
	sb.WriteString("**" + ent.Name + "** (`" + string(ent.Kind) + "`)\n\n")
	if len(ent.Properties) > 0 {
		for k, v := range ent.Properties {
			sb.WriteString("- " + k + ": " + v + "\n")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("Mentions: " + intStr(ent.Mentions) + "\n")
	return sb.String()
}

func intStr(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

func (s *State) aiCompletions(items []protocol.CompletionItem, p *protocol.CompletionParams) []protocol.CompletionItem {
	if s.AI == nil {
		return items
	}

	uri := string(p.TextDocument.URI)
	doc, ok := s.Documents[uri]
	if !ok {
		return items
	}

	lines := strings.Split(doc, "\n")
	lineIdx := int(p.Position.Line)
	if lineIdx >= len(lines) {
		return items
	}

	currentLine := lines[lineIdx]
	col := int(p.Position.Character)
	if col > len(currentLine) {
		col = len(currentLine)
	}
	linePrefix := currentLine[:col]

	if len(strings.TrimSpace(linePrefix)) < 3 {
		return items
	}

	start := lineIdx - 20
	if start < 0 {
		start = 0
	}
	precedingText := strings.Join(lines[start:lineIdx], "\n")

	ctx, cancel := context.WithTimeout(context.Background(), 5*1e9)
	defer cancel()

	completions, err := s.AI.CompleteText(ctx, uri, precedingText, linePrefix)
	if err != nil {
		return items
	}

	kindAI := protocol.CompletionItemKindText
	for i, comp := range completions {
		label := comp
		if len(label) > 60 {
			label = label[:57] + "..."
		}
		detail := "AI suggestion"
		sortText := "zzz" + string(rune('0'+i))
		items = append(items, protocol.CompletionItem{
			Label:      label,
			Kind:       &kindAI,
			Detail:     &detail,
			InsertText: &comp,
			SortText:   &sortText,
		})
	}
	return items
}

func (s *State) ItemResolve(c *glsp.Context, p *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return p, nil
}
