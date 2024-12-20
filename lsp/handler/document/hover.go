package document

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Hover(c *glsp.Context, p *protocol.HoverParams) (*protocol.Hover, error) {
	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind: protocol.MarkupKindMarkdown,
			Value: `
      # down-lsp

          info: down lsp

      ## Info

      ### Information

      #### About

      This is a hover tooltip

      `,
		},
		Range: &protocol.Range{
			Start: protocol.Position{
				Line:      0,
				Character: 0,
			},
			End: protocol.Position{
				Line:      0,
				Character: 0,
			},
		},
	}, nil
}
