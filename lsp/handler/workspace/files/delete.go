package files

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Delete(c *glsp.Context, p *protocol.DeleteFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}

func DidDelete(c *glsp.Context, p *protocol.DeleteFilesParams) error {
	return nil
}
