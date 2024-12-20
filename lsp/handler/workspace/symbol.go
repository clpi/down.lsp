package workspace

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Symbol(*glsp.Context, *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	return []protocol.SymbolInformation{}, nil
}
