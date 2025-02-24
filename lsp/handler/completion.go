package handler

import (
	"github.com/clpi/down.lsp/lsp/handler/completion/entries"
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
	return items, nil
}

func (s *State) ItemResolve(c *glsp.Context, p *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return p, nil
}
