package handler

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *State) Hover(c *glsp.Context, p *protocol.HoverParams) (*protocol.Hover, error) {
	// _, _, _ = util.ReadLnChar(p.TextDocument.URI, p.Position)
	// _, ch, err := util.ReadLnChar(p.TextDocument.URI, p.Position)
	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: "# hi\n\n### hover\n\n###\n\n#### ok\n",
			// string(p.Position.Line) + " " + string(p.Position.Character) +
			// `\n\n### ` + string(p.TextDocument.URI) + `\n\n` +
			// string(p.TextDocumentPositionParams.Position.Line) + " " + string(p.TextDocumentPositionParams.Position.Character) + `\n\n### ` + string(p.TextDocumentPositionParams.TextDocument.URI) + `\n\n`,
		},
		Range: &protocol.Range{
			Start: p.Position,
			End: protocol.Position{
				Line:      p.Position.Line,
				Character: p.Position.Character + 1,
			},
		},
	}, nil
}
