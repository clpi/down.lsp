package handler

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *State) Full(c *glsp.Context, p *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
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

func (s *State) Delta(c *glsp.Context, p *protocol.SemanticTokensDeltaParams) (any, error) {
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

func (s *State) Refresh(c *glsp.Context) error {
	var ()
	return nil
}

func (s *State) Range(c *glsp.Context, p *protocol.SemanticTokensRangeParams) (any, error) {
	var ()
	return nil, nil
}
