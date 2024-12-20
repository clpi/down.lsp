package document

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Moniker(c *glsp.Context, p *protocol.MonikerParams) ([]protocol.Moniker, error) {
	var (
		k1       protocol.MonikerKind = protocol.MonikerKindLocal
		monikers []protocol.Moniker
		m1       protocol.Moniker = protocol.Moniker{
			Unique:     "unique",
			Identifier: "identifier",
			Kind:       &k1,
			Scheme:     "file",
		}
	)
	return append(monikers, m1), nil
}
