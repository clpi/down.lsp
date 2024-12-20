package completion

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func ItemResolve(c *glsp.Context, p *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return p, nil
}
