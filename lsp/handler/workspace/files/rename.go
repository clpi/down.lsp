package files

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Rename(c *glsp.Context, p *protocol.RenameFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}

func DidRename(c *glsp.Context, p *protocol.RenameFilesParams) error {
	return nil
}
