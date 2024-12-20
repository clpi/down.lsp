package lens

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Resolve(c *glsp.Context, p *protocol.CodeLens) (*protocol.CodeLens, error) {
	return p, nil
}
