package semantic

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Full(c *glsp.Context, p *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	var (
		rid  string = "result-id"
		data        = []protocol.UInteger{
			protocol.UInteger(10),
			protocol.UInteger(20),
			protocol.UInteger(30),
		}
		st protocol.SemanticTokens = protocol.SemanticTokens{
			Data:     data,
			ResultID: &rid,
		}
	)
	return &st, nil
}
