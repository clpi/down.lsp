package document

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Color(c *glsp.Context, p *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	var ci []protocol.ColorInformation
	return append(ci, protocol.ColorInformation{}), nil
}

func ColorPresentation(c *glsp.Context, p *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	return []protocol.ColorPresentation{}, nil
}
