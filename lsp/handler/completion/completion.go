package completion

import (
	"github.com/clpi/down.lsp/lsp/handler/completion/entries"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	trueVal  = true
	falseVal = false

	Provider protocol.CompletionOptions = protocol.CompletionOptions{
		ResolveProvider: &trueVal,
		WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
			WorkDoneProgress: &trueVal,
		},
		AllCommitCharacters: []string{
			" ", "@", "#", "$", "%", "&",
		},
		TriggerCharacters: []string{
			" ",
			"@", "#", "$", "%", "&",
			"*", "+", "-", "/", "<", "=",
			">", "?", "^", "|", "~",
			"[", "(", "<", "{", "`",
			"]", ")", ">", "}",
			":", "=", ",",
			".", ";", "'",
			"\"", "'", "\\", "/",
			"!", "_",
			"~", "`",
		},
	}
)

func Completion(
	c *glsp.Context,
	p *protocol.CompletionParams,
) (interface{}, error) {
	var (
		preselect = true
		items     []protocol.CompletionItem
	)
	for s, sn := range entries.Snippets {
		kind := protocol.CompletionItemKindSnippet
		items = append(items, protocol.CompletionItem{
			Label:     s,
			Kind:      &kind,
			Preselect: &preselect,
			Documentation: &protocol.MarkupContent{
				Value: "# Snippets\n\n## Snippet\n_ _ _\n### Snippet: " + sn.Description + "\n---\n" + sn.Body,
				Kind:  protocol.MarkupKindMarkdown,
			},
			Detail:     &sn.Description,
			InsertText: &sn.Body,
		})

	}
	for w, e := range entries.Emojis {
		ec := e
		kind := protocol.CompletionItemKindConstant
		items = append(items, protocol.CompletionItem{
			Label:     w,
			Kind:      &kind,
			Preselect: &preselect,
			Documentation: &protocol.MarkupContent{
				Value: "# Emoji\n\n## Emoji\n_ _ _\n### Emoji: " + ec + "\n---\n" + ec,
				Kind:  protocol.MarkupKindMarkdown,
			},
			Detail:     &ec,
			InsertText: &ec,
		})
	}
	return items, nil
}
