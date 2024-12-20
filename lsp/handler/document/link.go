package document

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func LinkResolve(c *glsp.Context, p *protocol.DocumentLink) (*protocol.DocumentLink, error) {
	return p, nil
}
func Link(c *glsp.Context, p *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	var (
		dls []protocol.DocumentLink
		d1  protocol.DocumentLink = protocol.DocumentLink{
			Tooltip: &p.TextDocument.URI,
			Data:    nil,
			Target:  &p.TextDocument.URI,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      10,
					Character: 10,
				},
				End: protocol.Position{
					Line:      20,
					Character: 20,
				},
			},
		}
	)
	return append(dls, d1), nil
}
