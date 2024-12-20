package document

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func References(c *glsp.Context, p *protocol.ReferenceParams) ([]protocol.Location, error) {
	var refs []protocol.Location
	return append(refs, protocol.Location{}), nil
}
