package handler

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *State) DocumentHighlight(c *glsp.Context, p *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	var (
		kk                              = protocol.DocumentHighlightKindText
		hl []protocol.DocumentHighlight = []protocol.DocumentHighlight{}
		k1                              = protocol.DocumentHighlightKindText
		_                               = protocol.DocumentHighlight{
			Kind: &kk,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      0,
					Character: 5,
				},
				End: protocol.Position{
					Line:      0,
					Character: 5,
				},
			},
		}
		h1 protocol.DocumentHighlight = protocol.DocumentHighlight{
			Kind: &k1,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      0,
					Character: 10,
				},
				End: protocol.Position{
					Line:      0,
					Character: 10,
				},
			},
		}
	)
	return append(hl, h1), nil
}
