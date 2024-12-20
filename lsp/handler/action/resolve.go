package action

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Resolve(c *glsp.Context, p *protocol.CodeAction) (*protocol.CodeAction, error) {
	return p, nil
}
